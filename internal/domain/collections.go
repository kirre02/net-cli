package domain

import (
	"errors"
	"strings"
)

var ErrEmptyCollectionName = errors.New("collection name cannot be empty")

type Collection struct {
	ID   string
	Name string
}

func (c Collection) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrEmptyCollectionName
	}
	return nil
}
