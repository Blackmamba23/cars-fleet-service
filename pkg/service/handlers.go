package service

import (
	"fmt"
	"net/http"

	"github.com/Blackmamba23/cars-fleet-service/pkg/model"
)

func (s *CarsFleetService) GetCarByName(r *http.Request) (int, interface{}, error) {
	r.ParseForm()
	if len(r.FormValue("car_name")) == 0 {
		return http.StatusBadRequest, struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{"Validation Error", "car name is required"}, nil
	}

	car, err := s.fleetRepo.GetCarByName(r.FormValue("car_name"))
	if err != nil {
		return http.StatusBadRequest, struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{"Validation Error", fmt.Sprintf("could not fetch car named %v", r.FormValue("car_name"))}, nil
	}

	return http.StatusOK, struct {
		Status  string    `json:"status"`
		Message string    `json:"message"`
		Data    model.Car `json:"data"`
	}{"Success!", "successfully fetched car", car}, nil
}

func (s *CarsFleetService) GetCarsByName(r *http.Request) (int, interface{}, error) {
	r.ParseForm()
	if len(r.FormValue("car_name")) == 0 {
		return http.StatusBadRequest, struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{"Validation Error", "car name is required"}, nil
	}

	cars, err := s.fleetRepo.GetCarsByName(r.FormValue("car_name"))
	if err != nil {
		return http.StatusBadRequest, struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}{"Validation Error", fmt.Sprintf("could not fetch cars with name %v", r.FormValue("car_name"))}, nil
	}

	return http.StatusOK, struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    []model.Car `json:"data"`
	}{"Success!", "successfully fetched cars", cars}, nil
}
