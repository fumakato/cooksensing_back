-include .env

up:
	docker compose up -d

up-build:
	docker compose up --build -d

db:
	docker exec -it ${MYSQL_CONTAINER_HOST} mysql -u ${MYSQL_ROOT_USER} -p${MYSQL_ROOT_PASSWORD} -D ${MYSQL_DATABASE}

# 初期値追加
init-db:
	docker exec -it go_container go run scripts/manage_db.go -action init

# テーブル全削除とマイグレーション
dropmigrate-db:
	docker exec -it go_container go run scripts/manage_db.go -action dropmigrate