build:
	docker-compose build --no-cache
up:
	docker-compose up
up-detached:
	docker-compose up -d
test:
	go test -v ./crawler/ ./rates/ ./api/
destroy:
	docker-compose down