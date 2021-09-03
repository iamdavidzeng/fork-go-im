package pool

import (
	"fork_go_im/pkg/config"

	"github.com/panjf2000/ants/v2"
)

var AntsPool *ants.Pool

func ConnectPool() *ants.Pool {
	AntsPool, _ = ants.NewPool(config.GetInt("app.go_corotines"))
	return AntsPool
}
