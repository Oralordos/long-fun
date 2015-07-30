package game

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/julienschmidt/httprouter"
)

type game struct {
	Map   gameMap
	Units int
	Name  string
}

func (l *game) Load(ps []datastore.Property) error {
	bs := ps[0].Value.([]byte)
	err := json.Unmarshal(bs, l)
	return err
}

func (l *game) Save() ([]datastore.Property, error) {
	// TODO Only convert map data with json
	bs, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return []datastore.Property{
		datastore.Property{
			Name:    "Layers",
			Value:   bs,
			NoIndex: true,
		},
	}, nil
}

type gameID struct {
	ID   int64
	Name string
}

func addGame(ctx context.Context, g *game) (int64, error) {
	k := datastore.NewIncompleteKey(ctx, "Game", nil)
	ck, err := datastore.Put(ctx, k, g)
	if err != nil {
		return 0, err
	}
	return ck.IntID(), err
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
