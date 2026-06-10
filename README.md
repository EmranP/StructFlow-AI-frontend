# StructFlow — AI Project Structure Generator

<div align="center">

![StructFlow Banner](https://img.shields.io/badge/StructFlow-AI%20Powered-6366f1?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCI+PHRleHQgeT0iMjAiIGZvbnQtc2l6ZT0iMjAiPuKmoTwvdGV4dD48L3N2Zz4=)

![Angular](https://img.shields.io/badge/Angular-17-DD0031?style=flat-square&logo=angular)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat-square&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat-square&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

**Describe your project → AI designs the folder structure → Download and start coding**

[Demo](https://struct-flow-ai.vercel.app/)

</div>

---

## 📋 Table of Contents

- [About](#-about)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
  - [Environment Variables](#environment-variables)
- [Database Schema](#-database-schema)

---

## 🧠 About

**StructFlow** is a full-stack web application that uses AI to generate project directory structures based on a developer's description. Instead of manually creating dozens of folders and files when starting a new project, you describe what you're building and the AI instantly produces multiple architecture variants — from a simple starter layout to a full enterprise-grade modular structure.

### The Problem

Every developer knows the pain of starting a new project:

- Googling "go project structure clean architecture" for the 10th time
- Copy-pasting folder layouts from old projects
- Spending 30 minutes creating files before writing a single line of code

### The Solution

Describe your project in plain text → StructFlow generates 3 ready-to-use directory structures → Download as a ZIP and start coding immediately.

---

## ✨ Features

| Feature                      | Description                                                         |
| ---------------------------- | ------------------------------------------------------------------- |
| 🤖 **AI Generation**         | Describe your project in natural language, AI designs the structure |
| 📐 **3 Complexity Levels**   | Simple, Medium, Enterprise variants for every project size          |
| 📁 **Interactive File Tree** | Visual collapsible tree viewer with file type icons                 |
| 📦 **ZIP Download**          | Download the full directory structure instantly                     |
| 🔄 **Live Status Polling**   | Real-time generation progress with animated UI                      |
| 🗂️ **Project Management**    | Create, edit, delete projects with full history                     |
| 📜 **Generation History**    | Track all past generations per project with pagination              |
| 🔐 **JWT Auth**              | Secure authentication with protected routes                         |

---

## 🛠 Tech Stack

### Backend

| Technology     | Version | Purpose          |
| -------------- | ------- | ---------------- |
| **Go**         | 1.25    | Core language    |
| **Fiber**      | v2      | HTTP framework   |
| **PostgreSQL** | —       | Primary database |
| **JWT**        | —       | Authentication   |
| **Docker**     | —       | Containerization |

### Frontend

| Technology         | Version | Purpose                    |
| ------------------ | ------- | -------------------------- |
| **Angular**        | 17      | Core framework             |
| **TypeScript**     | 5.x     | Language                   |
| **SCSS**           | —       | Styling with CSS variables |
| **Angular Router** | —       | Client-side routing        |
| **HttpClient**     | —       | HTTP with JWT interceptor  |
| **Signals**        | —       | Reactive state             |

---

## 🚀 Getting Started

### Prerequisites

- **Go** `1.21+`
- **Node.js** `18+`
- **PostgreSQL** `14+`
- **Docker** (optional, recommended)

---

### Backend Setup

```bash
# Clone the repository
git clone https://github.com/your-username/structflow.git
cd structflow/backend

# Install dependencies
go mod download

# Copy environment config
cp .env.example .env
# Edit .env with your values

# Run database migrations
make migrate-up

# Start the server
make run
# or
go run cmd/main.go
```

**With Docker:**

```bash
cd backend
docker compose up -d
```

---

### Frontend Setup

```bash
cd structflow/frontend

# Install dependencies
npm install

# Start dev server
npm start
# Opens at http://localhost:4200

# Build for production
npm run build
```

---

### Environment Variables

#### Backend `.env`

```env
# Server
APP_PORT=3000
CLIENT_URL=http://localhost:4200

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=structflow
DB_USER=postgres
DB_PASSWORD=your_password

# JWT
JWT_SECRET=your_super_secret_key
JWT_EXPIRES_IN=24h


# AI Provider
AI_API_KEY=your_ai_api_key
```

---

## 🗄 Database Schema

```sql
-- Users
users (id UUID PK, email VARCHAR UNIQUE, password_hash TEXT, created_at, updated_at)

-- Projects
projects (
  id UUID PK,
  user_id UUID FK → users,
  title VARCHAR,
  project_type VARCHAR,
  stack TEXT,
  architecture TEXT,
  features TEXT,
  additional_info TEXT,
  created_at, updated_at
)

-- Generations
generations (
  id UUID PK,
  project_id UUID FK → projects,
  status VARCHAR,       -- pending | process | completed | failed
  error_message TEXT,
  created_at, updated_at
)

-- Templates
templates (
  id UUID PK,
  generation_id UUID FK → generations,
  name VARCHAR,
  description TEXT,
  structure_json JSONB,
  created_at
)

-- Generated Templates (final output)
generated_templates (
  id UUID PK,
  generation_id UUID FK → generations,
  type VARCHAR,         -- simple | medium | enterprise
  content JSONB,        -- { files: [...], directories: [...] }
  created_at
)
```

**Relationships:**

---

<div align="center">
Built with ❤️ by developers, for developers
</div>
