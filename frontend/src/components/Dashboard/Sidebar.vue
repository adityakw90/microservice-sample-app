<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-header">
      <transition name="logo">
        <div v-if="!collapsed" class="logo-container">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="url(#logo-gradient-1)"/>
              <path d="M2 17L12 22L22 17" stroke="url(#logo-gradient-2)" stroke-width="2" stroke-linecap="round"/>
              <path d="M2 12L12 17L22 12" stroke="url(#logo-gradient-2)" stroke-width="2" stroke-linecap="round"/>
              <defs>
                <linearGradient id="logo-gradient-1" x1="2" y1="2" x2="22" y2="12">
                  <stop offset="0%" stop-color="#6366f1"/>
                  <stop offset="100%" stop-color="#8b5cf6"/>
                </linearGradient>
                <linearGradient id="logo-gradient-2" x1="2" y1="7" x2="22" y2="22">
                  <stop offset="0%" stop-color="#6366f1"/>
                  <stop offset="100%" stop-color="#8b5cf6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <span class="logo-text">MicroAdmin</span>
        </div>
      </transition>
      <button class="toggle-btn" @click="$emit('toggle')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 19l-7-7 7-7M18 19l-7-7 7-7"/>
        </svg>
      </button>
    </div>

    <nav class="sidebar-nav">
      <div v-for="section in menuSections" :key="section.title" class="nav-section">
        <transition name="section-title">
          <div v-if="!collapsed" class="section-title">{{ section.title }}</div>
        </transition>
        <router-link
          v-for="item in section.items"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
        >
          <component :is="item.icon" class="nav-icon" />
          <transition name="label">
            <span v-if="!collapsed" class="nav-label">{{ item.label }}</span>
          </transition>
          <transition name="badge">
            <span v-if="!collapsed && item.badge" class="nav-badge">{{ item.badge }}</span>
          </transition>
        </router-link>
      </div>
    </nav>

    <div class="sidebar-footer">
      <div class="footer-gradient"></div>
      <transition name="user-info">
        <div v-if="!collapsed" class="user-info">
          <div class="user-avatar">
            <span>AD</span>
          </div>
          <div class="user-details">
            <div class="user-name">Admin User</div>
            <div class="user-status">Online</div>
          </div>
        </div>
      </transition>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

defineProps<{
  collapsed: boolean
}>()

defineEmits<{
  toggle: []
}>()

const route = useRoute()

const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/')
}

// Icon components
const HomeIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
      <path d="M9 22V12h6v10"/>
    </svg>
  `
}

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

const SettingsIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="3"/>
      <path d="M12 1v6M12 17v6M4.22 4.22l4.24 4.24M15.54 15.54l4.24 4.24M1 12h6M17 12h6M4.22 19.78l4.24-4.24M15.54 8.46l4.24-4.24"/>
    </svg>
  `
}

const ShieldIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
    </svg>
  `
}

const TestIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polyline points="9,11 12,14 22,4"/>
      <path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/>
    </svg>
  `
}

const FilesIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"/>
    </svg>
  `
}

const menuSections = computed(() => [
  {
    title: 'Overview',
    items: [
      { path: '/dashboard', label: 'Dashboard', icon: HomeIcon },
      { path: '/users', label: 'Users', icon: UsersIcon, badge: 'New' }
    ]
  },
  {
    title: 'Management',
    items: [
      { path: '/files', label: 'Files', icon: FilesIcon },
      { path: '/activity', label: 'Activity', icon: ActivityIcon },
      { path: '/security', label: 'Security', icon: ShieldIcon },
      { path: '/settings', label: 'Settings', icon: SettingsIcon }
    ]
  },
  {
    title: 'Developer',
    items: [
      { path: '/proto-test', label: 'Proto Test', icon: TestIcon, badge: 'Dev' }
    ]
  }
])
</script>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  width: 260px;
  background: linear-gradient(180deg, #1e293b 0%, #0f172a 100%);
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 100;
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.15);
}

.sidebar.collapsed {
  width: 72px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 700;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.toggle-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  color: #94a3b8;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.toggle-btn svg {
  width: 18px;
  height: 18px;
  transition: transform 0.3s ease;
}

.sidebar.collapsed .toggle-btn svg {
  transform: rotate(180deg);
}

.sidebar-nav {
  flex: 1;
  padding: 1rem 0.75rem;
  overflow-y: auto;
  overflow-x: hidden;
}

.nav-section {
  margin-bottom: 1.5rem;
}

.section-title {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #64748b;
  padding: 0.5rem 0.75rem;
  margin-bottom: 0.5rem;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 10px;
  color: #94a3b8;
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
  margin-bottom: 0.25rem;
}

.nav-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 0;
  background: linear-gradient(180deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 0 3px 3px 0;
  transition: height 0.2s ease;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e2e8f0;
}

.nav-item.active {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.15) 0%, rgba(139, 92, 246, 0.05) 100%);
  color: #e2e8f0;
}

.nav-item.active::before {
  height: 24px;
}

.nav-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.nav-label {
  font-size: 0.9rem;
  font-weight: 500;
  white-space: nowrap;
}

.nav-badge {
  margin-left: auto;
  padding: 0.15rem 0.5rem;
  font-size: 0.65rem;
  font-weight: 600;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border-radius: 10px;
}

.sidebar-footer {
  position: relative;
  padding: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.footer-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, #6366f1 50%, transparent 100%);
  opacity: 0.5;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  font-weight: 600;
  color: white;
  flex-shrink: 0;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 0.85rem;
  font-weight: 600;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-status {
  font-size: 0.7rem;
  color: #10b981;
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.user-status::before {
  content: '';
  width: 6px;
  height: 6px;
  background: #10b981;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.logo-enter-active,
.logo-leave-active,
.section-title-enter-active,
.section-title-leave-active,
.label-enter-active,
.label-leave-active,
.badge-enter-active,
.badge-leave-active,
.user-info-enter-active,
.user-info-leave-active {
  transition: all 0.2s ease;
}

.logo-enter-from,
.logo-leave-to,
.section-title-enter-from,
.section-title-leave-to,
.label-enter-from,
.label-leave-to,
.badge-enter-from,
.badge-leave-to,
.user-info-enter-from,
.user-info-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

.sidebar.collapsed .logo-enter-active,
.sidebar.collapsed .label-enter-active,
.sidebar.collapsed .badge-enter-active,
.sidebar.collapsed .section-title-enter-active,
.sidebar.collapsed .user-info-enter-active {
  position: absolute;
}
</style>
