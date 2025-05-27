open-db:
	docker exec -it quotes_postgres psql -U quotes_user -d quotes_db
clear-images:
	docker rmi $(shell docker images -q)
clear-containers:
	docker rm $(shell docker ps -aq)
start:
	docker-compose up -d
start-rebuild:
	docker-compose up --build
stop:
	docker-compose down