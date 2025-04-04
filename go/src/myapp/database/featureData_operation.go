package database

import (
	"myapp/model"
	"time"
)

// データの追加
func AddFeatureData(featureData model.FeatureData) error {
	if err := db.Create(&featureData).Error; err != nil {
		return err
	}
	return nil
}

// データの削除
func DeleteFeatureData(id uint) error {
	if err := db.Delete(&model.FeatureData{}, id).Error; err != nil {
		return err
	}
	return nil
}

// データの修正
func UpdateFeatureData(featureData model.FeatureData) error {
	if err := db.Save(&featureData).Error; err != nil {
		return err
	}
	return nil
}

// データの探索
func FindFeatureDataByID(id uint) (model.FeatureData, error) {
	var featureData model.FeatureData
	if err := db.First(&featureData, id).Error; err != nil {
		return featureData, err
	}
	return featureData, nil
}

// userIDで探索した全部（下のやつに0入れたら全く同じことするからこの関数は実質いらない子になってしまいました）
func GetFeatureDataByUserID(userID uint) ([]model.FeatureData, error) {
	var featureData []model.FeatureData
	if err := db.Where("user_id = ?", userID).Find(&featureData).Error; err != nil {
		return nil, err
	}
	return featureData, nil
}

// userIDで探索してdays以内のデータを探す
func GetFeatureDataByUserIDWithinDays(userID uint, days int) ([]model.FeatureData, error) {

	println("------GetFeatureDataByUserIDWithinDays------")
	var featureData []model.FeatureData

	// daysが0の場合、全てのデータを取得
	if days == 0 {
		println("------days=0------")
		if err := db.Where("user_id = ?", userID).Find(&featureData).Error; err != nil {
			return nil, err
		}
		return featureData, nil
	}

	println("------days=!0------")
	// 現在から指定された日数前の日時を計算
	daysAgo := time.Now().AddDate(0, 0, -days)

	// user_idが一致し、かつcreated_atが指定日数以内のデータを取得
	if err := db.Where("user_id = ? AND date >= ?", userID, daysAgo).Find(&featureData).Error; err != nil {
		return nil, err
	}

	return featureData, nil
}
