# Frontend Technical Documentation

## Overview

Vue 3 Single Page Application (SPA) built with TypeScript, consuming the API Gateway via HTTP/REST.

## Tech Stack

- **Vue 3** - Progressive JavaScript framework with Composition API
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool and dev server
- **Vue Router** - Client-side routing
- **Pinia** - State management
- **Axios** - HTTP client with interceptors

## Architecture

```
Component (Vue)
    ↓
Store (Pinia)
    ↓
API Service (Axios)
    ↓
API Gateway
```

## Directory Structure

```
frontend/
├── src/
│   ├── assets/          # Static assets
│   ├── components/      # Reusable components
│   ├── router/          # Route definitions
│   ├── services/        # API client layer
│   │   └── api.ts       # Axios configuration
│   ├── stores/          # Pinia stores
│   │   └── auth.ts      # Authentication state
│   ├── types/           # TypeScript types
│   ├── views/           # Page components
│   ├── App.vue          # Root component
│   └── main.ts          # Application entry
├── index.html
├── vite.config.ts
├── Dockerfile
└── README.md
```

## Authentication Flow

1. User logs in → store credentials in Pinia store
2. Axios interceptor adds `Authorization: Bearer <token>` to requests
3. On 401 response → attempt token refresh
4. Refresh success → retry original request
5. Refresh failure → redirect to login

## API Service Layer

All API calls go through `src/services/api.ts`:
- Base URL: `/api/v1` (proxied in dev to gateway:8080)
- Request interceptor: Adds auth token
- Response interceptor: Handles token refresh

## State Management

### Auth Store (`stores/auth.ts`)
- `token` - Access token
- `refreshToken` - Refresh token
- `user` - Current user object
- `isAuthenticated` - Computed property

## Development

```bash
cd frontend
npm install      # Install dependencies
npm run dev      # Start dev server (:3000)
npm run build    # Production build
```

## Docker

```bash
docker build -t frontend:latest .
# Runs nginx serving built static files on :80
```

## Routing

Protected routes require authentication. Redirect to `/login` if not authenticated.
