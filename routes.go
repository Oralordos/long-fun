package game

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := httprouter.New()
	r.GET("/", handleGame)
	r.GET("/games/:gameName", handleGetState)
	http.Handle("/", r)
}
