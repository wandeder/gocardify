up:
	docker compose -f docker-compose.yml  up -d --build
stop:
	docker compose stop
start:
	docker compose start
force:
	docker compose -f docker-compose.yml  up -d --build --force-recreate

