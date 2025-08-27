package utils

import (
	"errors"

	"github.com/gofrs/uuid"
)

func CreateUUID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", errors.Join(err)
	}

	return id.String(), nil
}
