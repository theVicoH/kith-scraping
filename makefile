ENV_FILE=.env.dev

DOCKER_ENVIRONMENT=docker-compose.dev.yml

DOCKER_COMPOSE=docker compose -f $(DOCKER_ENVIRONMENT) --env-file $(ENV_FILE)

.PHONY: up
up:
	@echo "Starting up the application..."
	$(DOCKER_COMPOSE) up
	@echo "Application is running."

.PHONY: build
build:
	@echo "Building the application..."
	$(DOCKER_COMPOSE) up --build
	@echo "Application has been built."

.PHONY: down
down:
	@echo "Stopping the application..."
	$(DOCKER_COMPOSE) down --remove-orphans
	@echo "Application has been stopped."

.PHONY: logs
logs:
	@echo "Fetching logs..."
	$(DOCKER_COMPOSE) logs -f
	@echo "Logs fetched."

.PHONY: shell
shell:
	@echo "Opening a shell in the application container..."
	$(DOCKER_COMPOSE) exec -ti backend sh
	@echo "Shell closed."

.PHONY: restart
restart:
	$(DOCKER_COMPOSE) restart



default: up