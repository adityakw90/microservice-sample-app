<template>
  <transition name="dialog">
    <div v-if="show" class="confirm-dialog-overlay" @click.self="$emit('cancel')">
      <div class="confirm-dialog">
        <div class="confirm-dialog__header">
          <div class="confirm-icon" :class="`confirm-icon--${type}`">
            <component :is="dialogIcon" />
          </div>
        </div>

        <div class="confirm-dialog__body">
          <h3 class="confirm-title">{{ title }}</h3>
          <p class="confirm-message">{{ message }}</p>
        </div>

        <div class="confirm-dialog__footer">
          <button @click="$emit('cancel')" class="confirm-btn confirm-btn--cancel">
            {{ cancelText }}
          </button>
          <button @click="$emit('confirm')" class="confirm-btn confirm-btn--confirm" :class="`confirm-btn--${type}`">
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineProps<{
  show: boolean
  title: string
  message: string
  type?: 'danger' | 'warning' | 'info'
  confirmText?: string
  cancelText?: string
}>()

defineEmits<{
  confirm: []
  cancel: []
}>()

const AlertIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0zM12 9v4M12 17h.01"/>
    </svg>
  `
}

const TrashIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M10 14v4M14 14v4"/>
    </svg>
  `
}

const InfoIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="10"/>
      <path d="M12 16v-4M12 8h.01"/>
    </svg>
  `
}

const dialogIcon = computed(() => {
  const icons = {
    danger: TrashIcon,
    warning: AlertIcon,
    info: InfoIcon
  }
  return icons[props.type] || InfoIcon
})
</script>

<style scoped>
.confirm-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.confirm-dialog {
  background: white;
  border-radius: 20px;
  padding: 2rem;
  max-width: 420px;
  width: 100%;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  animation: dialog-in 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@keyframes dialog-in {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.confirm-dialog__header {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
}

.confirm-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.confirm-icon svg {
  width: 32px;
  height: 32px;
}

.confirm-icon--danger {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(248, 113, 113, 0.1) 100%);
}

.confirm-icon--danger svg {
  color: #ef4444;
}

.confirm-icon--warning {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.15) 0%, rgba(251, 191, 36, 0.1) 100%);
}

.confirm-icon--warning svg {
  color: #f59e0b;
}

.confirm-icon--info {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15) 0%, rgba(96, 165, 250, 0.1) 100%);
}

.confirm-icon--info svg {
  color: #3b82f6;
}

.confirm-dialog__body {
  text-align: center;
  margin-bottom: 2rem;
}

.confirm-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 0.75rem 0;
}

.confirm-message {
  font-size: 0.9rem;
  color: #64748b;
  margin: 0;
  line-height: 1.5;
}

.confirm-dialog__footer {
  display: flex;
  gap: 0.75rem;
}

.confirm-btn {
  flex: 1;
  padding: 0.875rem 1.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.confirm-btn--cancel {
  background: #f1f5f9;
  color: #64748b;
}

.confirm-btn--cancel:hover {
  background: #e2e8f0;
  color: #475569;
}

.confirm-btn--confirm {
  color: white;
}

.confirm-btn--danger {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
}

.confirm-btn--danger:hover {
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
}

.confirm-btn--warning {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.confirm-btn--warning:hover {
  background: linear-gradient(135deg, #d97706 0%, #b45309 100%);
}

.confirm-btn--info {
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
}

.confirm-btn--info:hover {
  background: linear-gradient(135deg, #4f46e5 0%, #4338ca 100%);
}

.dialog-enter-active,
.dialog-leave-active {
  transition: all 0.2s ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-from .confirm-dialog,
.dialog-leave-to .confirm-dialog {
  transform: scale(0.95);
}
</style>
