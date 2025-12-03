# QuizLit - Smart Quiz Generator

A modern quiz application that allows users to upload PDF files and generate quizzes from the content using AI. Built with Next.js frontend and Go backend, powered by Supabase Auth.

## Tech Stack

| Frontend | Backend |
|----------|---------|
| Next.js 15 | Go 1.19+ |
| TypeScript | PostgreSQL (Supabase) |
| Tailwind CSS | OpenAI / Senopati LLM |
| Supabase Auth | JWT Authentication |

## Quick Start

### Prerequisites

- Node.js v18+
- Go v1.19+
- Supabase account (free tier available)

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/pbkk-quizlit.git
cd pbkk-quizlit
```

### 2. Setup Backend

```bash
cd be
cp .env.example .env
# Edit .env with your credentials
go mod tidy
go run *.go
```

Backend runs on: `http://localhost:8080`

### 3. Setup Frontend

```bash
cd fe
cp .env.local.example .env.local
# Edit .env.local with your credentials
npm install
npm run dev
```

Frontend runs on: `http://localhost:3000`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/quizzes/upload` | Upload PDF and generate quiz |
| POST | `/api/v1/quizzes/generate` | Generate quiz from text |
| GET | `/api/v1/quizzes` | List all quizzes |
| GET | `/api/v1/quizzes/:id` | Get quiz by ID |
| DELETE | `/api/v1/quizzes/:id` | Delete quiz |
| POST | `/api/v1/quizzes/:id/submit` | Submit quiz attempt |

## Project Structure

```
pbkk-quizlit/
├── be/                 # Go Backend
│   ├── internal/       # Internal packages
│   │   ├── api/        # Server setup
│   │   ├── handlers/   # HTTP handlers
│   │   ├── services/   # Business logic (AI, RAG, PDF)
│   │   └── models/     # Data models
│   ├── main.go         # Entry point
│   └── .env            # Backend config
│
├── fe/                 # Next.js Frontend
│   ├── app/            # App Router pages
│   │   ├── dashboard/  # Quiz management
│   │   ├── create/     # Quiz creation
│   │   ├── quiz/       # Quiz taking & results
│   │   ├── history/    # Quiz history
│   │   └── lib/        # API clients & utilities
│   ├── components/     # Reusable components
│   └── .env.local      # Frontend config
│
└── README.md           # This file
```

## Deployment

### Remote VM Deployment

1. SSH into your VM
2. Install Node.js and Go
3. Clone the repository
4. Update environment variables with VM IP addresses
5. Build and run:

```bash
# Backend
cd be
go build -o quizlit-backend
./quizlit-backend

# Frontend
cd fe
npm run build
npm start
```

### Important for Deployment

Update these in your environment files:
- `NEXT_PUBLIC_API_URL` → Your backend URL
- `CORS_ORIGIN` → Your frontend URL