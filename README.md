# FileShare 📂

<div align="center">
  <a href="https://dribbble.com/shots/15333756-File-Transfer-Animated-Button">
    <img src="https://cdn.dribbble.com/userupload/21251138/file/original-c99f57c05df70ffc48ab1b50bdc5c59f.gif" width="400" style="border-radius: 8px; object-fit: cover; height: 225px;">
  </a>
</div>
  *Secure, fast, and reliable file sharing made simple.*
  
  [![Docker](https://img.shields.io/badge/Docker-Ready-blue)](https://www.docker.com/)
</div>

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🔐 **Google OAuth2** | Secure authentication with Google accounts |
| 🔑 **JWT Auth** | Token-based authentication for API requests |
| 📤 **Upload/Download** | Seamless file transfer capabilities |
| ⚡ **Redis Caching** | High-performance response times |
| 💾 **PostgreSQL** | Reliable and scalable data storage |
| 🐳 **Docker Support** | Easy deployment with containerization |

## 🛠️ Tech Stack

<div align="center">

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/gin-%23008ECF.svg?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgresql-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

</div>

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/fileshare-pro.git
cd fileshare-pro

# Configure environment variables
cp .env.example .env
# Edit .env with your configuration

# Launch with Docker
docker-compose up --build
```

## 📚 API Endpoints

| Method    | Endpoint                    | Description                          | Authentication    |
|-----------|-----------------------------|--------------------------------------|-------------------|
| POST      | `/upload`                   | Upload files to the server           | JWT Required      |
| GET       | `/files/{id}`               | Download a specific file             | JWT Required      |
| GET       | `/auth/google`              | Initiate Google OAuth login          | None              |
| GET       | `/auth/google/callback`     | OAuth callback handler               | None              |
| GET       | `/api/files`                | List all user files                  | JWT Required      |
| DELETE    | `/api/files/{id}`           | Delete a specific file               | JWT Required      |

## 🔧 Configuration

Configure your .env file

```.env
GOOGLE_CLIENT_ID=Your_Client_ID
GOOGLE_CLIENT_SECRET=Your_Client_Secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
DB_CONNECTION=postgres://<DB_USERNAME>:<DB_PASSWORD>@localhost:5432/<DB_Name>?sslmode=disable
JWT_SECRET=your_jwt_secret
REDIS_ADDR=localhost:6379
FILE_EXPIRY_DAYS=set-file-expiry
RATE_LIMIT=set-rate-limit
RATE_LIMIT_WINDOW=set-rate-limit-window
```
