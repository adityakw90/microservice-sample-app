<template>
  <div class="dashboard-view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__content">
        <h1 class="page-title">Dashboard</h1>
        <p class="page-subtitle">Welcome back! Here's what's happening with your application.</p>
      </div>
      <div class="page-header__actions">
        <button @click="refreshData" class="refresh-btn" :class="{ 'refresh-btn--spinning': refreshing }">
          <svg class="refresh-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 00-9-9 9.75 9.75 0 00-6.74 2.74L3 8"/>
            <path d="M3 3v5h5M3 12a9 9 0 009 9 9.75 9.75 0 006.74-2.74L21 16"/>
            <path d="M16 21h5v-5"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="dashboard-loading">
      <LoadingSpinner text="Loading dashboard..." />
    </div>

    <!-- Dashboard Content -->
    <div v-else class="dashboard-content">
        <!-- Stats Grid -->
        <div class="stats-grid">
          <StatCard
            label="Total Users"
            :value="stats.totalUsers"
            suffix="users"
            :icon="UsersIcon"
            color="indigo"
            trend="+12.5%"
            :trend-up="true"
            :progress="totalUsersProgress"
            helper="of 5,000 capacity"
          />

          <StatCard
            label="Active Users"
            :value="stats.activeUsers"
            suffix="online"
            :icon="ActivityIcon"
            color="emerald"
            trend="+8.2%"
            :trend-up="true"
            :progress="activeUsersProgress"
            helper="currently active"
          />

          <StatCard
            label="New Signups"
            :value="stats.newRegistrations"
            suffix="today"
            :icon="TrendingUpIcon"
            color="amber"
            trend="+23.1%"
            :trend-up="true"
            :progress="newRegistrationsProgress"
            helper="vs yesterday"
          />

          <StatCard
            label="Active Sessions"
            :value="stats.activeSessions"
            suffix="sessions"
            :icon="ServerIcon"
            color="rose"
            trend="-2.4%"
            :trend-up="false"
            :progress="activeSessionsProgress"
            helper="across all servers"
          />
        </div>

        <!-- Main Grid -->
        <div class="dashboard-grid">
          <!-- Chart Section -->
          <div class="dashboard-card dashboard-card--chart">
            <ChartComponent
              title="User Registrations"
              subtitle="New user signups over the last 7 days"
              :data="chartData"
            />
          </div>

          <!-- Activity Feed -->
          <div class="dashboard-card dashboard-card--activity">
            <ActivityFeed :activities="activities" />
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="quick-actions">
          <div class="quick-actions__header">
            <h2 class="quick-actions__title">Quick Actions</h2>
          </div>
          <div class="quick-actions__grid">
            <router-link to="/users" class="quick-action">
              <div class="quick-action__icon quick-action__icon--users">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"/>
                </svg>
              </div>
              <div class="quick-action__content">
                <div class="quick-action__title">Manage Users</div>
                <div class="quick-action__desc">View and manage user accounts</div>
              </div>
              <svg class="quick-action__arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </router-link>

            <router-link to="/settings" class="quick-action">
              <div class="quick-action__icon quick-action__icon--settings">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3"/>
                  <path d="M12 1v6M12 17v6M4.22 4.22l4.24 4.24M15.54 15.54l4.24 4.24M1 12h6M17 12h6M4.22 19.78l4.24-4.24M15.54 8.46l4.24-4.24"/>
                </svg>
              </div>
              <div class="quick-action__content">
                <div class="quick-action__title">Settings</div>
                <div class="quick-action__desc">Configure application settings</div>
              </div>
              <svg class="quick-action__arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </router-link>

            <a href="#" class="quick-action">
              <div class="quick-action__icon quick-action__icon--security">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                </svg>
              </div>
              <div class="quick-action__content">
                <div class="quick-action__title">Security</div>
                <div class="quick-action__desc">Review security logs</div>
              </div>
              <svg class="quick-action__arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </a>

            <a href="#" class="quick-action">
              <div class="quick-action__icon quick-action__icon--reports">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                  <path d="M14 2v6h6M16 13H8M16 17H8M10 9H8"/>
                </svg>
              </div>
              <div class="quick-action__content">
                <div class="quick-action__title">Reports</div>
                <div class="quick-action__desc">View analytics reports</div>
              </div>
              <svg class="quick-action__arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </a>
          </div>
        </div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useDashboard } from '@/composables/useDashboard'
import StatCard from '@/components/Dashboard/StatCard.vue'
import ChartComponent from '@/components/Dashboard/ChartComponent.vue'
import ActivityFeed from '@/components/Dashboard/ActivityFeed.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

const {
  loading,
  stats,
  activities,
  chartData,
  totalUsersProgress,
  activeUsersProgress,
  newRegistrationsProgress,
  activeSessionsProgress,
  loadDashboardData
} = useDashboard()

const refreshing = ref(false)

const refreshData = async () => {
  refreshing.value = true
  await loadDashboardData()
  setTimeout(() => {
    refreshing.value = false
  }, 500)
}

// Icon components
const UsersIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
      <circle cx="9" cy="7" r="4"/>
      <path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"/>
    </svg>
  `
}

const ActivityIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polyline points="22,12 18,12 15,21 9,3 6,12 2,12"/>
    </svg>
  `
}

const TrendingUpIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polyline points="23,6 13.5,15.5 8.5,10.5 1,18"/>
      <polyline points="17,6 23,6 23,12"/>
    </svg>
  `
}

const ServerIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <rect x="2" y="2" width="20" height="8" rx="2" ry="2"/>
      <rect x="2" y="14" width="20" height="8" rx="2" ry="2"/>
      <line x1="6" y1="6" x2="6.01" y2="6"/>
      <line x1="6" y1="18" x2="6.01" y2="18"/>
    </svg>
  `
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.dashboard-view {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2.5rem;
  gap: 1.5rem;
}

.page-header__content {
  flex: 1;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 0.5rem 0;
  font-family: var(--font-display);
}

.page-subtitle {
  font-size: 1rem;
  color: #64748b;
  margin: 0;
}

.page-header__actions {
  display: flex;
  gap: 0.75rem;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  font-size: 0.9rem;
  font-weight: 500;
  color: #6366f1;
  background: white;
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.refresh-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  border-color: #6366f1;
}

.refresh-icon {
  width: 18px;
  height: 18px;
  transition: transform 0.5s ease;
}

.refresh-btn--spinning .refresh-icon {
  animation: spin 0.5s linear;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.dashboard-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.5rem;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

@media (max-width: 1024px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

.dashboard-card {
  background: white;
  border-radius: 20px;
  padding: 1.5rem;
  border: 1px solid rgba(0, 0, 0, 0.04);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
}

.quick-actions {
  background: white;
  border-radius: 20px;
  padding: 1.5rem;
  border: 1px solid rgba(0, 0, 0, 0.04);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
}

.quick-actions__header {
  margin-bottom: 1.5rem;
}

.quick-actions__title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.quick-actions__grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 0.75rem;
}

.quick-action {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
}

.quick-action::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.05) 0%, rgba(139, 92, 246, 0.02) 100%);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.quick-action:hover::before {
  opacity: 1;
}

.quick-action:hover {
  border-color: #6366f1;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.15);
}

.quick-action__icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.quick-action__icon svg {
  width: 20px;
  height: 20px;
}

.quick-action__icon--users {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
}

.quick-action__icon--settings {
  background: linear-gradient(135deg, #64748b 0%, #94a3b8 100%);
  color: white;
}

.quick-action__icon--security {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: white;
}

.quick-action__icon--reports {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  color: white;
}

.quick-action__content {
  flex: 1;
  min-width: 0;
  position: relative;
  z-index: 1;
}

.quick-action__title {
  font-size: 0.9rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 0.25rem;
}

.quick-action__desc {
  font-size: 0.8rem;
  color: #64748b;
}

.quick-action__arrow {
  width: 18px;
  height: 18px;
  color: #cbd5e1;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
  transition: color 0.2s ease;
}

.quick-action:hover .quick-action__arrow {
  color: #6366f1;
}

.content-enter-active,
.content-leave-active {
  transition: all 0.3s ease;
}

.content-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.content-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .quick-actions__grid {
    grid-template-columns: 1fr;
  }
}
</style>
