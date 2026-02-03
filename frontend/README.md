# Frontend - Vue 3 Application

A modern Vue 3 frontend application that demonstrates microservice architecture integration through an API Gateway (BFF).

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe development
- **Vite** - Next-generation build tool
- **Vue Router** - Official router for Vue.js
- **Pinia** - State management library
- **Axios** - HTTP client with interceptors

## Features

### Authentication
- JWT-based authentication with access and refresh tokens
- Automatic token refresh on 401 responses
- Protected routes with authentication guards
- Login/logout functionality

### User Management
- User listing with pagination and search
- Create new users
- Edit user details
- Delete users
- View user profiles

### User Profile
- View and edit personal profile
- Device management (view and revoke devices)
- Real-time profile updates

### Architecture
The frontend follows a clean architecture pattern:

```
src/
├── assets/          # Static assets and global styles
├── components/      # Reusable Vue components
├── router/          # Vue Router configuration
├── services/        # API service layer
├── stores/          # Pinia state management
├── views/           # Page-level components
├── App.vue          # Root component
└── main.ts          # Application entry point
```

## Local Development

### Prerequisites

- Node.js 20+
- npm or yarn

### Installation

```bash
cd frontend
npm install
```

### Running in Development Mode

```bash
npm run dev
```

The application will be available at `http://localhost:3000`

API requests will be proxied to the API Gateway via Vite's proxy configuration.

### Building for Production

```bash
npm run build
```

The built assets will be in the `dist/` directory.

### Type Checking

```bash
npm run type-check
```

## API Integration

The frontend communicates with the backend through the API Gateway:

- Base URL: `/api/v1` (proxied to API Gateway in development)
- All requests include JWT bearer tokens
- Automatic token refresh on expiration

### API Endpoints Used

**Authentication:**
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/validate` - Validate access token
- `POST /api/v1/auth/verify-pin` - Verify PIN

**Users:**
- `GET /api/v1/users` - List users
- `GET /api/v1/users/{uid}` - Get user details
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/{uid}` - Update user
- `DELETE /api/v1/users/{uid}` - Delete user
- `GET /api/v1/users/{uid}/profile` - Get user profile
- `PUT /api/v1/users/{uid}/profile` - Update user profile
- `GET /api/v1/users/{uid}/devices` - List user devices
- `DELETE /api/v1/users/{uid}/devices/{deviceUid}` - Revoke device

## State Management

### Auth Store
Manages authentication state:
- `accessToken` - Current JWT access token
- `refreshToken` - Current refresh token
- `user` - Authenticated user information
- `isAuthenticated` - Computed auth status

### Actions
- `login(identifier, identifier_type, password)` - Authenticate user
- `validateToken()` - Validate current token
- `refreshAccessToken()` - Refresh expired token
- `logout()` - Clear auth state

## Docker Deployment

The frontend is containerized using nginx for production:

```bash
docker build -t microservice-frontend .
docker run -p 80:80 microservice-frontend
```

The Dockerfile:
1. Builds the Vue application with Vite
2. Serves static files with nginx
3. Proxies API requests to the API Gateway
4. Enables gzip compression
5. Implements SPA routing with fallback to index.html

## Environment Variables

- `VITE_API_BASE_URL` - API base URL (default: `/api/v1`)

## Component Examples

### Using the Auth Store

```typescript
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// Login
await authStore.login('username', 'username', 'password')

// Check authentication status
if (authStore.isAuthenticated) {
  console.log('User:', authStore.user)
}

// Logout
authStore.logout()
```

### Making API Calls

```typescript
import { usersApi } from '@/services/api'

// List users
const { users, meta } = await usersApi.listUsers({
  page: 1,
  limit: 20,
  query: 'search term'
})

// Create user
const { uid } = await usersApi.createUser({
  username: 'johndoe',
  email: 'john@example.com',
  password: 'securepassword'
})
```

## Browser Support

- Chrome (last 2 versions)
- Firefox (last 2 versions)
- Safari (last 2 versions)
- Edge (last 2 versions)
