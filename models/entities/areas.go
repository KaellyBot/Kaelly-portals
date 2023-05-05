package entities

type Area struct {
	ID             string `gorm:"primaryKey"`
	DofusPortalsID string `gorm:"unique"`
}
