package database

import (
	"myapp/model"
)

// データの追加
func AddUser(user *model.User) error {
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// データの削除
func DeleteUser(id uint) error {
	if err := db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// データの修正
func UpdateUser(user model.User) error {
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func FindAllUser() ([]model.User, error) {
	var userDataList []model.User
	if err := db.Find(&userDataList).Error; err != nil {
		return nil, err
	}
	return userDataList, nil
}

// データの探索(First)
func FindUserByID(id uint) (model.User, error) {
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

// データの探索(Where)
func FindUsersByName(name string) ([]model.User, error) {
	var users []model.User
	if err := db.Where("name = ?", name).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindUserByEmail checks if a user with the specified email exists and returns the user data if found
func FindUserByFirebaseAuthUid(firebase_auth_uid string) (model.User, error) {
	var user model.User
	if err := db.Where("firebase_auth_uid = ?", firebase_auth_uid).First(&user).Error; err != nil {
		return user, err // returns the error if user is not found
	}
	return user, nil
}

// FindUserByEmail checks if a user with the specified email exists and returns the user data if found
func FindUserByEmail(email string) (model.User, error) {
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err // returns the error if user is not found
	}
	return user, nil
}

// レコードが存在するかどうか
func IsUserExists(id uint) (bool, error) {
	var count int64
	if err := db.Model(&model.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// 名前とfirebase_auth_uidで完全一致するか判定
func FindUserByNameAndFirebaseAuthUid(name string, firebaseAuthUid string) (model.User, error) {
	var user model.User
	err := db.Where("name = ? AND firebase_auth_uid = ?", name, firebaseAuthUid).First(&user).Error
	return user, err
}
