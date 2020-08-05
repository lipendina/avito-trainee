package handlers

import "net/http"

func getErrorStatus(isInternal bool) int {
	if !isInternal {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
