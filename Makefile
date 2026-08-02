.PHONY: backend-test backend-cover frontend-test frontend-cover up down

backend-test:
	cd backend && go test ./... -v

backend-cover:
	cd backend && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

frontend-test:
	cd frontend && npm run test:run -- --reporter verbose

frontend-cover:
	cd frontend && npx c8 --reporter=text --reporter=lcov npx vitest run

up:
	docker compose up -d --build

down:
	docker compose down
