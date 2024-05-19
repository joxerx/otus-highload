.PHONY: test
compose = docker compose
lines = 300

init:
	cp env.example .env

build:
	sudo $(compose) up --build -d

down:
	sudo $(compose) down

stop:
	sudo $(compose) stop

web-build:
	sudo $(compose) up --build -d web

web-logs:
	sudo $(compose) logs -f web

logs:
	sudo $(compose) logs -f

pg-shell:
	sudo $(compose) exec db bash
