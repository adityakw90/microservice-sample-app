<template>
  <div class="activity-feed">
    <div class="activity-header">
      <h3 class="activity-title">Recent Activity</h3>
      <button class="view-all-btn">View All</button>
    </div>

    <div class="activity-list">
      <transition-group name="activity" tag="div">
        <div
          v-for="(activity, index) in activities"
          :key="activity.id"
          class="activity-item"
          :style="{ transitionDelay: `${index * 50}ms` }"
        >
          <div class="activity-icon" :class="`activity-icon--${activity.type}`">
            <component :is="getActivityIcon(activity.type)" />
          </div>
          <div class="activity-content">
            <div class="activity-text">
              <span class="activity-user">{{ activity.user }}</span>
              <span class="activity-action">{{ activity.action }}</span>
            </div>
            <div class="activity-time">{{ activity.time }}</div>
          </div>
          <div v-if="activity.changes" class="activity-changes">
            <span class="change-badge">{{ activity.changes }}</span>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineProps<{
  activities: Array<{
    id: string
    user: string
    action: string
    time: string
    type: 'user' | 'create' | 'delete' | 'update' | 'security'
    changes?: string
  }>
}>()

const UserIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/>
      <circle cx="12" cy="7" r="4"/>
    </svg>
  `
}

const PlusIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="10"/>
      <path d="M12 6v12M6 12h12"/>
    </svg>
  `
}

const TrashIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
    </svg>
  `
}

const EditIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
      <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
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

const getActivityIcon = (type: string) => {
  const icons = {
    user: UserIcon,
    create: PlusIcon,
    delete: TrashIcon,
    update: EditIcon,
    security: ShieldIcon
  }
  return icons[type] || UserIcon
}
</script>

<style scoped>
.activity-feed {
  background: white;
  border-radius: 20px;
  padding: 1.5rem;
  border: 1px solid rgba(0, 0, 0, 0.04);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.activity-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.view-all-btn {
  padding: 0.5rem 1rem;
  font-size: 0.8rem;
  font-weight: 500;
  color: #6366f1;
  background: transparent;
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.view-all-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  border-color: #6366f1;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 0;
  border-bottom: 1px solid #f1f5f9;
  transition: all 0.2s ease;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-item:hover {
  background: rgba(99, 102, 241, 0.03);
  margin: 0 -1rem;
  padding-left: 1rem;
  padding-right: 1rem;
  border-radius: 10px;
}

.activity-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.activity-icon svg {
  width: 18px;
  height: 18px;
}

.activity-icon--user {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(139, 92, 246, 0.1) 100%);
  color: #6366f1;
}

.activity-icon--create {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(52, 211, 153, 0.1) 100%);
  color: #10b981;
}

.activity-icon--delete {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(248, 113, 113, 0.1) 100%);
  color: #ef4444;
}

.activity-icon--update {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.15) 0%, rgba(251, 191, 36, 0.1) 100%);
  color: #f59e0b;
}

.activity-icon--security {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.15) 0%, rgba(167, 139, 250, 0.1) 100%);
  color: #8b5cf6;
}

.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-text {
  font-size: 0.9rem;
  color: #1e293b;
  margin-bottom: 0.25rem;
}

.activity-user {
  font-weight: 600;
}

.activity-action {
  color: #64748b;
}

.activity-time {
  font-size: 0.75rem;
  color: #94a3b8;
}

.activity-changes {
  margin-left: auto;
}

.change-badge {
  padding: 0.25rem 0.625rem;
  font-size: 0.7rem;
  font-weight: 600;
  background: #f1f5f9;
  color: #64748b;
  border-radius: 12px;
}

.activity-enter-active {
  transition: all 0.3s ease;
}

.activity-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
