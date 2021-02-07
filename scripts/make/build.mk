.PHONY: build

build: ## Build docker image
	docker-compose build

push: ## Build and push to docker hub
	docker build --build-arg service=api -t whoisbi/api -f services/Dockerfile .
	docker build --build-arg service=manager -t whoisbi/manager -f services/Dockerfile .
	docker build --build-arg service=worker -t whoisbi/worker -f services/Dockerfile .
	docker push whoisbi/api
	docker push whoisbi/manager
	docker push whoisbi/worker

