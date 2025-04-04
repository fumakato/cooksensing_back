package database

import (
	"myapp/model"
)

// データの追加
func AddBestData(bestData model.BestData) error {
	if err := db.Create(&bestData).Error; err != nil {
		return err
	}
	return nil
}

// データの削除
func DeleteBestData(id uint) error {
	if err := db.Delete(&model.BestData{}, id).Error; err != nil {
		return err
	}
	return nil
}

// データの修正
func UpdateBestData(bestData model.BestData) error {
	if err := db.Save(&bestData).Error; err != nil {
		return err
	}
	return nil
}

// データの探索
func FindBestDataByID(id uint) (model.BestData, error) {
	var bestData model.BestData
	if err := db.First(&bestData, id).Error; err != nil {
		return bestData, err
	}
	return bestData, nil
}

// 全件取得 FindAllBestData retrieves all records from the BestData table
func FindAllBestData() ([]model.BestData, error) {
	var bestDataList []model.BestData
	if err := db.Find(&bestDataList).Error; err != nil {
		return nil, err
	}
	return bestDataList, nil
}

// bestdataの作成
func UpdateBestDataFromFeatureData() error {
	// db := database.GetDB()

	// Step 1: Retrieve the best AveragePace and AccelerationStandardDeviation for each UserID
	var results []struct {
		UserID                           uint
		MaxAveragePace                   float32
		MinAccelerationStandardDeviation float32
	}

	err := db.Table("feature_data").
		Select("user_id, MAX(average_pace) as max_average_pace, MIN(acceleration_standard_deviation) as min_acceleration_standard_deviation").
		Group("user_id").
		Scan(&results).Error

	if err != nil {
		return err
	}

	// Step 2: Update or insert the best data for each user in the BestData table
	for _, result := range results {
		bestData := model.BestData{
			UserID:                        result.UserID,
			AveragePace:                   result.MaxAveragePace,
			AccelerationStandardDeviation: result.MinAccelerationStandardDeviation,
		}

		// Use GORM's Upsert to insert the record if it doesn't exist, or update it if it does
		err = db.Where("user_id = ?", result.UserID).Assign(bestData).FirstOrCreate(&bestData).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// 平均の算出
func AveragePaceAndAccelerationStdDev() (float32, float32, error) {
	var averagePace float32
	var accelerationStdDev float32

	// Calculate the average of AveragePace
	err := db.Table("best_data").Select("AVG(average_pace)").Scan(&averagePace).Error
	if err != nil {
		return 0, 0, err
	}

	// Calculate the average of AccelerationStandardDeviation
	err = db.Table("best_data").Select("AVG(acceleration_standard_deviation)").Scan(&accelerationStdDev).Error
	if err != nil {
		return 0, 0, err
	}

	return averagePace, accelerationStdDev, nil
}

// func GetBestDataByUserID(userID uint) ([]model.BestData, error) {
// 	var besteData []model.BestData
// 	if err := db.Where("user_id = ?", userID).Find(&bestData).Error; err != nil {
// 		return nil, err
// 	}
// 	return featureData, nil
// }
