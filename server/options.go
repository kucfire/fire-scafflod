package server

import "time"

type Option func(*HttpServer)

func SetMaxHeaderBytes(param int) Option {
	return func(hs *HttpServer) {
		hs.Server.MaxHeaderBytes = param
	}
}

func SetReadTimeout(param time.Duration) Option {
	return func(hs *HttpServer) {
		hs.Server.ReadTimeout = param
	}
}

func SetWriteTimeout(param time.Duration) Option {
	return func(hs *HttpServer) {
		hs.Server.WriteTimeout = param
	}
}
