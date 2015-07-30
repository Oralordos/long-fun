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
	Units []unit
	Name  string
}

type datastoreGame struct {
	MapWidth   int     `datastore:",noindex"`
	MapHeight  int     `datastore:",noindex"`
	MapTileset tileset `datastore:",noindex"`
	MapLayers  []byte  `datastore:",noindex"`
	Units      []unit  `datastore:",noindex"`
	Name       string
}

func (g *game) Load(ps []datastore.Property) error {
	var d datastoreGame
	err := datastore.LoadStruct(&d, ps)
	if err != nil {
		return err
	}
	g.Name = d.Name
	g.Map.Width = d.MapWidth
	g.Map.Height = d.MapHeight
	g.Map.Tileset = d.MapTileset
	g.Units = d.Units
	var l mapLayers
	err = json.Unmarshal(d.MapLayers, &l)
	if err != nil {
		return err
	}
	g.Map.Layers = l
	return nil
}

func (g *game) Save() ([]datastore.Property, error) {
	d := datastoreGame{
		MapWidth:   g.Map.Width,
		MapHeight:  g.Map.Height,
		MapTileset: g.Map.Tileset,
		Name:       g.Name,
		Units:      g.Units,
	}
	bs, err := json.Marshal(g.Map.Layers)
	if err != nil {
		return nil, err
	}
	d.MapLayers = bs
	return datastore.SaveStruct(&d)
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
