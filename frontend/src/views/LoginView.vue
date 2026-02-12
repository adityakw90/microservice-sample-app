<template>
  <div class="login-view">
    <div class="login-container">
      <div class="login-card">
        <h1>Welcome Back</h1>
        <p class="subtitle">Sign in to your account</p>

        <form @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label for="identifier">Username or Email</label>
            <input
              id="identifier"
              v-model="formData.identifier"
              type="text"
              required
              class="form-input"
              placeholder="Enter your username or email"
            />
          </div>

          <div class="form-group">
            <label for="password">Password</label>
            <input
              id="password"
              v-model="formData.password"
              type="password"
              required
              class="form-input"
              placeholder="Enter your password"
            />
          </div>

          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="login-btn"
          >
            {{ loading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>

        <div class="oauth-divider">
          <span>or</span>
        </div>

        <button
          type="button"
          @click="handleGoogleOAuth"
          class="google-btn"
        >
          <svg class="google-icon" viewBox="0 0 24 24">
            <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.21-3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23c2.97 0 5.46-.98 7.28-2.66l2.85-2.22.81-.62c.87-2.6 3.3-4.53z"/>
            <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53z"/>
          </svg>
          Continue with Google
        </button>

        <div class="login-footer">
          <p>Don't have an account?</p>
          <router-link to="/users" class="register-link">Browse Users</router-link>
        </div>
      </div>

      <div class="info-card">
        <h2>About This Demo</h2>
        <p>This is a showcase application demonstrating:</p>
        <ul>
          <li>Vue 3 frontend with TypeScript</li>
          <li>API Gateway (BFF) pattern</li>
          <li>gRPC microservice communication</li>
          <li>JWT-based authentication</li>
          <li>User management features</li>
        </ul>
        <div class="tech-stack">
          <h3>Tech Stack</h3>
          <div class="tech-tags">
            <span class="tag">Vue 3</span>
            <span class="tag">TypeScript</span>
            <span class="tag">Pinia</span>
            <span class="tag">Axios</span>
            <span class="tag">Go</span>
            <span class="tag">gRPC</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const formData = ref({
  identifier: '',
  password: ''
})

const loading = ref(false)
const error = ref<string | null>(null)

const handleLogin = async () => {
  loading.value = true
  error.value = null

  try {
    const identifier = formData.value.identifier
    const identifier_type = identifier.includes('@') ? 'email' : 'username'

    const success = await authStore.login(
      identifier,
      identifier_type,
      formData.value.password
    )

    if (success) {
      router.push('/profile')
    } else {
      error.value = 'Invalid credentials'
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}

const handleGoogleOAuth = () => {
  // Simply redirect to gateway OAuth endpoint
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1'
  window.location.href = `${apiBaseUrl}/auth/google`
}

onMounted(() => {
  // Check for OAuth callback tokens in URL
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token')
  const refreshToken = urlParams.get('refresh')

  if (token && refreshToken) {
    // Store tokens from OAuth callback
    localStorage.setItem('access_token', token)
    localStorage.setItem('refresh_token', refreshToken)

    // Update store
    authStore.accessToken = token
    authStore.refreshToken = refreshToken

    // Validate token to get user info
    authStore.validateToken(token).then(() => {
      // Clean URL and redirect to profile
      window.history.replaceState({}, '', '/login')
      router.push('/profile')
    })
  }
})
</script>

<style scoped>
.login-view {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.login-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  max-width: 1000px;
  width: 100%;
}

.login-card {
  background: white;
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
}

.login-card h1 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
  font-size: 2rem;
}

.subtitle {
  color: #7f8c8d;
  margin-bottom: 2rem;
}

.login-form {
  margin-top: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2c3e50;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 0.875rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #3498db;
}

.error-message {
  background-color: #fee;
  color: #c33;
  padding: 0.75rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.login-btn {
  width: 100%;
  padding: 0.875rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.login-footer {
  margin-top: 2rem;
  text-align: center;
  color: #7f8c8d;
}

.register-link {
  color: #3498db;
  text-decoration: none;
  font-weight: 500;
  display: block;
  margin-top: 0.5rem;
}

.register-link:hover {
  text-decoration: underline;
}

.info-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2.5rem;
  border-radius: 12px;
}

.info-card h2 {
  margin: 0 0 1rem 0;
  font-size: 1.75rem;
}

.info-card p {
  margin-bottom: 1rem;
  opacity: 0.9;
}

.info-card ul {
  list-style: none;
  padding: 0;
  margin-bottom: 2rem;
}

.info-card li {
  padding: 0.5rem 0;
  display: flex;
  align-items: center;
}

.info-card li::before {
  content: '✓';
  margin-right: 0.75rem;
  font-weight: bold;
}

.tech-stack h3 {
  margin: 0 0 1rem 0;
  font-size: 1.25rem;
}

.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  background: rgba(255,255,255,0.2);
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
}

@media (max-width: 768px) {
  .login-container {
    grid-template-columns: 1fr;
  }
}

.oauth-divider {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
  color: #7f8c8d;
}

.oauth-divider::before,
.oauth-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #ddd;
}

.oauth-divider span {
  padding: 0 1rem;
  font-size: 0.875rem;
}

.google-btn {
  width: 100%;
  padding: 0.875rem;
  background: white;
  color: #2c3e50;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  transition: background-color 0.2s, border-color 0.2s;
}

.google-btn:hover {
  background: #f8f9fa;
  border-color: #ccc;
}

.google-icon {
  width: 20px;
  height: 20px;
}

.google-icon path:nth-child(1) {
  fill: #4285F4;
}

.google-icon path:nth-child(2) {
  fill: #34A853;
}

.google-icon path:nth-child(3) {
  fill: #FBBC05;
}

.google-icon path:nth-child(4) {
  fill: #EA4335;
}
</style>
