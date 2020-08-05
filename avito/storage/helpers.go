package storage

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func makeParamsFromUUID(paramIDs []uuid.UUID) (string, []interface{}) {
	params := make([]string, 0, len(paramIDs))
	result := make([]interface{}, 0, len(paramIDs))
	for idx, id := range paramIDs {
		params = append(params, fmt.Sprintf("$%d", idx+1))
		result = append(result, id)
	}

	return strings.Join(params, ","), result
}
