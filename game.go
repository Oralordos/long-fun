package game

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/julienschmidt/httprouter"
)

type game struct {
	Map   gameMap
	Units int
	Name  string
}

type gameID struct {
	ID   int
	Name string
}

func handleGame(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var gameData struct {
		GameID string
	}
	gameData.GameID = "test"
	err := tpl.ExecuteTemplate(res, "game", gameData)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		ctx := appengine.NewContext(req)
		log.Errorf(ctx, err.Error())
		return
	}
}
