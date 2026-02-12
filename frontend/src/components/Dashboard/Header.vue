<template>
  <header class="dashboard-header">
    <div class="header-left">
      <button class="mobile-menu-btn" @click="$emit('toggle-sidebar')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 12h18M3 6h18M3 18h18"/>
        </svg>
      </button>
      <div class="search-container">
        <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <path d="M21 21l-4.35-4.35"/>
        </svg>
        <input
          ref="searchInput"
          v-model="searchQuery"
          type="text"
          placeholder="Search users, settings..."
          class="search-input"
          @focus="showSearchResults = true"
          @blur="hideSearchResults"
        />
        <kbd v-if="!searchQuery" class="search-shortcut">⌘K</kbd>
        <transition name="search-results">
          <div v-if="showSearchResults && searchQuery" class="search-results">
            <div v-if="searchResults.length === 0" class="search-no-results">
              No results found
            </div>
            <a
              v-for="result in searchResults"
              :key="result.id"
              href="#"
              class="search-result-item"
              @click.prevent="handleSearchResult(result)"
            >
              <component :is="result.icon" class="result-icon" />
              <div class="result-content">
                <div class="result-title">{{ result.title }}</div>
                <div class="result-subtitle">{{ result.subtitle }}</div>
              </div>
            </a>
          </div>
        </transition>
      </div>
    </div>

    <div class="header-right">
      <button class="header-btn theme-toggle" @click="$emit('toggle-theme')">
        <transition name="icon-swap" mode="out-in">
          <svg v-if="!isDarkMode" key="moon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
          </svg>
          <svg v-else key="sun" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5"/>
            <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
          </svg>
        </transition>
      </button>

      <button class="header-btn notification-btn">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 01-3.46 0"/>
        </svg>
        <span class="notification-badge">3</span>
      </button>

      <div class="header-divider"></div>

      <div class="user-menu">
        <button class="user-menu-btn">
          <div class="user-avatar">
            <span>AD</span>
          </div>
          <svg class="chevron-down" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M6 9l6 6 6-6"/>
          </svg>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

defineProps<{
  isDarkMode: boolean
}>()

defineEmits<{
  toggleSidebar: []
  toggleTheme: []
}>()

const router = useRouter()
const searchQuery = ref('')
const showSearchResults = ref(false)
const searchInput = ref<HTMLInputElement>()

const UserIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/>
      <circle cx="12" cy="7" r="4"/>
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

const searchResults = computed(() => {
  if (!searchQuery.value) return []

  const results = []
  const query = searchQuery.value.toLowerCase()

  // Search users
  if (query.includes('user') || query.includes('profile')) {
    results.push({
      id: 'users',
      title: 'Users',
      subtitle: 'Manage user accounts',
      icon: UserIcon,
      path: '/users'
    })
  }

  // Search settings
  if (query.includes('setting') || query.includes('config')) {
    results.push({
      id: 'settings',
      title: 'Settings',
      subtitle: 'Application settings',
      icon: SettingsIcon,
      path: '/settings'
    })
  }

  return results
})

const handleSearchResult = (result: any) => {
  router.push(result.path)
  showSearchResults.value = false
  searchQuery.value = ''
}

const hideSearchResults = () => {
  setTimeout(() => {
    showSearchResults.value = false
  }, 200)
}

const handleKeydown = (e: KeyboardEvent) => {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    searchInput.value?.focus()
  }
  if (e.key === 'Escape') {
    showSearchResults.value = false
    searchQuery.value = ''
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.dashboard-header {
  height: 70px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  position: sticky;
  top: 0;
  z-index: 50;
  transition: all 0.3s ease;
}

:deep(.dark-mode) .dashboard-header {
  background: rgba(15, 23, 42, 0.8);
  border-bottom-color: rgba(255, 255, 255, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  flex: 1;
}

.mobile-menu-btn {
  display: none;
  width: 40px;
  height: 40px;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  color: #64748b;
  transition: all 0.2s ease;
}

.mobile-menu-btn:hover {
  background: #f1f5f9;
  color: #334155;
}

.mobile-menu-btn svg {
  width: 20px;
  height: 20px;
}

.search-container {
  position: relative;
  max-width: 400px;
  width: 100%;
}

.search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  color: #94a3b8;
  pointer-events: none;
}

.search-input {
  width: 100%;
  height: 44px;
  padding: 0 3.5rem 0 2.75rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  font-size: 0.9rem;
  color: #1e293b;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  background: white;
  border-color: #6366f1;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.1);
}

.search-input::placeholder {
  color: #94a3b8;
}

.search-shortcut {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  padding: 0.25rem 0.5rem;
  font-size: 0.7rem;
  font-weight: 500;
  color: #94a3b8;
  background: #e2e8f0;
  border-radius: 6px;
  pointer-events: none;
}

.search-results {
  position: absolute;
  top: calc(100% + 0.5rem);
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  z-index: 100;
}

.search-no-results {
  padding: 1.5rem;
  text-align: center;
  color: #94a3b8;
  font-size: 0.9rem;
}

.search-result-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  text-decoration: none;
  color: inherit;
  transition: background 0.15s ease;
}

.search-result-item:hover {
  background: #f8fafc;
}

.result-icon {
  width: 18px;
  height: 18px;
  color: #6366f1;
  flex-shrink: 0;
}

.result-content {
  flex: 1;
  min-width: 0;
}

.result-title {
  font-size: 0.9rem;
  font-weight: 500;
  color: #1e293b;
}

.result-subtitle {
  font-size: 0.75rem;
  color: #64748b;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.header-btn {
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  color: #64748b;
  transition: all 0.2s ease;
  position: relative;
}

.header-btn:hover {
  background: #f1f5f9;
  color: #334155;
}

.header-btn svg {
  width: 20px;
  height: 20px;
}

.notification-btn .notification-badge {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 18px;
  height: 18px;
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  border: 2px solid white;
  border-radius: 50%;
  font-size: 0.65rem;
  font-weight: 600;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-divider {
  width: 1px;
  height: 24px;
  background: #e2e8f0;
  margin: 0 0.5rem;
}

.user-menu-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem;
  background: transparent;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.user-menu-btn:hover {
  background: #f1f5f9;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 600;
  color: white;
}

.chevron-down {
  width: 16px;
  height: 16px;
  color: #64748b;
}

.icon-swap-enter-active,
.icon-swap-leave-active {
  transition: all 0.2s ease;
}

.icon-swap-enter-from {
  opacity: 0;
  transform: rotate(-90deg) scale(0.5);
}

.icon-swap-leave-to {
  opacity: 0;
  transform: rotate(90deg) scale(0.5);
}

.search-results-enter-active,
.search-results-leave-active {
  transition: all 0.2s ease;
}

.search-results-enter-from,
.search-results-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 768px) {
  .dashboard-header {
    padding: 0 1rem;
  }

  .mobile-menu-btn {
    display: flex;
  }

  .search-container {
    max-width: 200px;
  }

  .search-shortcut {
    display: none;
  }
}
</style>
