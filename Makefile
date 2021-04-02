SOURCE := $(shell git rev-parse --show-toplevel)

.PHONY: build status logs start stop clean nuke

build: ## Build docker image
	docker-compose -f build/dev.docker-compose.yml build

start: ## Start
	docker-compose -f build/dev.docker-compose.yml up

status: ## Get status of containers
	docker-compose -f build/dev.docker-compose.yml ps

logs: ## Get logs of containers
	docker-compose -f build/dev.docker-compose.yml logs -f

update: ## Pull latest images
	docker-compose -f build/dev.docker-compose.yml pull

stop: ## Stop docker containers
	docker-compose -f build/dev.docker-compose.yml stop

clean:stop ## Stop docker containers, clean data and workspace
	docker-compose -f build/dev.docker-compose.yml down -v --remove-orphans

prune:stop ## Stop and prune obsolete images
	docker image prune

nuke: ## Stop docker containers, clean data, workspace and volumes
	docker-compose -f build/dev.docker-compose.yml down
	docker rm -f $(docker ps -a -q)
	docker volume rm $(docker volume ls -q)

newtoolbox: ## create a new toolbox
	docker-compose -f build/dev.docker-compose.yml up -d --no-deps --build  toolbox
