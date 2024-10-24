.PHONY: run

run:
	cp .env.example .env
	docker-compose up -d --remove-orphans --build;