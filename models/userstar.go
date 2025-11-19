package models

type UserStar struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	UserID  uint   `json:"userid"`
	Type    string `json:"type"`
	AssetID uint   `json:"assetid"`
}
