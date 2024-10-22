package api

import (
	"fmt"
	"net/http"
	"strconv"

	"kiln-exercice/pkg/api"
)

// PaginationFromRequest takes a request and returns a sorting Pagination.
func PaginationFromRequest(req *http.Request) (pagination api.Pagination, err error) {
	pagination.PageNumber, err = getPaginationValueFromRequest(req, api.PageNumberKey, api.DefaultPageNumber)
	if err != nil {
		return api.Pagination{}, err
	}

	pagination.PageSize, err = getPaginationValueFromRequest(req, api.PageSizeKey, api.DefaultPageSize)
	if err != nil {
		return api.Pagination{}, err
	}

	return pagination, nil
}

func getPaginationValueFromRequest(req *http.Request, name string, defaultValue int) (int, error) {
	var value int

	var err error

	if rawVal := req.URL.Query().Get(name); rawVal == "" {
		value = defaultValue
	} else {
		value, err = strconv.Atoi(rawVal)
		if err != nil {
			return value, fmt.Errorf("strconv.Atoi(%s): %w", rawVal, err)
		}

		if value < 1 {
			return value, fmt.Errorf("%s must be greater than 0", name)
		}
	}

	return value, nil
}
