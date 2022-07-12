SHELL := /bin/bash
.PHONY: init

default: init

init:
	docker compose down -v
	docker compose up -d