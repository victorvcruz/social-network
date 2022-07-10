package utils

import (
	"os"
	"strconv"
)

func GetStringEnvOrElse(envName string, defaultValue string) string {
	value, found := os.LookupEnv(envName)
	if !found {
		value = defaultValue
	}
	return value
}

func GetIntEnvOrElse(envName string, defaultValue int) (value int) {
	valueStr, found := os.LookupEnv(envName)
	if !found {
		value = defaultValue
	} else {
		value, _ = strconv.Atoi(valueStr)
	}
	return value
}
