# 🧩 Simple Queue System (Go + React)

A full-stack task queue management app built with **Go** for the backend and **React** for the frontend. It includes features like task retry, filtering, pagination, and auto-refresh.

## 📦 Features

- Queue task data via API
- View and filter task status (pending, failed, done)
- Retry failed tasks
- Pagination and auto-refresh
- Built with Docker and Docker Compose for easy deployment

---

## 📁 Project Structure

```
project-root/
├── backend/        # Go backend
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── frontend/       # React frontend
│   ├── src/
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🚀 Quick Start (Docker Compose)

Make sure you have **Docker** and **Docker Compose** installed.

### 1. Clone the repository
```bash
git clone https://github.com/ezzycreative1/simple-queue-system.git
cd simple-queue-system
```

### 2. Build and run containers
```bash
docker-compose up --build
```

- React app will run at: `http://localhost:3000`
- Go backend API will run at: `http://localhost:8080`

### 3. Stop containers
```bash
docker-compose down
```

---

## 🛠 Backend (Go)

### Build and Run Manually (Optional)
```bash
cd backend
go mod tidy
go build -o queue-backend
./queue-backend
```

### API Endpoints
- `GET /api/tasks` — fetch task list
- `POST /api/enqueue` — enqueue a new task
- `POST /api/retry/{id}` — retry a failed task

---

## 💻 Frontend (React)

### Run in Development (Optional)
```bash
cd frontend
npm install
npm start
```

### Build for Production
```bash
npm run build
```

---

## 🐳 Docker Compose

### `docker-compose.yml`
```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: queue-backend
    ports:
      - "8080:8080"
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: queue-frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped
```

---

## 📝 License

MIT License © 2025

---

For improvements or issues, feel free to open a PR or issue!
