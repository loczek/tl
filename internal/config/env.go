package config

import "os"

func IsProd() bool {
	return os.Getenv("ENV") != "development"
}

func IsDev() bool {
	return !IsProd()
}
