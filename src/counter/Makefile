.PHONY: test
compose = docker compose
lines = 300

ifeq (,$(wildcard .env))
    $(error .env file not found)
endif

include .env
export $(shell sed 's/=.*//' .env)

# DATABASE_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_MASTER_PORT)/$(POSTGRES_DB)?sslmode=disable

init:
	cp env.example .env && sudo docker build -f patroni/Dockerfile.citus -t patroni-citus patroni

build:
	sudo $(compose) up --build -d
#	 sudo $(compose) up --build -d && sudo $(compose) -f patroni/docker-compose-citus.yml up -d

down:
	sudo $(compose) down

stop:
	sudo $(compose) stop

# migrate-up:
# 	sudo $(compose) exec web migrate -database $(DATABASE_URL) -path /code/db/migrations up

# migrate-down:
# 	sudo $(compose) exec web migrate -database $(DATABASE_URL) -path /code/db/migrations down

web-logs:
	sudo $(compose) logs -f web

logs:
	sudo $(compose) logs -f

pg-shell:
	sudo $(compose) exec db bash

redis-cli:
	sudo $(compose) exec redis redis-cli