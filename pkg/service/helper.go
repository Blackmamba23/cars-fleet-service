package service

// ApplicationConfig ...sets the application ENV varibles
type ApplicationConfig struct {
	ApplicationEnv   string `required:"true" default:"DEV"`
	JsonFileDataPath string `required:"true" envconfig:"JSON_FILE_DATA_PATH" default:"cars.json"`
}
