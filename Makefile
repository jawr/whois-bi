SOURCE := $(shell git rev-parse --show-toplevel)

.PHONY: build start stop clean nuke newtoolbox test

build: ## Build docker image
	docker-compose build

start: ## Start
	docker-compose up

test: ## test
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --remove-orphans
	docker-compose -f docker-compose.test.yml down --volumes

stop: ## Stop docker containers
	docker-compose stop

clean:stop ## Stop docker containers, clean data and workspace
	docker-compose down -v --remove-orphans

nuke:clean ## Stop docker containers, clean data, workspace and volumes
	docker image prune

newtoolbox: ## create a new toolbox
	docker-compose up -d --no-deps --build  toolbox
