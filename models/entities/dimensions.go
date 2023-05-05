package entities

type Dimension struct {
	ID             string `gorm:"primaryKey"`
	DofusPortalsID string `gorm:"unique"`
}
