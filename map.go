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
		Image          string
		Imagewidth     int
		Imageheight    int
		Tileproperties map[string]map[string]string
	}
}

type tileset struct {
	Tilewidth       int
	Tileheight      int
	Width           int
	Height          int
	RedTeam         []int64
	BlueTeam        []int64
	YellowTeam      []int64
	GreenTeam       []int64
	Filename        string
	MoveIndicator   int64
	AttackIndicator int64
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
			RedTeam:    []int64{},
			BlueTeam:   []int64{},
			GreenTeam:  []int64{},
			YellowTeam: []int64{},
		},
	}
	gameLayers := mapLayers{}
	units := []unit{}
	for _, l := range jsonMapData.Layers {
		if l.Name == "Unit Layer" {
			for x := 0; x < jsonMapData.Width; x++ {
				for y := 0; y < jsonMapData.Height; y++ {
					ct := l.Data[x+y*jsonMapData.Width]
					if ct != 0 {
						newUnit := newUnit(ct, x, y)
						units = append(units, *newUnit)
					}
				}
			}
		} else {
			gameLayers = append(gameLayers, l.Data)
		}
	}
	mapData.Layers = gameLayers
	for k, v := range jsonMapData.Tilesets[0].Tileproperties {
		key, err := strconv.ParseInt(k, 10, 64)
		key++
		if err != nil {
			return nil, err
		}
		for prop, propValue := range v {
			switch prop {
			case "team":
				switch propValue {
				case "red":
					mapData.Tileset.RedTeam = append(mapData.Tileset.RedTeam, key)
				case "blue":
					mapData.Tileset.BlueTeam = append(mapData.Tileset.BlueTeam, key)
				case "green":
					mapData.Tileset.GreenTeam = append(mapData.Tileset.GreenTeam, key)
				case "yellow":
					mapData.Tileset.YellowTeam = append(mapData.Tileset.YellowTeam, key)
				}
			case "indicator":
				switch propValue {
				case "move":
					mapData.Tileset.MoveIndicator = key
				case "attack":
					mapData.Tileset.AttackIndicator = key
				}
			}
		}
	}
	g := &game{
		Map:   mapData,
		Units: units,
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
	var gameMap *game
	if gameID == 1 {
		gameMap, err = loadMap("test")
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
	} else {
		gameMap, err = getGame(ctx, gameID)
		if err != nil {
			http.NotFound(res, req)
			log.Warningf(ctx, err.Error())
			return
		}
	}
	err = json.NewEncoder(res).Encode(gameMap)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
}
