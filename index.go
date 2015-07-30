package game

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"github.com/julienschmidt/httprouter"
)

func handleIndex(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	u := getUser(ctx)

	var indexData struct {
		LoginURL string
		Username string
	}
	l, err := user.LoginURL(ctx, "/")
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	indexData.LoginURL = l
	if u != nil {
		indexData.Username = u.Username
	}
	err = tpl.ExecuteTemplate(res, "index", indexData)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}
