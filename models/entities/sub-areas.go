package entities

type SubArea struct {
	ID             string `gorm:"primaryKey"`
	DofusPortalsID string `gorm:"unique"`
}
