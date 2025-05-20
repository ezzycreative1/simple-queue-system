# 🚀 Task Queue App

A modern task queue system built with **Go** (backend) and **ReactJS** (frontend). It supports creating, managing, and retrying tasks. Features include status filtering, auto-refresh, pagination, and a clean UI with TailwindCSS.

---

## 📌 Features

- ✅ Add new tasks with data input
- 🔁 Auto-refresh task list every 5 seconds
- 🔍 Filter tasks by status: `pending`, `done`, `failed`
- ♻️ Retry failed tasks instantly
- 📃 Pagination support for large datasets
- ✨ Responsive UI with TailwindCSS

---

## 📁 Project Structure

task-queue-app/
├── backend/ # Golang REST API
│ ├── main.go
│ └── Dockerfile
├── frontend/ # ReactJS frontend
│ ├── src/
│ └── Dockerfile
├── docker-compose.yml # Compose for full app
└── README.md

---

## 🐳 Getting Started with Docker

### 📋 Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### ▶️ Run the app

```bash
docker-compose up --build
Frontend runs on: http://localhost:3000

Backend API: http://localhost:8080

💻 Running Locally (Dev Mode)
Backend

cd backend
go run main.go
Frontend

cd frontend
npm install
npm run dev

📡 API Endpoints

Method	Endpoint	    Description
POST	/api/enqueue	Add a new task
GET	    /api/tasks	    Get list of tasks
POST	/api/retry/:id	Retry a failed task by ID

📸 Screenshot



⚙️ Built With
Backend: Go 1.22, native HTTP, JSON, net/http
Frontend: ReactJS, TailwindCSS, React Toastify
Docker: Multi-service development with Docker Compose

🧑‍💻 Developer
Created by Adi Setyadharma

🪪 License
Licensed under the MIT License.