package model

import "time"

type ShareToken struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Token     string     `json:"token" gorm:"uniqueIndex;size:64"`
	Path      string     `json:"path"`        // e.g. "/MyFolder/ProjectX"
	Label     string     `json:"label"`       // human-readable name
	ExpiresAt *time.Time `json:"expires_at"` // nil = never expires
	CreatedAt time.Time  `json:"created_at"`
}
