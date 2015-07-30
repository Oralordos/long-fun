package game

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"github.com/julienschmidt/httprouter"
)

func handleCreate(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

func handleCreatePage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	u := getUser(ctx)
	if u != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(res, "create", nil)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}
