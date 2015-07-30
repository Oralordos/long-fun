package game

import (
	"net/http"
	"strconv"

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
	ID   int64
	Name string
}

func handleGame(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var gameData struct {
		GameID int64
	}
	id, err := strconv.ParseInt(p.ByName("gameID"), 10, 64)
	gameData.GameID = id
	err = tpl.ExecuteTemplate(res, "game", gameData)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		ctx := appengine.NewContext(req)
		log.Errorf(ctx, err.Error())
		return
	}
}
