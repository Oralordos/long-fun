package game

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/julienschmidt/httprouter"
)

type jsonGameMap struct {
	Height int
	Width  int
	Layers []struct {
		Name string
		Data []int
	}
	Tileheight int
	Tilewidth  int
	Tilesets   []struct {
		Image       string
		Imagewidth  int
		Imageheight int
	}
}

type tileset struct {
	Tilewidth  int
	Tileheight int
	Width      int
	Height     int
	Filename   string
}

type mapLayers [][]int

type gameMap struct {
	Width   int
	Height  int
	Layers  mapLayers
	Tileset tileset
}

func getMaps() ([]string, error) {
	var maps []string
	err := filepath.Walk("maps", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".json") {
			name := info.Name()
			maps = append(maps, name[:len(name)-5])
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return maps, nil
}

func loadMap(filename string) (*game, error) {
	f, err := os.Open("maps/" + filename + ".json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var jsonMapData jsonGameMap
	err = json.NewDecoder(f).Decode(&jsonMapData)
	if err != nil {
		return nil, err
	}
	mapData := gameMap{
		Width:  jsonMapData.Width,
		Height: jsonMapData.Height,
		Tileset: tileset{
			Width:      jsonMapData.Tilesets[0].Imagewidth / jsonMapData.Tilewidth,
			Height:     jsonMapData.Tilesets[0].Imageheight / jsonMapData.Tileheight,
			Tilewidth:  jsonMapData.Tilewidth,
			Tileheight: jsonMapData.Tileheight,
			Filename:   jsonMapData.Tilesets[0].Image[2:],
		},
	}
	gameLayers := mapLayers{}
	for _, l := range jsonMapData.Layers {
		if l.Name == "Unit Layer" {
			// TODO Load the units.
		} else {
			gameLayers = append(gameLayers, l.Data)
		}
	}
	mapData.Layers = gameLayers
	g := &game{
		Map: mapData,
	}
	return g, nil
}

func getGame(ctx context.Context, gameID int64) (*game, error) {
	k := datastore.NewKey(ctx, "Game", "", gameID, nil)
	var g game
	err := datastore.Get(ctx, k, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func handleGetState(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(req)
	gn := p.ByName("gameName")
	gameID, err := strconv.ParseInt(gn, 10, 64)
	if err != nil {
		http.NotFound(res, req)
		return
	}
	gameMap, err := getGame(ctx, gameID)
	if err != nil {
		http.NotFound(res, req)
		log.Warningf(ctx, err.Error())
		return
	}
	err = json.NewEncoder(res).Encode(gameMap)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}
