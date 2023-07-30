package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"check-in/api/internal/constants"
)

func ReadUUIDURLParam(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return "", err
	}

	value := id.String()
	return value, nil
}

func ReadIntURLParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func ReadStrQueryParam(r *http.Request, name string, defaultValue string) string {
	param := r.URL.Query().Get(name)

	if param == "" {
		return defaultValue
	}

	return param
}

func ReadIntQueryParam(
	r *http.Request,
	name string,
	defaultValue int64,
) (int64, error) {
	param := r.URL.Query().Get(name)

	if param == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("invalid %s query param", name)
	}

	return value, nil
}

func ReadDateQueryParam(
	r *http.Request,
	name string,
	defaultValue *time.Time,
) (*time.Time, error) {
	param := r.URL.Query().Get(name)

	if param == "" {
		return defaultValue, nil
	}

	value, err := time.Parse(constants.DateFormatTz, param)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid %s query param, needs timezone (yyyy-MM-ddx)",
			name,
		)
	}

	return &value, nil
}