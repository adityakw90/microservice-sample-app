<template>
  <div class="stat-card" :class="[`stat-card--${color}`, { 'stat-card--loading': loading }]">
    <div class="stat-card__background">
      <div class="stat-card__gradient-1"></div>
      <div class="stat-card__gradient-2"></div>
      <div class="stat-card__noise"></div>
    </div>

    <div class="stat-card__content">
      <div class="stat-card__header">
        <div class="stat-card__icon-wrapper">
          <component :is="icon" class="stat-card__icon" />
          <div class="stat-card__icon-glow"></div>
        </div>
        <div v-if="trend" class="stat-card__trend" :class="trendClass">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline :points="trendUp ? '23 6 13.5 15.5 8.5 10.5 1 18' : '23 18 13.5 8.5 8.5 13.5 1 6'"/>
            <polyline v-if="trendUp" points="17 6 23 6 23 12"/>
            <polyline v-else points="17 18 23 18 23 12"/>
          </svg>
          <span>{{ trend }}</span>
        </div>
      </div>

      <div class="stat-card__body">
        <div class="stat-card__value">
          <transition name="count" mode="out-in">
            <span :key="displayValue" class="stat-card__number">{{ displayValue }}</span>
          </transition>
          <span v-if="suffix" class="stat-card__suffix">{{ suffix }}</span>
        </div>
        <div class="stat-card__label">{{ label }}</div>
      </div>

      <div class="stat-card__footer">
        <div class="stat-card__progress">
          <div class="stat-card__progress-bar" :style="{ width: progress + '%' }"></div>
        </div>
        <div class="stat-card__helper">{{ helper }}</div>
      </div>
    </div>

    <div v-if="loading" class="stat-card__skeleton">
      <div class="skeleton-icon"></div>
      <div class="skeleton-value"></div>
      <div class="skeleton-label"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, onMounted } from 'vue'

const props = defineProps<{
  label: string
  value: number
  suffix?: string
  icon: any
  color?: 'indigo' | 'emerald' | 'amber' | 'rose' | 'violet'
  trend?: string
  trendUp?: boolean
  progress?: number
  helper?: string
  loading?: boolean
}>()

const displayValue = ref('0')
const targetValue = computed(() => props.value)

const animateValue = (start: number, end: number, duration: number = 1000) => {
  const startTime = performance.now()
  const range = end - start

  const update = (currentTime: number) => {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)

    // Easing function
    const easeOut = 1 - Math.pow(1 - progress, 3)
    const current = start + (range * easeOut)

    displayValue.value = Math.round(current).toLocaleString()

    if (progress < 1) {
      requestAnimationFrame(update)
    }
  }

  requestAnimationFrame(update)
}

watch(targetValue, (newVal, oldVal) => {
  const start = oldVal || 0
  animateValue(start, newVal)
}, { immediate: true })

onMounted(() => {
  animateValue(0, props.value)
})

const trendClass = computed(() => ({
  'stat-card__trend--up': props.trendUp,
  'stat-card__trend--down': !props.trendUp
}))
</script>

<style scoped>
.stat-card {
  position: relative;
  background: white;
  border-radius: 20px;
  padding: 1.5rem;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.stat-card__background {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.stat-card__gradient-1 {
  position: absolute;
  width: 200%;
  height: 200%;
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 30% 30%, var(--card-color-1) 0%, transparent 50%);
  opacity: 0.1;
  animation: gradientShift 8s ease-in-out infinite;
}

.stat-card__gradient-2 {
  position: absolute;
  width: 150%;
  height: 150%;
  bottom: -25%;
  right: -25%;
  background: radial-gradient(circle at 70% 70%, var(--card-color-2) 0%, transparent 50%);
  opacity: 0.08;
  animation: gradientShift 10s ease-in-out infinite reverse;
}

.stat-card__noise {
  position: absolute;
  inset: 0;
  opacity: 0.02;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}

.stat-card--indigo {
  --card-color-1: #6366f1;
  --card-color-2: #8b5cf6;
}

.stat-card--emerald {
  --card-color-1: #10b981;
  --card-color-2: #34d399;
}

.stat-card--amber {
  --card-color-1: #f59e0b;
  --card-color-2: #fbbf24;
}

.stat-card--rose {
  --card-color-1: #ef4444;
  --card-color-2: #f87171;
}

.stat-card--violet {
  --card-color-1: #8b5cf6;
  --card-color-2: #a78bfa;
}

@keyframes gradientShift {
  0%, 100% { transform: translate(0, 0) rotate(0deg); }
  50% { transform: translate(5%, 5%) rotate(5deg); }
}

.stat-card__content {
  position: relative;
  z-index: 1;
}

.stat-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.25rem;
}

.stat-card__icon-wrapper {
  position: relative;
  width: 52px;
  height: 52px;
}

.stat-card__icon {
  position: relative;
  z-index: 2;
  width: 100%;
  height: 100%;
  color: var(--card-color-1);
}

.stat-card__icon-glow {
  position: absolute;
  inset: 0;
  background: var(--card-color-1);
  filter: blur(16px);
  opacity: 0.3;
  animation: glow 3s ease-in-out infinite;
}

@keyframes glow {
  0%, 100% { opacity: 0.2; transform: scale(0.95); }
  50% { opacity: 0.4; transform: scale(1.05); }
}

.stat-card__trend {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.625rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
}

.stat-card__trend svg {
  width: 14px;
  height: 14px;
}

.stat-card__trend--up {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(52, 211, 153, 0.1) 100%);
  color: #059669;
}

.stat-card__trend--down {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(248, 113, 113, 0.1) 100%);
  color: #dc2626;
}

.stat-card__body {
  margin-bottom: 1.25rem;
}

.stat-card__value {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
  margin-bottom: 0.5rem;
}

.stat-card__number {
  font-size: 2.25rem;
  font-weight: 700;
  color: #1e293b;
  line-height: 1;
  font-family: var(--font-mono);
}

.count-enter-active {
  transition: all 0.3s ease;
}

.count-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.stat-card__suffix {
  font-size: 1rem;
  font-weight: 500;
  color: #64748b;
}

.stat-card__label {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
}

.stat-card__footer {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.stat-card__progress {
  flex: 1;
  height: 4px;
  background: #f1f5f9;
  border-radius: 2px;
  overflow: hidden;
}

.stat-card__progress-bar {
  height: 100%;
  background: linear-gradient(90deg, var(--card-color-1) 0%, var(--card-color-2) 100%);
  border-radius: 2px;
  transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card__helper {
  font-size: 0.7rem;
  color: #94a3b8;
  white-space: nowrap;
}

.stat-card__skeleton {
  position: absolute;
  inset: 0;
  background: white;
  z-index: 10;
  display: flex;
  flex-direction: column;
  padding: 1.5rem;
  gap: 1rem;
}

.skeleton-icon {
  width: 52px;
  height: 52px;
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  border-radius: 12px;
  animation: skeleton-loading 1.5s infinite;
}

.skeleton-value {
  width: 60%;
  height: 36px;
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  border-radius: 8px;
  animation: skeleton-loading 1.5s infinite;
}

.skeleton-label {
  width: 40%;
  height: 16px;
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  border-radius: 8px;
  animation: skeleton-loading 1.5s infinite;
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
