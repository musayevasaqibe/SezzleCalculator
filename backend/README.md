# Sezzle Calculator

A small example calculator app with a React + TypeScript frontend and a Go backend.

Tech stack
- Frontend: React + TypeScript (Vite). Tests: Vitest. Coverage: c8 (LCOV).
- Backend: Go (net/http). Tests: go test. Coverage: go tool cover.

Quick start — local (development)
1. Frontend
cd frontend
npm install
npm run dev           # starts dev server at http://localhost:5173
npm run test          # run unit tests
# generate coverage
npx c8 --reporter=lcov --reporter=text npx vitest run

2. Backend
cd backend
go mod tidy
go run .              # runs server on :8080
go test ./... -v
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage/backend/index.html

API
POST /api/calc
Content-Type: application/json
Body: {"op":"add","a":11,"b":12}
Response: {"result":23,"error":null}

Docker (optional)
docker compose build
docker compose up -d
# view backend logs:
docker compose logs -f backend

CI
GitHub Actions workflow runs frontend and backend tests and collects coverage artifacts.

