package entity

import "github.com/google/uuid"

type ID = uuid.UUID

func NewUUID() ID {
	return ID(uuid.New())
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return id, err
}
