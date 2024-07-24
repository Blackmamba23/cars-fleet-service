package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Blackmamba23/cars-fleet-service/internal/fleet"
	"github.com/Blackmamba23/cars-fleet-service/pkg/model"
	"github.com/Blackmamba23/cars-fleet-service/pkg/service"
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/blevesearch/bleve/v2"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var cfg struct {
		Server *server.Config
	}

	// set up application ENVs
	var appEnvs service.ApplicationConfig
	err := envconfig.Process("cars-fleet-service", &appEnvs)
	if err != nil {
		log.Fatal(err.Error())
	}

	// load from the local JSON file into a config.Config struct
	config.LoadJSONFile("./config.json", &cfg)

	// NOTE: We can easily change out the server for another eg gin, echo, mux, martini

	// SetConfigOverrides will allow us to override some of the values in
	// the JSON file with CLI flags.
	server.SetConfigOverrides(cfg.Server)

	// initialize server with given configs
	server.Init("cars-fleet-service", cfg.Server)

	index, err := bleve.Open("cars.bleve")
	if err == bleve.ErrorIndexPathDoesNotExist {
		// initialize index
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New("cars.bleve", mapping)
		if err != nil {
			server.Log.Fatal(err.Error())
		}

		// index data in another thread in the background
		go func() {
			// read the json file to be indexed
			file, err := ioutil.ReadFile(appEnvs.JsonFileDataPath)
			if err != nil {
				server.Log.Fatal(err.Error())
			}

			data := []model.Car{}

			err = json.Unmarshal([]byte(file), &data)
			if err != nil {
				server.Log.Fatal(err.Error())
			}
			count := 0
			for i := 0; i < len(data); i++ {
				count++
				// index the data
				index.Index(data[i].Name, data[i])
			}
			server.Log.Info(fmt.Sprintf("indexed %d documents", count))
		}()

	} else if err != nil {
		server.Log.Fatal(err)
	} else {
		server.Log.Info("opening existing index...")
	}

	// instantiate a new fleet service
	svc, err := service.NewCarsFleetService(fleet.NewRepository(index))
	if err != nil {
		server.Log.Fatal("unable to create cars fleet service: ", err)
	}

	// register our saved item service with the Gizmo server
	err = server.Register(svc)
	if err != nil {
		server.Log.Fatal("unable to register cars fleet service: ", err)
	}

	// run the Gizmo server
	err = server.Run()
	if err != nil {
		server.Log.Fatal("unable to run cars fleet service: ", err)
	}
}
