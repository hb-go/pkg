package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/urfave/cli/v2"

	"github.com/hb-go/pkg/config/source"
)

func test(t *testing.T, withContext bool) {
	var src source.Source

	// setup app
	app := newCmd().App()
	app.Name = "testapp"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "db-host",
		},
		&cli.StringFlag{
			Name:    "db-host-as",
			Aliases: []string{"db_host"},
		},
	}

	// with context
	if withContext {
		// set action
		app.Action = func(c *cli.Context) error {
			src = WithContext(c)
			return nil
		}

		// run app
		if err := app.Run([]string{"run", "-db-host", "localhost", "-db-host-as", "localhost"}); err != nil {
			t.Error(err)
		}

		// no context
	} else {
		// set args
		os.Args = []string{"run", "-db-host", "localhost", "-db-host-as", "localhost"}
		src = NewSource(app)
	}

	// test config
	c, err := src.Read()
	if err != nil {
		t.Error(err)
	}
	if len(c.Data) == 0 {
		t.Fatal()
	}

	t.Log(string(c.Data))

	var actual *simplejson.Json
	if actual, err = simplejson.NewJson(c.Data); err != nil {
		t.Error(err)
	}

	if actual.Get("db-host").MustString() != "localhost" {
		t.Errorf("expected localhost, got %v", actual.Get("db-host").MustString())
	}

	if actual.Get("db").Get("host").MustString() != "localhost" {
		t.Errorf("expected localhost, got %v", actual.Get("db").Get("host").MustString())
	}
}

func TestCliSource(t *testing.T) {
	// without context
	test(t, false)
}

func TestCliSourceWithContext(t *testing.T) {
	// with context
	test(t, true)
}

func TestCliSource_cmd(t *testing.T) {
	// setup app
	app := newCmd().App()
	app.Name = "testcmd"
	app.Flags = DefaultFlags

	var src source.Source
	// set action

	cmd := &cli.Command{
		Name: "cmd",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cmd-option",
				Aliases: []string{"cmd_option"},
			},
		},
	}

	cmd.Action = func(ctx *cli.Context) error {
		src = WithContext(ctx)

		// test config
		c, err := src.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(c.Data) == 0 {
			t.Fatal()
		}

		t.Log(string(c.Data))

		var conf config
		if err := json.Unmarshal(c.Data, &conf); err != nil {
			t.Error(err)
		}

		// test default
		if conf.Config != "1" {
			t.Fatal(fmt.Errorf("config should be [1], not: [%s]", conf.Config))
		}

		// test the config from cmd
		if conf.Cmd.Option != "1" {
			t.Fatal(fmt.Errorf("cmd option should be [1] which is cmd value, not: [%s]", conf.Cmd.Option))
		}

		return nil
	}

	app.Commands = cli.Commands{
		cmd,
	}

	// set args
	os.Args = []string{"run"}
	for _, v := range DefaultFlags {
		os.Args = append(os.Args, fmt.Sprintf("--%s", v.Names()[0]), "1")
	}
	os.Args = append(os.Args, "cmd")
	for _, v := range app.Commands[0].Flags {
		os.Args = append(os.Args, fmt.Sprintf("--%s", v.Names()[0]), "1")
	}

	// run app
	if err := app.Run(os.Args); err != nil {
		t.Error(err)
	}
	// src := NewSource(app)
}

type CmdOption struct {
	Option string `json:"option"`
}

type config struct {
	Config string    `json:"config"`
	Cmd    CmdOption `json:"cmd"`
}

type Cmd interface {
	// The cli app within this cmd
	App() *cli.App
	// Adds options, parses flags and initialise
	// exits on error
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
	// ConfigFile path. This is not good
	ConfigFile() string
}

type cmd struct {
	opts Options
	app  *cli.App
	conf string
}

var (
	DefaultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			EnvVars: []string{"CONFIG"},
			Usage:   "config file",
			Value:   "/opt/config.yml",
		},
	}
)

func init() {
	rand.Seed(time.Now().Unix())
	help := cli.HelpPrinter
	cli.HelpPrinter = func(writer io.Writer, templ string, data interface{}) {
		help(writer, templ, data)
		os.Exit(0)
	}
}

func newCmd(opts ...Option) Cmd {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Description) == 0 {
		options.Description = "a stack-rpc service"
	}

	cmd := new(cmd)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = cmd.opts.Name
	cmd.app.Version = cmd.opts.Version
	cmd.app.Usage = cmd.opts.Description
	cmd.app.Flags = DefaultFlags
	cmd.app.Before = cmd.before
	cmd.app.Action = func(c *cli.Context) error { return nil }
	if len(options.Version) == 0 {
		cmd.app.HideVersion = true
	}

	return cmd
}

func (c *cmd) ConfigFile() string {
	return c.conf
}

func (c *cmd) before(ctx *cli.Context) error {
	// set the config file path
	if name := ctx.String("config"); len(name) > 0 {
		c.conf = name
	}
	return nil
}

func (c *cmd) App() *cli.App {
	return c.app
}

func (c *cmd) Options() Options {
	return c.opts
}

func (c *cmd) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	c.app.Name = c.opts.Name
	c.app.Version = c.opts.Version
	c.app.HideVersion = len(c.opts.Version) == 0
	c.app.Usage = c.opts.Description
	return c.app.Run(os.Args)
}

type Option func(o *Options)

type Options struct {
	// For the Command Line itself
	Name        string
	Description string
	Version     string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// endregion
