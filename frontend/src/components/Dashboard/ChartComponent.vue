<template>
  <div class="chart-component">
    <div class="chart-header">
      <div>
        <h3 class="chart-title">{{ title }}</h3>
        <p class="chart-subtitle">{{ subtitle }}</p>
      </div>
      <div class="chart-legend">
        <div class="legend-item">
          <div class="legend-dot legend-dot--primary"></div>
          <span>Registrations</span>
        </div>
      </div>
    </div>

    <div class="chart-container">
      <div class="chart-bars">
        <div
          v-for="(item, index) in data"
          :key="index"
          class="chart-bar-wrapper"
        >
          <div class="chart-bar-container">
            <transition name="bar-grow">
              <div
                class="chart-bar"
                :style="{
                  height: getBarHeight(item.value) + '%',
                  background: getBarColor(index)
                }"
              >
                <div class="chart-bar-glow"></div>
                <div class="chart-bar-tooltip">
                  <div class="tooltip-label">{{ item.label }}</div>
                  <div class="tooltip-value">{{ item.value }} users</div>
                </div>
              </div>
            </transition>
          </div>
          <div class="chart-label">{{ item.label }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineProps<{
  title: string
  subtitle: string
  data: { label: string; value: number }[]
}>()

const maxValue = computed(() => {
  return Math.max(...[...arguments[0]?.data?.map((d: any) => d.value) || []])
})

const getBarHeight = (value: number) => {
  const max = maxValue.value || 1
  return Math.max((value / max) * 100, 5)
}

const getBarColor = (index: number) => {
  const colors = [
    'linear-gradient(180deg, #6366f1 0%, #8b5cf6 100%)',
    'linear-gradient(180deg, #8b5cf6 0%, #a78bfa 100%)',
    'linear-gradient(180deg, #6366f1 0%, #818cf8 100%)'
  ]
  return colors[index % colors.length]
}
</script>

<style scoped>
.chart-component {
  background: white;
  border-radius: 20px;
  padding: 1.5rem;
  border: 1px solid rgba(0, 0, 0, 0.04);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.chart-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 0.25rem 0;
}

.chart-subtitle {
  font-size: 0.85rem;
  color: #64748b;
  margin: 0;
}

.chart-legend {
  display: flex;
  gap: 1rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8rem;
  color: #64748b;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.legend-dot--primary {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
}

.chart-container {
  height: 200px;
  position: relative;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  height: 100%;
  gap: 1rem;
  padding: 0 0.5rem;
}

.chart-bar-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.chart-bar-container {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.chart-bar {
  position: relative;
  width: 100%;
  max-width: 48px;
  border-radius: 8px 8px 4px 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.chart-bar:hover {
  transform: scaleY(1.05);
  filter: brightness(1.1);
}

.chart-bar-glow {
  position: absolute;
  inset: 0;
  background: inherit;
  filter: blur(12px);
  opacity: 0.4;
  border-radius: inherit;
}

.chart-bar-tooltip {
  position: absolute;
  bottom: calc(100% + 0.75rem);
  left: 50%;
  transform: translateX(-50%);
  background: #1e293b;
  color: white;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  font-size: 0.75rem;
  white-space: nowrap;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s ease;
}

.chart-bar-tooltip::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 5px solid transparent;
  border-top-color: #1e293b;
}

.chart-bar:hover .chart-bar-tooltip {
  opacity: 1;
}

.tooltip-label {
  font-weight: 500;
  margin-bottom: 0.15rem;
}

.tooltip-value {
  color: #94a3b8;
}

.chart-label {
  margin-top: 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
  color: #64748b;
  text-align: center;
}

.bar-grow-enter-active {
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.bar-grow-enter-from {
  height: 0 !important;
  opacity: 0;
}
</style>
