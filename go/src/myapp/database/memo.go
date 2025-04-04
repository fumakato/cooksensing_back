package database

// "myapp/model"
/*

databaseフォルダ内のdatabase.go以外について
それぞれのテーブルに対して追加・削除・修正・探索を行うもの
それ以外にも単純なものについては各テーブルごとに機能持たせるかも．
userのexistみたいな
userにはちょっと詳しい操作も入れてる


複雑なものはcontrollerにトランザクションでつくる．
単純なものはcontrollerで呼び出すだけ

databaseの中では
『　var db *gorm.DB　』
が自由に使えるように設定してある．
実際にdbを操作するのはdatabaseフォルダ内のみとなる．
逆に，de操作以外はdatabaseフォルダ内では行ってはいけない

探索についてはidでするのはあまり考えられない．
FirstじゃなくてWhereで，id以外のカラムについても探索の関数を作る必要があり


//トランザクションの例
// トランザクションの開始
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // パニックが発生した場合、ロールバック
		}
	}()

	// ユーザーの検索
	userID := uint(1) // 例としてID 1を使用
	user, err := repository.FindUserByID(tx, userID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error finding user: %v", err)
	}

	// ユーザーの情報を更新
	user.Name = "Updated Name"
	if err := repository.UpdateUser(tx, user); err != nil {
		tx.Rollback()
		log.Fatalf("Error updating user: %v", err)
	}

	// トランザクションのコミット
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	log.Println("Transaction committed successfully")

*/
