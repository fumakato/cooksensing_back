-include .env

up:
	docker compose up -d

db:
	docker exec -it ${MYSQL_CONTAINER_HOST} mysql -u ${MYSQL_ROOT_USER} -p${MYSQL_ROOT_PASSWORD} -D ${MYSQL_DATABASE}