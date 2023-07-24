package config

import (
	"fmt"
	"os"
	"strconv"
)

const errorMessage = "Can't convert env var '%s' with value '%s' to %s"

func GetEnvStr(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

func GetEnvInt(key string, defaultValue int) int {
	strVal := GetEnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "int"))
	}

	return intVal
}

func GetEnvBool(key string, defaultValue bool) bool {
	strVal := GetEnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "bool"))
	}

	return boolVal
}
