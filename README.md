# Go Kafka Microservices API

This project is a backend service built with Go, PostgreSQL, and Kafka, following a microservices-oriented architecture. It provides basic user and wallet transaction features with JWT authentication.

---

## ⚠️ Note from the Author

I would like to apologize as I did not have enough time to fully complete the test with the level of detail I originally intended, nor to properly follow a full Git workflow (such as MR/PR structure, commits, etc.).

However, I was able to implement some interesting features and a solid project structure. I hope you find it valuable.

Also, error handling and proper HTTP status code management were not fully refined due to time constraints.

---

## 🚀 Requirements

- Docker installed
- Docker Compose installed

---

## 🔐 Environment Variables

The file `.env-example` contains all required environment variables to run the project locally.

You should copy it:

```bash
cp .env-example .env
```

Then adjust values if needed.

## 📌 Authentication Rules

Only the POST `/users` route is publicly accessible (no authentication required)

All other routes require a valid JWT token

## 🐳 Running the project

After copy `.env`, you will be able to run `make up`, and use the project with the ports specified `3001` and `3002`.

Note: The first run may take a while to build the project.

### Here is another Make commands:

Start services `make up`

Stop services `make down`

View logs `make logs`

Restart services `make restart`

Build containers `make build`

Check running containers `make ps`

## 📦 Tech Stack

Go (Chi Router)

PostgreSQL

Kafka

Docker / Docker Compose

JWT Authentication

## 💡 Suggestions for Improvement

If continuing this project, the following improvements would be recommended:

Better structured error handling (custom error types + HTTP status mapping)

Idempotency key support for transactions

Pagination for listing endpoints

Unit and integration tests

Swagger UI integration for API documentation

Improved logging strategy (structured logs)

CI/CD pipeline (GitHub Actions)
