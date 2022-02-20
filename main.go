package main

import (
	"fire-scaffold/conf"
	"fire-scaffold/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf.InitConfig("./conf/dev-env.yaml")

	server := server.New(
		conf.GlobalConfig.HTTP.GinMode,
		conf.GlobalConfig.HTTP.Addr,
		server.SetReadTimeout(conf.GlobalConfig.HTTP.ReadTimeout),
		server.SetWriteTimeout(conf.GlobalConfig.HTTP.WriteTimeout),
		server.SetMaxHeaderBytes(1<<uint(conf.GlobalConfig.HTTP.MaxHeaderBytes)),
	)

	server.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.Stop()
}
