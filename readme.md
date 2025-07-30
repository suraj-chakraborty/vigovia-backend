# 🧳 Itinerary PDF Generator (Go + Gin)

This project generates a downloadable PDF itinerary from a list of travel plan using the Go `gofpdf` library and Gin web framework.

---

## 📦 Prerequisites

Make sure you have the following installed:

- [Go 1.24+](https://golang.org/dl/)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/) *(optional, for containerized run)*

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/suraj-chakraborty/vigovia-backend.git
cd vigovia-backend

### Install Go Dependencies
go mod download

##🧪 Run the Server Locally
go run main.go
```

## Server will start at

# <http://localhost:8080>

# 📤 API Usage

```bash
Endpoint
POST /generate
```

### Content-Type

### application/json

## Sample.json file contain some demo json to try it

# 🐳 Run with Docker

```bash
git clone https://github.com/suraj-chakraborty/vigovia-backend.git

cd vigovia-backend

run docker-compose up --build
 ```
