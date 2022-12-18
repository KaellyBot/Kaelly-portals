package entities

type Area struct {
	Id             string `gorm:"primaryKey"`
	DofusPortalsId string `gorm:"unique"`
}
