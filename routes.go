package game

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := httprouter.New()
	r.GET("/", handleIndex)
	r.GET("/create", handleCreatePage)
	r.POST("/create", handleCreate)
	r.GET("/gamesList", handleGamesList)
	r.GET("/game/:gameID", handleGame)
	r.GET("/api/game/:gameName", handleGetState)
	http.Handle("/", r)
}
