package game

import (
	"encoding/json"
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/julienschmidt/httprouter"
)

type gameMap struct {
	Height int
	Width  int
	Layers []struct {
		Name string
		Data []int32
	}
	Tileheight int
	Tilewidth  int
	Tilesets   []struct {
		Imagewidth  int
		Imageheight int
	}
}

func loadMap(filename string) (*gameMap, error) {
	f, err := os.Open("maps/" + filename + ".json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var mapData gameMap
	err = json.NewDecoder(f).Decode(&mapData)
	if err != nil {
		return nil, err
	}
	return &mapData, nil
}

func handleMap(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(req)
	mapName := p.ByName("mapName")
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
