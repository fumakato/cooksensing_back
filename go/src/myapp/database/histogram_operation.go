package database

import (
	"fmt"
	"myapp/model"

	"gorm.io/gorm"
)

// 全件取得 FindAllBestData retrieves all records from the BestData table
func FindAllHistogram() ([]model.Histogram, error) {
	var histogramDataList []model.Histogram
	if err := db.Find(&histogramDataList).Error; err != nil {
		return nil, err
	}
	return histogramDataList, nil
}

// データの追加
func AddHistogram(histogram model.Histogram) error {
	if err := db.Create(&histogram).Error; err != nil {
		return err
	}
	return nil
}

// データの削除
func DeleteHistogram(id uint) error {
	if err := db.Delete(&model.Histogram{}, id).Error; err != nil {
		return err
	}
	return nil
}

// データの修正
func UpdateHistogram(histogram model.Histogram) error {
	if err := db.Save(&histogram).Error; err != nil {
		return err
	}
	return nil
}

// データの探索
func FindHistogramByID(id uint) (model.Histogram, error) {
	var histogram model.Histogram
	if err := db.First(&histogram, id).Error; err != nil {
		return histogram, err
	}
	return histogram, nil
}

// import (
// 	"fmt"
// 	"gorm.io/gorm"
// 	"your_project/database"
// 	"your_project/models"
// )

// GenerateAndStoreHistogramData generates histograms for average pace and acceleration stddev,
// stores them in the histogram table, and updates BestData with class information.
func GenerateAndStoreHistogramData() error {
	// Generate histogram for Average Pace (DisplayItemID: 1)
	histPace, err := generateAndStoreHistogram("average_pace", 1)
	if err != nil {
		return err
	}

	// Generate histogram for Acceleration StdDev (DisplayItemID: 2)
	histAcc, err := generateAndStoreHistogram("acceleration_standard_deviation", 2)
	if err != nil {
		return err
	}

	// Update BestData with class information
	var bestDataList []model.BestData
	if err := db.Find(&bestDataList).Error; err != nil {
		return err
	}

	for _, bd := range bestDataList {
		bd.AveragePaceClass = calculateClass(bd.AveragePace, histPace.Min, histPace.Range)
		bd.AccelerationStdDevClass = calculateClass(bd.AccelerationStandardDeviation, histAcc.Min, histAcc.Range)
		if err := db.Save(&bd).Error; err != nil {
			return err
		}
	}

	return nil
}

// generateAndStoreHistogram calculates and saves histogram data, and returns the histogram.
func generateAndStoreHistogram(field string, displayItemID uint) (model.Histogram, error) {
	var max, min, dataRange float32

	if err := db.Table("best_data").Select(fmt.Sprintf("MAX(%s)", field)).Scan(&max).Error; err != nil {
		return model.Histogram{}, err
	}
	if err := db.Table("best_data").Select(fmt.Sprintf("MIN(%s)", field)).Scan(&min).Error; err != nil {
		return model.Histogram{}, err
	}
	dataRange = max - min

	histogram := calculateHistogram(field, min, max, dataRange)
	histogram.DisplayItemID = displayItemID
	histogram.Max = max
	histogram.Min = min
	histogram.Range = dataRange
	histogram.ActionID = 1

	var existing model.Histogram
	result := db.Where("display_item_id = ? AND action_id = ?", displayItemID, 1).First(&existing)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return model.Histogram{}, result.Error
	}

	if result.RowsAffected > 0 {
		histogram.ID = existing.ID
		if err := db.Save(&histogram).Error; err != nil {
			return model.Histogram{}, err
		}
	} else {
		if err := db.Create(&histogram).Error; err != nil {
			return model.Histogram{}, err
		}
	}

	return histogram, nil
}

// calculateClass maps a value to its corresponding bin (0–9) based on min and range.
func calculateClass(value, min, dataRange float32) uint {
	if dataRange == 0 {
		return 0
	}
	binSize := dataRange / 10.0
	class := int((value - min) / binSize)
	if class > 9 {
		class = 9
	}
	return uint(class)
}

// // ヒストグラムを作る
// // GenerateAndStoreHistogramData calculates and stores histogram data for AveragePace and AccelerationStandardDeviation
// func GenerateAndStoreHistogramData() error {
// 	// db := database.GetDB()

// 	// Generate and store histogram data for AveragePace
// 	err := generateAndStoreHistogram("average_pace", 1)
// 	if err != nil {
// 		return err
// 	}

// 	// Generate and store histogram data for AccelerationStandardDeviation
// 	err = generateAndStoreHistogram("acceleration_standard_deviation", 2)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // 実際にヒストグラムを作成する部分
// // generateAndStoreHistogram calculates and stores histogram data for a given field
// func generateAndStoreHistogram(field string, displayItemID uint) error {

// 	var max, min, dataRange float32

// 	// Retrieve the maximum and minimum values
// 	err := db.Table("best_data").Select(fmt.Sprintf("MAX(%s)", field)).Scan(&max).Error
// 	if err != nil {
// 		return err
// 	}
// 	err = db.Table("best_data").Select(fmt.Sprintf("MIN(%s)", field)).Scan(&min).Error
// 	if err != nil {
// 		return err
// 	}
// 	dataRange = max - min

// 	// Output the results
// 	fmt.Printf("%s - Max: %f, Min: %f, Range: %f\n", field, max, min, dataRange)

// 	// Calculate histogram data
// 	histogram := calculateHistogram(field, min, max, dataRange)
// 	histogram.DisplayItemID = displayItemID
// 	histogram.Max = max
// 	histogram.Min = min
// 	histogram.Range = dataRange
// 	histogram.ActionID = 1 // Set ActionID to 1
// 	// ActionIDは1が輪切りになってる．
// 	histogram.CreatedAt = time.Now() // 手動で現在時刻を設定
// 	histogram.UpdatedAt = time.Now() // 手動で現在時刻を設定

// 	// Check if a histogram already exists for this DisplayItemID and ActionID
// 	var existingHistogram model.Histogram
// 	result := db.Where("display_item_id = ? AND action_id = ?", displayItemID, 1).First(&existingHistogram)

// 	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
// 		return result.Error
// 	}

// 	// fmt.Printf("Before Save - CreatedAt: %v, UpdatedAt: %v", histogram.CreatedAt, histogram.UpdatedAt)

// 	if result.RowsAffected > 0 {
// 		// Update the existing histogram record
// 		histogram.ID = existingHistogram.ID
// 		if err := db.Save(&histogram).Error; err != nil {
// 			return err
// 		}
// 	} else {
// 		// Create a new histogram record
// 		if err := db.Create(&histogram).Error; err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// 計算部分
// 0,1件しかない場合の対策済み
func calculateHistogram(field string, minValue, max, dataRange float32) model.Histogram {
	var histogram model.Histogram

	// Special case: all values are the same
	if dataRange == 0 {
		var count int64
		db.Table("best_data").
			Where(fmt.Sprintf("%s = ?", field), minValue).
			Count(&count)

		histogram.Time1 = uint(count)
		// 他のビンは0（デフォルト値）
		return histogram
	}

	// 通常ケース
	binSize := dataRange / 10.0

	for i := 0; i < 10; i++ {
		minBin := minValue + binSize*float32(i)
		maxBin := minValue + binSize*float32(i+1)

		var count int64
		if i == 9 {
			db.Table("best_data").
				Where(fmt.Sprintf("%s >= ? AND %s <= ?", field, field), minBin, max).
				Count(&count)
		} else {
			db.Table("best_data").
				Where(fmt.Sprintf("%s >= ? AND %s < ?", field, field), minBin, maxBin).
				Count(&count)
		}

		switch i {
		case 0:
			histogram.Time1 = uint(count)
		case 1:
			histogram.Time2 = uint(count)
		case 2:
			histogram.Time3 = uint(count)
		case 3:
			histogram.Time4 = uint(count)
		case 4:
			histogram.Time5 = uint(count)
		case 5:
			histogram.Time6 = uint(count)
		case 6:
			histogram.Time7 = uint(count)
		case 7:
			histogram.Time8 = uint(count)
		case 8:
			histogram.Time9 = uint(count)
		case 9:
			histogram.Time10 = uint(count)
		}
	}

	return histogram
}
