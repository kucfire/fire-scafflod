package main

import (
	"fmt"
	"log"

	"fire-scaffold/cache"
	"fire-scaffold/conf"
	"fire-scaffold/pkg/shutdown"
	"fire-scaffold/server"
)

func main() {
	conf.InitConfig("./conf/dev-env.yaml")

	if err := cache.InitRedis(conf.GlobalConfig.Redis); err != nil {
		panic(fmt.Sprintf("init redis error: %v", err))
	}

	server := server.New(
		conf.GlobalConfig.HTTP.GinMode,
		conf.GlobalConfig.HTTP.Addr,
		server.SetReadTimeout(conf.GlobalConfig.HTTP.ReadTimeout),
		server.SetWriteTimeout(conf.GlobalConfig.HTTP.WriteTimeout),
		server.SetMaxHeaderBytes(1<<uint(conf.GlobalConfig.HTTP.MaxHeaderBytes)),
	)

	server.Run()

	// graceful close 优雅关闭
	shutdown.New().Close(
		func() {
			if err := cache.Close(); err != nil {
				log.Fatalf(" [ERROR] cache close:%s err:%v\n", conf.GlobalConfig.Redis.Addr, err)
			}
		},
		func() {
			server.Stop()
		},
	)
}
