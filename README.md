# WebPageAnalyzer

## Main Steps to Build & Deploy

### 1. **Build Locally**
```sh
go build -o webpageanalyzer ./cmd/server
```

### 2. Run Tests
```sh
go test -v ./...
```

### 3. **Run Locally**
```sh
go run ./cmd/server/main.go
```
- The server will start on `localhost:8080`.

### 4. **Run React Frontend (in `fe` folder)**
```sh
cd fe
npm install
npm run dev
```
- The frontend will start on `http://localhost:5173` by default.

### 5. **Run with Docker**
Build the backend Docker image:
```sh
docker build -t webpageanalyzer .
```
Run the backend container:
```sh
docker run -p 8080:8080 webpageanalyzer
```

To build and run the frontend (from `fe`):
```sh
cd fe
docker build -t webpageanalyzer-frontend .
docker run -p 5173:80 webpageanalyzer-frontend
```

### 6. **Run with Docker Compose (Backend + Frontend)**
If you have a `docker-compose.yml`:
```sh
docker-compose up --build
```
- Backend will be available at `http://localhost:8080`
- Frontend will be available at `http://localhost:5173`

---

## Assumptions & Decisions

- The API expects a POST request to `/analyzer` with JSON body:  
  `{ "webpageUrl": "https://example.com" }`
- Only basic HTML analysis is performed (title, headings, links, login form detection, etc.).
- CORS is enabled for `http://localhost:5173` (assumed frontend).
- Only public, accessible URLs are supported.
- Error handling is basic; invalid URLs or unreachable pages return HTTP 400.

---

## Suggestions for Improvement

- Add authentication and rate limiting for production use.
- Improve error messages and validation.
- Add more detailed analysis (e.g., images, scripts, SEO tags).
- Add more unit and integration tests.
- Support for more HTML versions and edge cases.
- Add OpenAPI/Swagger documentation.
- Allow configuration of CORS and server port via environment variables.
- Add Correlation ID for tracing requests.

---