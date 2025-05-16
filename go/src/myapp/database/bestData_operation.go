package database

import (
	"fmt"
	"myapp/model"
	"time"
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
			AveragePaceClass:              100,
			AccelerationStdDevClass:       100,
		}

		// Use GORM's Upsert to insert the record if it doesn't exist, or update it if it does
		err = db.Where("user_id = ?", result.UserID).Assign(bestData).FirstOrCreate(&bestData).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateBestDataFromFeatureDataOneMonth() error {
	// 🔸 1ヶ月前の日付を計算（基準日）
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	// -------------------------------
	// Step 1: 直近1ヶ月のデータがあるユーザーに対して、
	//        ・AveragePaceの最大
	//        ・AccelerationStandardDeviationの最小
	//        を集計する
	// -------------------------------
	var recentResults []struct {
		UserID                           uint
		MaxAveragePace                   float32
		MinAccelerationStandardDeviation float32
	}
	err := db.Table("feature_data").
		Where("date >= ?", oneMonthAgo). // ← 1ヶ月以内のデータに限定
		Select("user_id, MAX(average_pace) as max_average_pace, MIN(acceleration_standard_deviation) as min_acceleration_standard_deviation").
		Group("user_id").
		Scan(&recentResults).Error
	if err != nil {
		return err
	}

	// 🔸 直近1ヶ月のデータがあったユーザーを記録しておく
	userIDsWithRecentData := make(map[uint]bool)
	for _, r := range recentResults {
		userIDsWithRecentData[r.UserID] = true

		// ベストデータを構成して保存
		best := model.BestData{
			UserID:                        r.UserID,
			AveragePace:                   r.MaxAveragePace,
			AccelerationStandardDeviation: r.MinAccelerationStandardDeviation,
			AveragePaceClass:              100, // 仮クラス
			AccelerationStdDevClass:       100,
		}
		if err := db.Where("user_id = ?", r.UserID).Assign(best).FirstOrCreate(&best).Error; err != nil {
			return err
		}
	}

	// -------------------------------
	// Step 2: 直近1ヶ月以内にデータがないユーザーに対して、
	//        「最新の1件（全期間中の直近）だけ」を使ってBestDataを作成
	// -------------------------------

	// まず全ユーザーの user_id 一覧を取得
	var allResults []struct {
		UserID                           uint
		MaxAveragePace                   float32
		MinAccelerationStandardDeviation float32
	}
	err = db.Table("feature_data").
		Select("user_id, MAX(average_pace) as max_average_pace, MIN(acceleration_standard_deviation) as min_acceleration_standard_deviation").
		Group("user_id").
		Scan(&allResults).Error
	if err != nil {
		return err
	}

	// 🔸 1ヶ月以内にデータがなかったユーザーに対してのみ処理
	for _, r := range allResults {
		if userIDsWithRecentData[r.UserID] {
			continue // すでに処理済みのユーザーはスキップ
		}

		// 全期間の中から「最新の1件」を取得
		var latest model.FeatureData
		err := db.Where("user_id = ?", r.UserID).
			Order("date DESC").
			Limit(1).
			First(&latest).Error

		if err != nil {
			return fmt.Errorf("failed to get latest feature_data for user_id=%d: %w", r.UserID, err)
		}

		// その1件からBestDataを構成して保存
		best := model.BestData{
			UserID:                        latest.UserID,
			AveragePace:                   latest.AveragePace,
			AccelerationStandardDeviation: latest.AccelerationStandardDeviation,
			AveragePaceClass:              100,
			AccelerationStdDevClass:       100,
		}

		if err := db.Where("user_id = ?", latest.UserID).Assign(best).FirstOrCreate(&best).Error; err != nil {
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

// AssignBestClassByUserID assigns histogram class info to a user's BestData based on current histogram bins.
func AssignBestClassByUserID(userID uint) error {
	// ① ユーザーのBestDataを取得
	var bestData model.BestData
	if err := db.Where("user_id = ?", userID).First(&bestData).Error; err != nil {
		return fmt.Errorf("failed to find BestData: %w", err)
	}

	// ② ヒストグラム情報を取得（DisplayItemID 1: pace, 2: acc）
	var paceHist, accHist model.Histogram

	if err := db.Where("display_item_id = ? AND action_id = ?", 1, 1).First(&paceHist).Error; err != nil {
		return fmt.Errorf("failed to get pace histogram: %w", err)
	}
	if err := db.Where("display_item_id = ? AND action_id = ?", 2, 1).First(&accHist).Error; err != nil {
		return fmt.Errorf("failed to get acceleration histogram: %w", err)
	}

	// ③ クラス計算
	bestData.AveragePaceClass = calculateClass(bestData.AveragePace, paceHist.Min, paceHist.Range)
	bestData.AccelerationStdDevClass = calculateClass(bestData.AccelerationStandardDeviation, accHist.Min, accHist.Range)

	// ④ 更新保存
	if err := db.Save(&bestData).Error; err != nil {
		return fmt.Errorf("failed to save BestData class info: %w", err)
	}

	fmt.Printf("✅ UserID %d のクラス情報を更新: pace=%d, acc=%d\n",
		bestData.UserID, bestData.AveragePaceClass, bestData.AccelerationStdDevClass)

	return nil
}

// AssignBestClassToAll assigns class info to all users' BestData based on current histogram.
func AssignBestClassToAll() error {
	// ① 全 BestData を取得
	var bestDataList []model.BestData
	if err := db.Find(&bestDataList).Error; err != nil {
		return fmt.Errorf("failed to fetch BestData list: %w", err)
	}

	// ② 必要なヒストグラム（DisplayItemID 1と2）を取得
	var paceHist, accHist model.Histogram
	if err := db.Where("display_item_id = ? AND action_id = ?", 1, 1).First(&paceHist).Error; err != nil {
		return fmt.Errorf("failed to get pace histogram: %w", err)
	}
	if err := db.Where("display_item_id = ? AND action_id = ?", 2, 1).First(&accHist).Error; err != nil {
		return fmt.Errorf("failed to get acceleration histogram: %w", err)
	}

	// ③ 各ユーザーの BestData に対してクラス計算＆保存
	for _, bd := range bestDataList {
		bd.AveragePaceClass = calculateClass(bd.AveragePace, paceHist.Min, paceHist.Range)
		bd.AccelerationStdDevClass = calculateClass(bd.AccelerationStandardDeviation, accHist.Min, accHist.Range)

		if err := db.Save(&bd).Error; err != nil {
			fmt.Printf("❌ Failed to update user %d: %v\n", bd.UserID, err)
		} else {
			fmt.Printf("✅ Updated class for user %d: pace=%d, acc=%d\n",
				bd.UserID, bd.AveragePaceClass, bd.AccelerationStdDevClass)
		}
	}

	return nil
}

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

// func GetBestDataByUserID(userID uint) ([]model.BestData, error) {
// 	var besteData []model.BestData
// 	if err := db.Where("user_id = ?", userID).Find(&bestData).Error; err != nil {
// 		return nil, err
// 	}
// 	return featureData, nil
// }
