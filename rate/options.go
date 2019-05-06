package rate

import (
	"golang.org/x/time/rate"
)

type WindowType int

const (
	WindowFixedSET WindowType = iota
	WindowFixedHSET
	WindowRolling
)

type Option func(o *Options)

type Options struct {
	window WindowType

	addr     string
	password string
	poolSize int

	fallback *rate.Limiter
}

func DefaultOptions() *Options {
	return &Options{
		window: WindowFixedSET,
		addr:   "127.0.0.1:6379",
	}
}

func Window(windowType WindowType) Option {
	return func(o *Options) {
		o.window = windowType
	}
}

func Addr(addr string) Option {
	return func(o *Options) {
		o.addr = addr
	}
}

func Password(pw string) Option {
	return func(o *Options) {
		o.password = pw
	}
}

func PoolSize(size int) Option {
	return func(o *Options) {
		o.poolSize = size
	}
}

func Fallback(f *rate.Limiter) Option {
	return func(o *Options) {
		o.fallback = f
	}
}
