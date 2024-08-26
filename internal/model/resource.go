package model

import "gorm.io/gorm"

// Resource represents resource data.

type Resource struct {
	gorm.Model
	ResourceType string
	Payload      string
}
