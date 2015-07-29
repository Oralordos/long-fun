package game

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	r := httprouter.New()
	r.GET("/", handleGame)
	r.GET("/maps/:mapName", handleMap)
	http.Handle("/", r)
}
