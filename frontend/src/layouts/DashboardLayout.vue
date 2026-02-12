<template>
  <div class="dashboard-layout" :class="{ 'dark-mode': isDarkMode }">
    <DashboardSidebar
      :collapsed="sidebarCollapsed"
      @toggle="sidebarCollapsed = !sidebarCollapsed"
    />

    <div class="main-wrapper" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
      <DashboardHeader
        @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed"
        @toggle-theme="toggleTheme"
        :is-dark-mode="isDarkMode"
      />

      <main class="dashboard-main">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <NotificationToast
      v-if="notificationStore.notification.show"
      :message="notificationStore.notification.message"
      :type="notificationStore.notification.type"
      @close="notificationStore.hide"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DashboardSidebar from '@/components/Dashboard/Sidebar.vue'
import DashboardHeader from '@/components/Dashboard/Header.vue'
import NotificationToast from '@/components/common/NotificationToast.vue'
import { useNotificationStore } from '@/stores/notification'

const sidebarCollapsed = ref(false)
const isDarkMode = ref(false)
const notificationStore = useNotificationStore()

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  document.documentElement.classList.toggle('dark', isDarkMode.value)
  localStorage.setItem('darkMode', String(isDarkMode.value))
}

const showNotification = (message: string, type: 'success' | 'error' | 'warning' | 'info') => {
  notificationStore.show(message, type)
}

onMounted(() => {
  const savedDarkMode = localStorage.getItem('darkMode')
  if (savedDarkMode === 'true') {
    isDarkMode.value = true
    document.documentElement.classList.add('dark')
  }
})

defineExpose({
  showNotification
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap');

:where(*) {
  --font-display: 'Outfit', -apple-system, BlinkMacSystemFont, sans-serif;
  --font-mono: 'JetBrains Mono', 'SF Mono', monospace;
}

.dashboard-layout {
  display: flex;
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  font-family: var(--font-display);
  transition: background 0.3s ease;
}

.dashboard-layout.dark-mode {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}

.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 260px;
  transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-wrapper.sidebar-collapsed {
  margin-left: 72px;
}

.dashboard-main {
  flex: 1;
  padding: 2rem;
  min-height: calc(100vh - 70px);
}

@media (max-width: 768px) {
  .main-wrapper {
    margin-left: 0;
  }

  .main-wrapper.sidebar-collapsed {
    margin-left: 0;
  }

  .dashboard-main {
    padding: 1rem;
  }
}

.page-enter-active,
.page-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>
