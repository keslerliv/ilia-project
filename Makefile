APP_NAME=go-kafka-microservices

up:
	docker-compose up --build

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose down && docker-compose up --build

build:
	docker-compose build

ps:
	docker-compose ps

create-message:
	curl -X POST http://localhost:8080/place-order \
	-H "Content-Type: application/json" \
	-d '{ "customer_name": "Hello Kafka last", "coffee_type": "coffe Super" }'
