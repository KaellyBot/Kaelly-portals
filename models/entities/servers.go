package entities

type Server struct {
	Id             string `gorm:"primaryKey"`
	DofusPortalsId string `gorm:"unique"`
}
