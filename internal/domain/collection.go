package domain

import "errors"

var ErrEmptyCollectionName = errors.New("collection name cannot be empty")

type Collection struct {
	ID   string
	Name string
}

func (c Collection) Validate() error {
	if c.Name == "" {
		return ErrEmptyCollectionName
	}
	return nil
}
