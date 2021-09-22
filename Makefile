start-dev: ## Starting environments for development in docker compose
	@echo " > Start Development ENV..."
	docker-compose up -d
	@echo " > Done Start"

run: ## run app
	@echo " > starting the app .............."
	go run main.go http