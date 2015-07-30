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

type datastoreGame struct {
	MapWidth             int
	MapHeight            int
	MapTilesetTilewidth  int
	MapTilesetTileheight int
	MapTilesetWidth      int
	MapTilesetHeight     int
	MapTilesetFilename   string
	MapLayers            []byte
	Units                int
	Name                 string
}

func (g *game) Load(ps []datastore.Property) error {
	var d datastoreGame
	err := datastore.LoadStruct(&d, ps)
	if err != nil {
		return err
	}
	g.Name = d.Name
	g.Units = d.Units
	g.Map.Width = d.MapWidth
	g.Map.Height = d.MapHeight
	g.Map.Tileset.Width = d.MapTilesetWidth
	g.Map.Tileset.Height = d.MapTilesetHeight
	g.Map.Tileset.Tilewidth = d.MapTilesetTilewidth
	g.Map.Tileset.Tileheight = d.MapTilesetTileheight
	g.Map.Tileset.Filename = d.MapTilesetFilename
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
		MapWidth:             g.Map.Width,
		MapHeight:            g.Map.Height,
		MapTilesetTilewidth:  g.Map.Tileset.Tilewidth,
		MapTilesetTileheight: g.Map.Tileset.Tileheight,
		MapTilesetWidth:      g.Map.Tileset.Width,
		MapTilesetHeight:     g.Map.Tileset.Height,
		MapTilesetFilename:   g.Map.Tileset.Filename,
		Units:                g.Units,
		Name:                 g.Name,
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
