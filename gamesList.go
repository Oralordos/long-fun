package game

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/julienschmidt/httprouter"
)

func handleGamesList(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	u := getUser(ctx)
	if u == nil {
		http.Error(res, "Not logged in", http.StatusUnauthorized)
		return
	}
	var gamesListData struct {
		Username string
		Games    []gameID
	}
	gamesListData.Username = u.Username
	gamesListData.Games = u.Games
	err := tpl.ExecuteTemplate(res, "games-list", gamesListData)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}
