# pkg
Go common package.

- Feature
	- `conv`
		- `struct` 使用reflect实现struct同名字段转换，支持Struct、Ptr迭代
    - `dispatcher` 协程调度器
    - [ ] `lock`
    - `log` 日志
        - `zap`
    - `pool` 协程池
    - `rate` 分布式限流 [分布式系统限流服务-Golang&Redis](http://hbchen.com/post/distributed/2019-05-05-rate-limit/)
        - `fixed window` 固定窗口
        - `rolling window` 滚动窗口
    - `sync`
        - [ ] `watcher`
