.PHONY: nuke

nuke: ## Stop docker containers, clean data, workspace and volumes
	docker-compose down
	docker rm -f $(docker ps -a -q)
	docker volume rm $(docker volume ls -q)
