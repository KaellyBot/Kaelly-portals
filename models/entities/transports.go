package entities

type TransportType struct {
	ID             string `gorm:"primaryKey"`
	DofusPortalsID string `gorm:"unique"`
}
