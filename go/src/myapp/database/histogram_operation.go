package database

import (
	"fmt"
	"myapp/model"

	"time"

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

// ヒストグラムを作る
// GenerateAndStoreHistogramData calculates and stores histogram data for AveragePace and AccelerationStandardDeviation
func GenerateAndStoreHistogramData() error {
	// db := database.GetDB()

	// Generate and store histogram data for AveragePace
	err := generateAndStoreHistogram("average_pace", 1)
	if err != nil {
		return err
	}

	// Generate and store histogram data for AccelerationStandardDeviation
	err = generateAndStoreHistogram("acceleration_standard_deviation", 2)
	if err != nil {
		return err
	}

	return nil
}

// 実際にヒストグラムを作成する部分
// generateAndStoreHistogram calculates and stores histogram data for a given field
func generateAndStoreHistogram(field string, displayItemID uint) error {

	var max, min, dataRange float32

	// Retrieve the maximum and minimum values
	err := db.Table("best_data").Select(fmt.Sprintf("MAX(%s)", field)).Scan(&max).Error
	if err != nil {
		return err
	}
	err = db.Table("best_data").Select(fmt.Sprintf("MIN(%s)", field)).Scan(&min).Error
	if err != nil {
		return err
	}
	dataRange = max - min

	// Output the results
	fmt.Printf("%s - Max: %f, Min: %f, Range: %f\n", field, max, min, dataRange)

	// Calculate histogram data
	histogram := calculateHistogram(field, min, max, dataRange)
	histogram.DisplayItemID = displayItemID
	histogram.Max = max
	histogram.Min = min
	histogram.Range = dataRange
	histogram.ActionID = 1 // Set ActionID to 1
	// ActionIDは1が輪切りになってる．
	histogram.CreatedAt = time.Now() // 手動で現在時刻を設定
	histogram.UpdatedAt = time.Now() // 手動で現在時刻を設定

	// Check if a histogram already exists for this DisplayItemID and ActionID
	var existingHistogram model.Histogram
	result := db.Where("display_item_id = ? AND action_id = ?", displayItemID, 1).First(&existingHistogram)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	// fmt.Printf("Before Save - CreatedAt: %v, UpdatedAt: %v", histogram.CreatedAt, histogram.UpdatedAt)

	if result.RowsAffected > 0 {
		// Update the existing histogram record
		histogram.ID = existingHistogram.ID
		if err := db.Save(&histogram).Error; err != nil {
			return err
		}
	} else {
		// Create a new histogram record
		if err := db.Create(&histogram).Error; err != nil {
			return err
		}
	}

	return nil
}

// 計算部分
// calculateHistogram calculates the histogram data for a given field in the BestData table
func calculateHistogram(field string, minValue, max, dataRange float32) model.Histogram {
	var histogram model.Histogram

	// Calculate the bin size
	binSize := dataRange / 10.0

	// Calculate the count for each bin
	for i := 0; i < 10; i++ {
		minBin := minValue + binSize*float32(i)
		maxBin := minValue + binSize*float32(i+1)

		var count int64
		// db.Table("best_data").Where(fmt.Sprintf("%s >= ? AND %s < ?", field, field), minBin, maxBin).Count(&count)

		if i == 9 {
			// Ensure the maximum value is included in the last bin
			db.Table("best_data").Where(fmt.Sprintf("%s >= ? AND %s <= ?", field, field), minBin, max).Count(&count)
		} else {
			db.Table("best_data").Where(fmt.Sprintf("%s >= ? AND %s < ?", field, field), minBin, maxBin).Count(&count)
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
