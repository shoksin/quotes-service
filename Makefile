open-db:
	docker exec -it quotes_postgres psql -U quotes_user -d quotes_db
clean-images:
	docker rmi $(shell docker images -q)
clean-containers:
	docker rm $(shell docker ps -aq)