package str

import "github.com/google/uuid"

func GenerateUUIDName() string {
	return uuid.New().String()
}