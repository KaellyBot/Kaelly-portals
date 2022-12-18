package entities

type TransportType struct {
	Id             string `gorm:"primaryKey"`
	DofusPortalsId string `gorm:"unique"`
}
