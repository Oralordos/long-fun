package game

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := httprouter.New()
	r.GET("/", handleIndex)
	r.GET("/create/user", handleCreateUserPage)
	r.POST("/create/user", handleCreateUser)
	r.GET("/create/game", handleCreateGamePage)
	r.POST("/create/game", handleCreateGame)
	r.GET("/gamesList", handleGamesList)
	r.GET("/game/:gameID", handleGame)
	r.GET("/api/game/:gameName", handleGetState)
	http.Handle("/", r)
}
