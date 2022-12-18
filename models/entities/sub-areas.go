package entities

type SubArea struct {
	Id             string `gorm:"primaryKey"`
	DofusPortalsId string `gorm:"unique"`
}
