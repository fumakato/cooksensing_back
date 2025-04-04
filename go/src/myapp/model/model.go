package model

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	// "github.com/jinzhu/gorm"
)

/*
データベース一覧
User
FeatureData
BestData
Histogram
Action
DisplayItem
*/

type BaseModel struct {
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BaseModel
	TsukurepoID              string `gorm:"size:255" json:"tsukurepo_id"`
	Name                     string `gorm:"size:255;not null" json:"name"`
	FirebaseAuthUID          string `gorm:"size:255;unique;autoIncrement" json:"firebase_auth_uid"`
	IsIdentification         bool   `gorm:"default:false" json:"is_identification"`
	IdentityVerificationText string `gorm:"size:255" json:"identity_verification_text"`
	// Email                    string `gorm:"size:255;unique;" json:"email"`
	// Password                 string `gorm:"size:255" json:"password"`
}

type FeatureData struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BaseModel
	UserID                        uint      `gorm:"not null" json:"user_id"`
	ActionID                      uint      `json:"action_id"`
	Date                          time.Time `json:"date"`
	AveragePace                   float32   `gorm:"type:double" json:"average_pace"`
	AccelerationStandardDeviation float32   `gorm:"type:double" json:"acceleration_standard_deviation"`
}

type Return struct {
	Date        []string  `json:"date"`
	AveragePace []float32 `json:"average_pace"`
}

type BestData struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	BaseModel
	AveragePace                   float32 `gorm:"type:double" json:"average_pace"`
	AccelerationStandardDeviation float32 `gorm:"type:double" json:"acceleration_standard_deviation"`
	// AveragePaceFeatureDataID                   uint    `json:"average_pace_feature_data_id"`
	// AccelerationStandardDeviationFeatureDataID uint    `json:"acceleration_standard_deviation_feature_data_id"`
}

type Histogram struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BaseModel
	DisplayItemID uint    `json:"display_item_id"`
	ActionID      uint    `json:"action_id"`
	Time1         uint    `gorm:"default:0" json:"time1"`
	Time2         uint    `gorm:"default:0" json:"time2"`
	Time3         uint    `gorm:"default:0" json:"time3"`
	Time4         uint    `gorm:"default:0" json:"time4"`
	Time5         uint    `gorm:"default:0" json:"time5"`
	Time6         uint    `gorm:"default:0" json:"time6"`
	Time7         uint    `gorm:"default:0" json:"time7"`
	Time8         uint    `gorm:"default:0" json:"time8"`
	Time9         uint    `gorm:"default:0" json:"time9"`
	Time10        uint    `gorm:"default:0" json:"time10"`
	Max           float32 `json:"max"`
	Min           float32 `json:"min"`
	Range         float32 `json:"range"`
}

// 輪切り　とか
type Action struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BaseModel
	Type string `gorm:"size:255" json:"type"`
}

// 平均ペース　とか
type DisplayItem struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BaseModel
	Item string `gorm:"size:255" json:"item"`
}

// type Label struct {
// 	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
// 	BaseModel
// 	ActionID    uint   `json:"action_id"`
// 	DisplayItem string `gorm:"size:255" json:"display_item"`
// 	MainTitle   string `gorm:"size:255" json:"main_title"`
// 	SubTitleX   string `gorm:"size:255" json:"sub_title_x"`
// 	SubtitleY   string `gorm:"size:255" json:"sub_title_y"`
// 	Scale       `json:"scale"`
// }
