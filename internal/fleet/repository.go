package fleet

import (
	"errors"

	"github.com/Blackmamba23/cars-fleet-service/pkg/model"
	"github.com/blevesearch/bleve/v2"
)

// Data ...reposiory for the fleet that fetches data from a json file
type repository struct {
	searchIndex bleve.Index
}

// NewRepository ...init a new repository with the provided .json file
func NewRepository(searchIndex bleve.Index) Repository {
	return repository{
		searchIndex: searchIndex,
	}
}

func (r repository) GetCarByName(carName string) (model.Car, error) {
	car := model.Car{}
	// search for some text
	query := bleve.NewMatchQuery(carName)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{
		"name",
		"miles_per_gallon",
		"cylinders",
		"displacement",
		"horsepower",
		"weight_in_lbs",
		"acceleration",
		"year",
		"origin",
	}
	search.Size = 1
	searchResults, err := r.searchIndex.Search(search)
	if err != nil {
		return model.Car{}, err
	}

	if len(searchResults.Hits) == 0 {
		return car, errors.New("no results found")
	}

	car.Name = searchResults.Hits[0].Fields["name"].(string)
	car.MilesPerGallon = searchResults.Hits[0].Fields["miles_per_gallon"].(float64)
	car.Cylinders = int(searchResults.Hits[0].Fields["cylinders"].(float64))
	car.Displacement = searchResults.Hits[0].Fields["displacement"].(float64)
	car.Horsepower = int(searchResults.Hits[0].Fields["horsepower"].(float64))
	car.WeightInLbs = int(searchResults.Hits[0].Fields["weight_in_lbs"].(float64))
	car.Acceleration = searchResults.Hits[0].Fields["acceleration"].(float64)
	car.Year = searchResults.Hits[0].Fields["year"].(string)
	car.Origin = searchResults.Hits[0].Fields["origin"].(string)

	return car, nil
}

func (r repository) GetCarsByName(carName string) ([]model.Car, error) {
	cars := []model.Car{}
	// search for some text
	query := bleve.NewMatchQuery(carName)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{
		"name",
		"miles_per_gallon",
		"cylinders",
		"displacement",
		"horsepower",
		"weight_in_lbs",
		"acceleration",
		"year",
		"origin",
	}
	search.Size = 1000
	searchResults, err := r.searchIndex.Search(search)
	if err != nil {
		return cars, err
	}

	if len(searchResults.Hits) == 0 {
		return cars, errors.New("no results found")
	}

	for i := range searchResults.Hits {
		car := model.Car{}
		car.Name = searchResults.Hits[i].Fields["name"].(string)
		car.MilesPerGallon = searchResults.Hits[i].Fields["miles_per_gallon"].(float64)
		car.Cylinders = int(searchResults.Hits[i].Fields["cylinders"].(float64))
		car.Displacement = searchResults.Hits[i].Fields["displacement"].(float64)
		car.Horsepower = int(searchResults.Hits[i].Fields["horsepower"].(float64))
		car.WeightInLbs = int(searchResults.Hits[i].Fields["weight_in_lbs"].(float64))
		car.Acceleration = searchResults.Hits[i].Fields["acceleration"].(float64)
		car.Year = searchResults.Hits[i].Fields["year"].(string)
		car.Origin = searchResults.Hits[i].Fields["origin"].(string)

		cars = append(cars, car)
	}
	return cars, nil
}
