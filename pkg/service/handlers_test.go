package service_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	// "github.com/Blackmamba23/cars-fleet-service/internal/fleet/entities"
	"github.com/Blackmamba23/cars-fleet-service/internal/fleet/mocks"
	"github.com/Blackmamba23/cars-fleet-service/pkg/model"
	"github.com/Blackmamba23/cars-fleet-service/pkg/service"
	"github.com/NYTimes/gizmo/server"

	// "github.com/Blackmamba23/cars-fleet-service/internal/fleet/utils"

	"github.com/stretchr/testify/suite"
)

type fleetHandlerSuite struct {
	// we need this to use the suite functionalities from testify
	suite.Suite
	// the mocked version of the usecase
	fleetRepo *mocks.Repository
	// the functionalities we need to test
	service service.CarsFleetService
	// testing server to be used the handler
	testingServer *httptest.Server
}

func (s *fleetHandlerSuite) SetupSuite() {

	fmt.Println("From SetupSuite")

	fleetRepo := new(mocks.Repository)
	svc, _ := service.NewCarsFleetService(fleetRepo)

	ss := server.NewSimpleServer(nil)
	ss.Register(svc)

	// create and run the testing server
	testingServer := httptest.NewServer(ss)

	// assign the dependencies we need as the suite properties
	// we need this to run the tests
	s.testingServer = testingServer
	s.fleetRepo = fleetRepo
	s.service = *svc
}

func (s *fleetHandlerSuite) TearDownSuite() {
	fmt.Println("From TearDownSuite")

	defer s.testingServer.Close()
}

func (s *fleetHandlerSuite) TestGetCarByName_Positive() {
	fmt.Println("From GetCarByName_Positive")
	carNameVal := "mustang"
	// an example Car for the test
	car := model.Car{
		Name:           "ford mustang",
		MilesPerGallon: 18,
		Cylinders:      6,
		Displacement:   250,
		Horsepower:     88,
		WeightInLbs:    3139,
		Acceleration:   14.5,
		Year:           "1971-01-01T00:00:00Z",
		Origin:         "USA",
	}

	// fleetRepo's GetCarByName method will be called
	s.fleetRepo.On("GetCarByName", carNameVal).Return(car, nil)

	// calling the testing server given the provided request body
	response, err := http.Get(fmt.Sprintf("%s/svc/fleet/v1/car?car_name=%v", s.testingServer.URL, carNameVal))
	s.NoError(err, "no error when calling the endpoint")
	defer response.Body.Close()

	// unmarshalling the response
	responseBody := struct {
		Status  string    `json:"status"`
		Message string    `json:"message"`
		Data    model.Car `json:"data"`
	}{"Success!", "successfully fetched car", model.Car{}}
	json.NewDecoder(response.Body).Decode(&responseBody)

	// running assertions to make sure that our method does the correct thing
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(responseBody.Message, "successfully fetched car")
	s.Equal(responseBody.Data.Name, "ford mustang")
	s.fleetRepo.AssertExpectations(s.T())
}

func (s *fleetHandlerSuite) TestGetCarByName_Negative() {
	fmt.Println("From TestGetCarByName_Negative")
	carNameVal := "must"
	// an example Car for the test
	car := model.Car{}

	// fleetRepo's GetCarByName method will be called
	s.fleetRepo.On("GetCarByName", carNameVal).Return(car, errors.New("no results"))

	// calling the testing server given the provided request body
	response, err := http.Get(fmt.Sprintf("%s/svc/fleet/v1/car?car_name=%v", s.testingServer.URL, carNameVal))
	s.NoError(err, "no error when calling the endpoint")
	defer response.Body.Close()

	// unmarshalling the response
	responseBody := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{"Validation Error", fmt.Sprintf("could not fetch car named %v", carNameVal)}
	json.NewDecoder(response.Body).Decode(&responseBody)

	// running assertions to make sure that our method does the correct thing
	s.Equal(http.StatusBadRequest, response.StatusCode)
	s.Equal(responseBody.Message, fmt.Sprintf("could not fetch car named %v", carNameVal))
	s.Equal(car.Name, "")
	s.fleetRepo.AssertExpectations(s.T())
}

func (s *fleetHandlerSuite) TestGetCarsByName_Positive() {
	fmt.Println("From TestGetCarsByName_Positive")
	carNameVal := "mustang"
	// an example Car for the test
	car := []model.Car{
		{
			Name:           "ford mustang",
			MilesPerGallon: 18,
			Cylinders:      6,
			Displacement:   250,
			Horsepower:     88,
			WeightInLbs:    3139,
			Acceleration:   14.5,
			Year:           "1971-01-01T00:00:00Z",
			Origin:         "USA",
		},
		{
			Name:           "ford mustang ii",
			MilesPerGallon: 13,
			Cylinders:      6,
			Displacement:   250,
			Horsepower:     88,
			WeightInLbs:    3139,
			Acceleration:   14.5,
			Year:           "1971-01-01T00:00:00Z",
			Origin:         "USA",
		},
		{
			Name:           "ford mustang cobra",
			MilesPerGallon: 23.6,
			Cylinders:      6,
			Displacement:   250,
			Horsepower:     88,
			WeightInLbs:    3139,
			Acceleration:   14.5,
			Year:           "1971-01-01T00:00:00Z",
			Origin:         "USA",
		},
		{
			Name:           "ford mustang gl",
			MilesPerGallon: 27,
			Cylinders:      4,
			Displacement:   250,
			Horsepower:     88,
			WeightInLbs:    3139,
			Acceleration:   14.5,
			Year:           "1971-01-01T00:00:00Z",
			Origin:         "USA",
		},
	}

	// fleetRepo's GetCarByName method will be called
	s.fleetRepo.On("GetCarsByName", carNameVal).Return(car, nil)

	// calling the testing server given the provided request body
	response, err := http.Get(fmt.Sprintf("%s/svc/fleet/v1/cars?car_name=%v", s.testingServer.URL, carNameVal))
	s.NoError(err, "no error when calling the endpoint")
	defer response.Body.Close()

	// unmarshalling the response
	responseBody := struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    []model.Car `json:"data"`
	}{"Success!", "successfully fetched cars", []model.Car{}}
	json.NewDecoder(response.Body).Decode(&responseBody)

	// running assertions to make sure that our method does the correct thing
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(responseBody.Message, "successfully fetched cars")
	s.Equal(len(responseBody.Data), 4)
	s.fleetRepo.AssertExpectations(s.T())
}

func (s *fleetHandlerSuite) TestGetCarsByName_Negative() {
	fmt.Println("From TestGetCarsByName_Negative")
	carNameVal := "must"
	// an example Car for the test
	cars := []model.Car{}

	// fleetRepo's GetCarByName method will be called
	s.fleetRepo.On("GetCarsByName", carNameVal).Return(cars, errors.New("no results"))

	// calling the testing server given the provided request body
	response, err := http.Get(fmt.Sprintf("%s/svc/fleet/v1/cars?car_name=%v", s.testingServer.URL, carNameVal))
	s.NoError(err, "no error when calling the endpoint")
	defer response.Body.Close()

	// unmarshalling the response
	responseBody := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{"Validation Error", fmt.Sprintf("could not fetch cars with name %v", carNameVal)}
	json.NewDecoder(response.Body).Decode(&responseBody)

	// running assertions to make sure that our method does the correct thing
	s.Equal(http.StatusBadRequest, response.StatusCode)
	s.Equal(responseBody.Message, fmt.Sprintf("could not fetch cars with name %v", carNameVal))
	s.Equal(len(cars), 0)
	s.fleetRepo.AssertExpectations(s.T())
}

func TestFleetHandlerSuite(t *testing.T) {
	suite.Run(t, new(fleetHandlerSuite))
}
