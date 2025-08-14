package utils

import (
	"strings"

	"github.com/gofrs/uuid"
)

func UUID() string {
	return strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
}
