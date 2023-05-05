package entities

type Server struct {
	ID             string `gorm:"primaryKey"`
	DofusPortalsID string `gorm:"unique"`
}
