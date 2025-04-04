package database

import (
	"myapp/model"
)

// データの追加
func AddAction(action model.Action) error {
	if err := db.Create(&action).Error; err != nil {
		return err
	}
	return nil
}

// データの削除
func DeleteAction(id uint) error {
	if err := db.Delete(&model.Action{}, id).Error; err != nil {
		return err
	}
	return nil
}

// データの修正
func UpdateAction(action model.Action) error {
	if err := db.Save(&action).Error; err != nil {
		return err
	}
	return nil
}

// データの探索
func FindActionByID(id uint) (model.Action, error) {
	var action model.Action
	if err := db.First(&action, id).Error; err != nil {
		return action, err
	}
	return action, nil
}
