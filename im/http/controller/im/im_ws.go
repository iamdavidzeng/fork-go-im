package im

import (
	"fork_go_im/im/service"
	"fork_go_im/pkg/jwt"
	"fork_go_im/pkg/pool"
	"fork_go_im/pkg/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IMService struct{}

func (*IMService) Connect(c *gin.Context) {
	conn, err := ws.App(c.Writer, c.Request)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	client := &service.ImClient{ID: claims.ID, Socket: conn, Send: make(chan []byte)}
	service.ImManager.Register <- client

	pool.AntsPool.Submit(func() {
		client.ImRead()
	})

	pool.AntsPool.Submit(func() {
		client.ImWrite()
	})
}
