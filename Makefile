
.PHONY: help dev prod staging build clean logs shell test migrate

RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m  
BLUE=\033[0;34m
NC=\033[0m

DEV_COMPOSE=docker-compose.dev.yaml

up:
	@docker-compose -f $(DEV_COMPOSE) --env-file .env.dev up -d --build

start:
	@docker-compose -f $(DEV_COMPOSE) --env-file .env.dev start

stop:
	@docker-compose -f $(DEV_COMPOSE) --env-file .env.dev stop

logs:
	@docker-compose -f $(DEV_COMPOSE) --env-file .env.dev logs

reset-up:
	@echo "$(RED)Hard reset - this will delete all data!!!!!$(NC)"
	@read -p "Continue? (y/N):" confirm && [ "$$confirm" = "y" ]
	@docker-compose -f $(DEV_COMPOSE) --env-file .env.dev down -v
	@docker volume prune -a
	@docker ps -a
	@docker volume ls
	@make up

dev:
	@echo "$(GREEN)Starting development environment... $(NC)"
	@make start
	air

compose-config:
	docker-compose --file $(DEV_COMPOSE) --env-file .env.dev config
