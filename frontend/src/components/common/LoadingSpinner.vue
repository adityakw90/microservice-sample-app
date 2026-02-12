<template>
  <div class="loading-spinner" :class="`loading-spinner--${size}`">
    <div class="spinner-ring">
      <div class="spinner-segment spinner-segment--1"></div>
      <div class="spinner-segment spinner-segment--2"></div>
      <div class="spinner-segment spinner-segment--3"></div>
    </div>
    <div v-if="text" class="spinner-text">{{ text }}</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  size?: 'sm' | 'md' | 'lg'
  text?: string
}>()
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.spinner-ring {
  position: relative;
  width: 48px;
  height: 48px;
}

.loading-spinner--sm .spinner-ring {
  width: 32px;
  height: 32px;
}

.loading-spinner--lg .spinner-ring {
  width: 64px;
  height: 64px;
}

.spinner-segment {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  animation: spin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
}

.spinner-segment::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 3px solid transparent;
  border-top-color: #6366f1;
}

.spinner-segment--1::before {
  border-top-color: #6366f1;
  border-width: 3px;
}

.spinner-segment--2::before {
  border-top-color: #8b5cf6;
  border-width: 3px;
  filter: blur(1px);
}

.spinner-segment--3::before {
  border-top-color: #a78bfa;
  border-width: 2px;
  filter: blur(2px);
}

.spinner-segment--1 {
  animation-delay: 0s;
}

.spinner-segment--2 {
  animation-delay: -0.4s;
}

.spinner-segment--3 {
  animation-delay: -0.8s;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.spinner-text {
  font-size: 0.875rem;
  font-weight: 500;
  color: #64748b;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
