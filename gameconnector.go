package gameconnector

import (
	"github.com/coolaxie/gameconnector/gate"
	"github.com/coolaxie/gameconnector/router"
)

type GameConnector struct {
	gate   *gate.Gate
	router *router.Router
}
