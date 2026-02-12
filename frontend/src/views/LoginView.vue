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
import { ref } from 'vue'
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
</style>
