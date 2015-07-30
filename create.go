package game

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"github.com/julienschmidt/httprouter"
)

func handleCreateUser(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)
	username := req.FormValue("username")
	err := createUser(ctx, username, u.Email)
	// TODO Check for specific errors
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func handleCreateUserPage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	u := getUser(ctx)
	if u != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(res, "create-user", nil)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}

func handleCreateGame(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	u := getUser(ctx)
	if u == nil {
		http.Error(res, "Must be logged in", http.StatusUnauthorized)
		return
	}
	name, mapName := req.FormValue("name"), req.FormValue("map")
	g, err := loadMap(mapName)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	id, err := addGame(ctx, g)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	gInfo := &gameID{
		Name: name,
		ID:   id,
	}
	err = addGameToProfile(ctx, gInfo)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	http.Redirect(res, req, fmt.Sprintf("/game/%d", id), http.StatusSeeOther)
}

func handleCreateGamePage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	maps, err := getMaps()
	if err != nil {
		ctx := appengine.NewContext(req)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	var createGameData struct {
		Maps []string
	}
	createGameData.Maps = maps
	err = tpl.ExecuteTemplate(res, "create-game", createGameData)
	if err != nil {
		ctx := appengine.NewContext(req)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}
