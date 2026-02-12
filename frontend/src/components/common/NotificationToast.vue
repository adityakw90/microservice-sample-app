<template>
  <transition name="toast">
    <div class="notification-toast" :class="`notification-toast--${type}`">
      <div class="toast-icon-wrapper">
        <component :is="toastIcon" class="toast-icon" />
      </div>
      <div class="toast-content">
        <div class="toast-message">{{ message }}</div>
      </div>
      <button @click="$emit('close')" class="toast-close">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6L6 18M6 6l12 12"/>
        </svg>
      </button>
      <div class="toast-progress"></div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineProps<{
  message: string
  type: 'success' | 'error' | 'warning' | 'info'
}>()

defineEmits<{
  close: []
}>()

const CheckIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
      <path d="M20 6L9 17l-5-5"/>
    </svg>
  `
}

const ErrorIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
      <circle cx="12" cy="12" r="10"/>
      <path d="M15 9l-6 6M9 9l6 6"/>
    </svg>
  `
}

const WarningIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
      <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0zM12 9v4M12 17h.01"/>
    </svg>
  `
}

const InfoIcon = {
  template: `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
      <circle cx="12" cy="12" r="10"/>
      <path d="M12 16v-4M12 8h.01"/>
    </svg>
  `
}

const toastIcon = computed(() => {
  const icons = {
    success: CheckIcon,
    error: ErrorIcon,
    warning: WarningIcon,
    info: InfoIcon
  }
  return icons[props.type]
})
</script>

<style scoped>
.notification-toast {
  position: fixed;
  top: 2rem;
  right: 2rem;
  min-width: 320px;
  max-width: 420px;
  background: white;
  border-radius: 12px;
  padding: 1rem 1rem 1rem 0.875rem;
  display: flex;
  align-items: center;
  gap: 0.875rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  overflow: hidden;
}

.notification-toast::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
}

.notification-toast--success::before {
  background: linear-gradient(180deg, #10b981 0%, #34d399 100%);
}

.notification-toast--error::before {
  background: linear-gradient(180deg, #ef4444 0%, #f87171 100%);
}

.notification-toast--warning::before {
  background: linear-gradient(180deg, #f59e0b 0%, #fbbf24 100%);
}

.notification-toast--info::before {
  background: linear-gradient(180deg, #3b82f6 0%, #60a5fa 100%);
}

.toast-icon-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notification-toast--success .toast-icon-wrapper {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(52, 211, 153, 0.1) 100%);
}

.notification-toast--success .toast-icon {
  color: #10b981;
}

.notification-toast--error .toast-icon-wrapper {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(248, 113, 113, 0.1) 100%);
}

.notification-toast--error .toast-icon {
  color: #ef4444;
}

.notification-toast--warning .toast-icon-wrapper {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.15) 0%, rgba(251, 191, 36, 0.1) 100%);
}

.notification-toast--warning .toast-icon {
  color: #f59e0b;
}

.notification-toast--info .toast-icon-wrapper {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15) 0%, rgba(96, 165, 250, 0.1) 100%);
}

.notification-toast--info .toast-icon {
  color: #3b82f6;
}

.toast-icon {
  width: 20px;
  height: 20px;
}

.toast-content {
  flex: 1;
  min-width: 0;
}

.toast-message {
  font-size: 0.9rem;
  font-weight: 500;
  color: #1e293b;
}

.toast-close {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: #94a3b8;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.toast-close:hover {
  background: #f1f5f9;
  color: #64748b;
}

.toast-close svg {
  width: 16px;
  height: 16px;
}

.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--progress-color) 0%, var(--progress-color-light) 100%);
  animation: progress 3s ease-out forwards;
}

.notification-toast--success {
  --progress-color: #10b981;
  --progress-color-light: #34d399;
}

.notification-toast--error {
  --progress-color: #ef4444;
  --progress-color-light: #f87171;
}

.notification-toast--warning {
  --progress-color: #f59e0b;
  --progress-color-light: #fbbf24;
}

.notification-toast--info {
  --progress-color: #3b82f6;
  --progress-color-light: #60a5fa;
}

@keyframes progress {
  from { width: 100%; }
  to { width: 0%; }
}

.toast-enter-active {
  animation: toast-in 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.toast-leave-active {
  animation: toast-out 0.2s ease-in;
}

@keyframes toast-in {
  from {
    opacity: 0;
    transform: translateX(100%) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateX(0) scale(1);
  }
}

@keyframes toast-out {
  from {
    opacity: 1;
    transform: translateX(0) scale(1);
  }
  to {
    opacity: 0;
    transform: translateX(20px) scale(0.95);
  }
}

@media (max-width: 480px) {
  .notification-toast {
    right: 1rem;
    left: 1rem;
    min-width: 0;
  }
}
</style>
