# Sezzle Calculator

A small example calculator application with a React + TypeScript frontend and a Go backend.

Tech stack
- Frontend: React + TypeScript (Vite). Tests: Vitest. Coverage: c8 / lcov.
- Backend: Go (net/http). Tests: `go test`. Coverage: `go tool cover`.

Local development

1) Frontend
cd frontend
npm install
npm run dev           # development server (default: http://localhost:5173)
npm run test          # run unit tests
# generate coverage
npx c8 --reporter=lcov --reporter=text npx vitest run

2) Backend
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

Notes
- The API returns JSON in the form: { "result": <number|null>, "error": <string|null> }.
- For local frontend dev, ensure the frontend is configured to call http://localhost:8080 or update the API URL using environment variables.
