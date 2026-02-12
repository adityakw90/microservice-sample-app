<template>
  <div class="service-section" :class="{ 'not-implemented': notImplemented }">
    <div class="service-header" @click="collapsed = !collapsed">
      <div class="service-info">
        <component :is="icon" class="service-icon" />
        <h2 class="service-title">{{ title }}</h2>
        <span v-if="notImplemented" class="badge not-implemented-badge">Not Implemented</span>
      </div>
      <div class="service-actions">
        <span class="endpoint-count">{{ endpoints.length }} endpoints</span>
        <svg class="collapse-icon" :class="{ collapsed }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="6,9 12,15 18,9"/>
        </svg>
      </div>
    </div>

    <div v-show="!collapsed" class="endpoints-container">
      <div
        v-for="(endpoint, key) in endpoints"
        :key="key"
        class="endpoint-card"
        :class="{
          'has-response': endpoint.response,
          'has-error': endpoint.error,
          'loading': loading[key],
          'expanded': endpoint.expanded
        }"
      >
        <div class="endpoint-header" @click="endpoint.expanded = !endpoint.expanded">
          <div class="endpoint-info">
            <span class="method-badge" :class="endpoint.method.toLowerCase()">{{ endpoint.method }}</span>
            <span class="endpoint-name">{{ endpoint.name }}</span>
            <span class="endpoint-description">{{ endpoint.description }}</span>
          </div>
          <div class="endpoint-status">
            <span v-if="loading[key]" class="status-indicator loading">Loading...</span>
            <span v-else-if="endpoint.response" class="status-indicator success">Success</span>
            <span v-else-if="endpoint.error" class="status-indicator error">Error</span>
            <svg class="expand-icon" :class="{ 'expanded': endpoint.expanded }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6,9 12,15 18,9"/>
            </svg>
          </div>
        </div>

        <div v-show="endpoint.expanded" class="endpoint-body">
          <!-- Parameters Form -->
          <form @submit.prevent="handleSubmit(key, $event)" class="endpoint-form">
            <div class="form-fields">
              <div
                v-for="param in endpoint.params"
                :key="param.name"
                class="form-field"
              >
                <label class="field-label">
                  <span class="field-name">{{ param.name }}</span>
                  <span class="field-type">{{ param.type }}</span>
                  <span v-if="param.required" class="field-required">*</span>
                </label>
                <input
                  v-if="param.type === 'boolean'"
                  type="checkbox"
                  :name="param.name"
                  :checked="param.default === true"
                  class="field-input checkbox"
                />
                <input
                  v-else
                  :type="param.type === 'number' ? 'number' : 'text'"
                  :name="param.name"
                  :placeholder="param.default !== undefined ? String(param.default) : ''"
                  :required="param.required"
                  class="field-input"
                />
                <span v-if="param.description" class="field-description">{{ param.description }}</span>
              </div>
            </div>

            <div class="form-actions">
              <button
                type="submit"
                :disabled="loading[key] || notImplemented"
                class="execute-btn"
                :class="{ loading: loading[key] }"
              >
                <span v-if="loading[key]" class="btn-spinner"></span>
                <span v-else>{{ notImplemented ? 'Not Implemented' : 'Execute' }}</span>
              </button>
            </div>
          </form>

          <!-- Response Display -->
          <div v-if="endpoint.response || endpoint.error" class="response-section">
            <div class="response-header">
              <h4>Response</h4>
              <button @click="clearResponse(key)" class="clear-btn">Clear</button>
            </div>
            <pre v-if="endpoint.response" class="response-content success">{{ JSON.stringify(endpoint.response, null, 2) }}</pre>
            <pre v-else class="response-content error">{{ JSON.stringify(endpoint.error, null, 2) }}</pre>
          </div>
        </div>
      </div>

      <div v-if="Object.keys(endpoints).length === 0" class="no-endpoints">
        <p>No endpoints defined</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface ParamDefinition {
  name: string
  type: string
  required: boolean
  default?: any
  description?: string
}

interface Endpoint {
  name: string
  method: string
  description: string
  params: ParamDefinition[]
  response?: any
  error?: any
  expanded?: boolean
}

interface Endpoints {
  [key: string]: Endpoint
}

const props = defineProps<{
  title: string
  icon: any
  endpoints: Endpoints
  loading: Record<string, boolean>
  notImplemented?: boolean
}>()

const emit = defineEmits<{
  execute: [endpointName: string, formData: FormData]
}>()

const collapsed = ref(false)

const endpointList = computed(() => Object.values(props.endpoints))

const handleSubmit = (endpointName: string, event: Event) => {
  const form = event.target as HTMLFormElement
  const formData = new FormData(form)
  emit('execute', endpointName, formData)
}

const clearResponse = (key: string) => {
  const endpoint = props.endpoints[key]
  if (endpoint) {
    endpoint.response = null
    endpoint.error = null
  }
}
</script>

<style scoped>
.service-section {
  background: var(--color-surface);
  border-radius: 16px;
  border: 1px solid var(--color-border);
  overflow: hidden;
  transition: all 0.3s ease;
}

.service-section:hover {
  border-color: rgba(99, 102, 241, 0.3);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.service-section.not-implemented {
  opacity: 0.7;
}

.service-section.not-implemented .service-header {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(251, 146, 60, 0.05) 100%);
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s ease;
}

.service-header:hover {
  background: var(--color-bg);
}

.service-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.service-icon {
  width: 24px;
  height: 24px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.service-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.badge {
  padding: 0.2rem 0.5rem;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.not-implemented-badge {
  background: linear-gradient(135deg, #f59e0b 0%, #fb923c 100%);
  color: white;
}

.service-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.endpoint-count {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.collapse-icon {
  width: 20px;
  height: 20px;
  color: var(--color-text-secondary);
  transition: transform 0.3s ease;
}

.collapse-icon.collapsed {
  transform: rotate(-90deg);
}

.endpoints-container {
  padding: 0 1rem 1rem 1rem;
}

.endpoint-card {
  background: var(--color-bg);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  margin-bottom: 1rem;
  overflow: hidden;
  transition: all 0.3s ease;
}

.endpoint-card:last-child {
  margin-bottom: 0;
}

.endpoint-card.has-response {
  border-color: var(--color-success);
}

.endpoint-card.has-error {
  border-color: var(--color-error);
}

.endpoint-card.expanded {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.endpoint-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.875rem 1rem;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s ease;
}

.endpoint-header:hover {
  background: rgba(0, 0, 0, 0.02);
}

.endpoint-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
  min-width: 0;
}

.method-badge {
  padding: 0.2rem 0.5rem;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  flex-shrink: 0;
}

.method-badge.get {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.method-badge.post {
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: white;
}

.method-badge.put {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.method-badge.delete {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
}

.endpoint-name {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text-primary);
  flex-shrink: 0;
}

.endpoint-description {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.endpoint-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-indicator {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 0.15rem 0.5rem;
  border-radius: 6px;
}

.status-indicator.loading {
  background: var(--color-bg);
  color: var(--color-text-secondary);
}

.status-indicator.success {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.status-indicator.error {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
}

.expand-icon {
  width: 16px;
  height: 16px;
  color: var(--color-text-secondary);
  transition: transform 0.3s ease;
  flex-shrink: 0;
}

.expand-icon.expanded {
  transform: rotate(180deg);
}

.endpoint-body {
  border-top: 1px solid var(--color-border);
}

.endpoint-form {
  padding: 1rem;
}

.form-fields {
  display: grid;
  gap: 0.875rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.field-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8rem;
}

.field-name {
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.field-type {
  padding: 0.1rem 0.35rem;
  background: var(--color-bg);
  border-radius: 4px;
  font-size: 0.7rem;
  color: var(--color-text-secondary);
  font-family: 'JetBrains Mono', monospace;
}

.field-required {
  color: var(--color-error);
  font-weight: 700;
}

.field-input {
  padding: 0.6rem 0.875rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-surface);
  color: var(--color-text-primary);
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.field-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.field-input.checkbox {
  width: auto;
  cursor: pointer;
}

.field-description {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  font-style: italic;
}

.form-actions {
  margin-top: 1rem;
  display: flex;
  justify-content: flex-end;
}

.execute-btn {
  padding: 0.6rem 1.25rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.execute-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.execute-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid white;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.response-section {
  border-top: 1px solid var(--color-border);
  background: var(--color-bg);
}

.response-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
}

.response-header h4 {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.clear-btn {
  padding: 0.3rem 0.6rem;
  border-radius: 6px;
  font-size: 0.75rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  background: var(--color-error);
  color: white;
  border-color: var(--color-error);
}

.response-content {
  margin: 0;
  padding: 1rem;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.75rem;
  line-height: 1.6;
  max-height: 300px;
  overflow-y: auto;
  border-top: 1px solid var(--color-border);
}

.response-content.success {
  color: #10b981;
  background: rgba(16, 185, 129, 0.05);
}

.response-content.error {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.05);
}

.no-endpoints {
  padding: 2rem;
  text-align: center;
  color: var(--color-text-secondary);
}

@media (max-width: 768px) {
  .service-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .endpoint-info {
    flex-wrap: wrap;
  }

  .endpoint-description {
    display: none;
  }
}
</style>
