package game

import (
	"encoding/json"
	"net/http"
	"os"

	"google.golang.org/appengine"
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

type gameMap struct {
	Width   int
	Height  int
	Layers  [][]int
	Tileset tileset
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
	gameLayers := [][]int{}
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

func handleGetState(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(req)
	mapName := p.ByName("gameName")
	// TODO Get the map name from the datastore
	gameMap, err := loadMap(mapName)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	err = json.NewEncoder(res).Encode(gameMap)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}
