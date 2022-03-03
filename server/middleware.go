package server

import "github.com/gin-gonic/gin"

// time out
func TimeoutMiddleware(c *gin.Context) {

}

func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
}
