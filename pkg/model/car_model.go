package model

// Car ...a struct to represent a car:
type Car struct {
	Name           string  `json:"name" validate:"required"`
	MilesPerGallon float64 `json:"miles_per_gallon" validate:"required"`
	Cylinders      int     `json:"cylinders" validate:"required"`
	Displacement   float64 `json:"displacement" validate:"required"`
	Horsepower     int     `json:"horsepower" validate:"required"`
	WeightInLbs    int     `json:"weight_in_lbs" validate:"required"`
	Acceleration   float64 `json:"acceleration" validate:"required"`
	Year           string  `json:"year" validate:"required"`
	Origin         string  `json:"origin" validate:"required"`
}
