package services

import (
	"layersapi/entities"
	"time"
)

func CreateMetadata() entities.Metadata {
	now := time.Now().Format(time.RFC3339)

	return entities.Metadata{
		CreatedAt: now,
		CreatedBy: "webapp",
		UpdatedAt: now,
		UpdatedBy: "webapp",
	}
}
