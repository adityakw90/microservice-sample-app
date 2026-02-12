<template>
  <div class="home">
    <div class="hero">
      <h1>Welcome to the Microservice Showcase</h1>
      <p class="subtitle">A demonstration of modern microservice architecture</p>
    </div>

    <div class="features">
      <div class="feature-card">
        <h3>🔐 Authentication</h3>
        <p>Secure JWT-based authentication with access and refresh tokens.</p>
        <router-link to="/login" class="btn">Try Login</router-link>
      </div>

      <div class="feature-card">
        <h3>👥 User Management</h3>
        <p>Complete CRUD operations for user accounts with profile management.</p>
        <router-link to="/users" class="btn">View Users</router-link>
      </div>

      <div class="feature-card">
        <h3>🔧 gRPC Services</h3>
        <p>Backend communication via gRPC with a Go-based microservice.</p>
        <a href="https://grpc.io" target="_blank" class="btn">Learn More</a>
      </div>
    </div>

    <div class="architecture">
      <h2>Architecture Overview</h2>
      <div class="arch-diagram">
        <div class="arch-box frontend">
          <h4>Frontend</h4>
          <p>Vue 3 Application</p>
        </div>
        <div class="arch-arrow">↓ HTTP/REST</div>
        <div class="arch-box gateway">
          <h4>API Gateway</h4>
          <p>Backend For Frontend (BFF)</p>
        </div>
        <div class="arch-arrow">↓ gRPC</div>
        <div class="arch-box service">
          <h4>User Service</h4>
          <p>Go gRPC Microservice</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home {
  max-width: 1000px;
  margin: 0 auto;
}

.hero {
  text-align: center;
  margin: 3rem 0;
}

.hero h1 {
  font-size: 2.5rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.subtitle {
  font-size: 1.2rem;
  color: #7f8c8d;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
  margin: 3rem 0;
}

.feature-card {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  text-align: center;
}

.feature-card h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #2c3e50;
}

.feature-card p {
  color: #7f8c8d;
  margin-bottom: 1.5rem;
}

.btn {
  display: inline-block;
  background-color: #3498db;
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  text-decoration: none;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #2980b9;
}

.architecture {
  margin: 4rem 0;
  text-align: center;
}

.architecture h2 {
  color: #2c3e50;
  margin-bottom: 2rem;
}

.arch-diagram {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.arch-box {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1.5rem 2rem;
  border-radius: 8px;
  min-width: 250px;
}

.arch-box.frontend {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.arch-box.gateway {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.arch-box.service {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.arch-box h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.2rem;
}

.arch-box p {
  margin: 0;
  opacity: 0.9;
}

.arch-arrow {
  color: #2c3e50;
  font-size: 1.5rem;
  font-weight: bold;
}
</style>
