package database

import (
	"myapp/model"
)

// データの追加
func AddDisplayItem(displayItem model.DisplayItem) error {
	if err := db.Create(&displayItem).Error; err != nil {
		return err
	}
	return nil
}

// データの削除
func DeleteDisplayItem(id uint) error {
	if err := db.Delete(&model.DisplayItem{}, id).Error; err != nil {
		return err
	}
	return nil
}

// データの修正
func UpdateDisplayItem(displayItem model.DisplayItem) error {
	if err := db.Save(&displayItem).Error; err != nil {
		return err
	}
	return nil
}

// データの探索
func FindDisplayItemByID(id uint) (model.DisplayItem, error) {
	var displayItem model.DisplayItem
	if err := db.First(&displayItem, id).Error; err != nil {
		return displayItem, err
	}
	return displayItem, nil
}
