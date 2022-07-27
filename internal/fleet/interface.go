package fleet

import (
	"github.com/Blackmamba23/cars-fleet-service/pkg/model"
)

// Repository Interface
type Repository interface {
	GetCarByName(carName string) (model.Car, error)
	GetCarsByName(carName string) ([]model.Car, error)
}
