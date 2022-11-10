package util

import (
	"os"
	"path"

	"github.com/rs/zerolog"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const pathToLogConfig = "configuration/livesettings.json"
const logLevel = "Logging.LogLevel.Default"

// GetConfigValues gets constant values from the file and injects them
func GetConfigValues(pathToConfig string) (map[string]string, error) {
	currentPath, _ := os.Getwd()
	fullPath := path.Join(currentPath, pathToConfig)
	viper.SetConfigFile(fullPath)
	viper.SetConfigType("json")
	log.Info().Msgf("Path: %s", fullPath)
	err := viper.ReadInConfig() // Find and read the config file
	// just use the default value(s) if the config file was not found
	if _, ok := err.(*os.PathError); ok {
		log.Warn().Msgf("No config file '%s' not found. Using default values", fullPath)
	} else if err != nil { // Handle other errors that occurred while reading the config file
		log.Err(err).Msgf("Error while reading the config file")
	}
	var configuration map[string]string
	err = viper.Unmarshal(&configuration)
	return configuration, err
}

// SetLogLevels gets configuration values from the file and injects them
func SetLogLevels() {
	currentPath, _ := os.Getwd()
	fullPath := path.Join(currentPath, pathToLogConfig)
	viper.SetConfigFile(fullPath)
	viper.SetConfigType("json")
	err := viper.ReadInConfig() // Find and read the config file
	// just use the default value(s) if the config file was not found
	if _, ok := err.(*os.PathError); ok {
		log.Warn().Msgf("No config file '%s' not found. Using default values", fullPath)
	} else if err != nil { // Handle other errors that occurred while reading the config file
		log.Err(err).Msgf("Error while reading the config file")
	} else {
		log.Info().Msgf("Log Level from config: %s", viper.GetString(logLevel))
		setLogLevel(viper.GetString(logLevel))
		// monitor the changes in the config file
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Info().Msgf("Log Level from config: %s", viper.GetString(logLevel))
			setLogLevel(viper.GetString(logLevel))
		})
	}
}

func setLogLevel(level string) {
	switch level {
	case "Debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		break
	case "Info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break
	case "Warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		break
	case "Error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		break
	default:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}
}
