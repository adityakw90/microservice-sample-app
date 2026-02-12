<template>
  <div :class="['file-card', { 'list-view': listView }]">
    <div class="file-icon" :class="fileTypeClass">
      <svg v-if="isImage" viewBox="0 0 24 24" fill="currentColor">
        <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
      </svg>
      <svg v-else viewBox="0 0 24 24" fill="currentColor">
        <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.89 2-2V8l-6-6zm2 16H8v-2h4v2h-3v3h3V7h2v5h2V7h2v2h-2V5h-2z"/>
      </svg>
    </div>

    <div class="file-info">
      <h3 class="file-name">{{ file.file_name }}</h3>
      <p class="file-meta">{{ formatFileSize(file.size) }} · {{ formatDate(file.created_at) }}</p>
      <span v-if="file.visibility === 1" class="visibility-badge public">Public</span>
      <span v-else class="visibility-badge private">Private</span>
    </div>

    <div class="file-actions">
      <button @click="$emit('update', file)" class="action-btn" title="Rename">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h7m-9 1V7a2 2 0 0 12 2v4a2 2 0 0 1 2 2 2 2 2 2-4 4 1 1 4 4-9.5 9.5z"/>
        </svg>
      </button>
      <button @click="handleDownload" class="action-btn" title="Download">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2 2v-4"/>
          <polyline points="7,10 12,15 17,10"/>
          <line x1="12" y1="15" x2="12" y2="3"/>
        </svg>
      </button>
      <button @click="$emit('delete', file)" class="action-btn danger" title="Delete">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="3,6 5,6 21,6"/>
          <path d="M19,6v14a2,2,0,0,1-2,2H7a2,2,0,0,1-2-2V6m3,0V4a2,2,0,0,1,2-2,2H7a2,2,0,0,1,2,2v14a2,2,0,0,1-2,2H7a2,2,0,0,1,2,2V6"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserFile } from '@/services/api'

const props = defineProps<{
  file: UserFile
  listView?: boolean
}>()

const emit = defineEmits<{
  delete: [file: UserFile]
  update: [file: UserFile]
}>()

const isImage = computed(() => props.file.mime_type.startsWith('image/'))

const fileTypeClass = computed(() => {
  if (props.file.mime_type.startsWith('image/')) return 'image'
  if (props.file.mime_type.startsWith('video/')) return 'video'
  if (props.file.mime_type.includes('pdf')) return 'pdf'
  return 'document'
})

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString()
}

function handleDownload() {
  // TODO: Implement file download
  console.log('Download file:', props.file.uid)
}
</script>

<style scoped>
.file-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.file-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
}

.file-card.list-view {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
}

.file-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
  color: white;
}

.file-card.list-view .file-icon {
  margin-bottom: 0;
  flex-shrink: 0;
}

.file-icon.image { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.file-icon.video { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.file-icon.pdf { background: linear-gradient(135deg, #fa709a 0%, #fee140 100%); }
.file-icon.document { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }

.file-icon svg {
  width: 24px;
  height: 24px;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-meta {
  margin: 0;
  font-size: 0.875rem;
  color: #7f8c8d;
}

.visibility-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  margin-top: 0.5rem;
}

.visibility-badge.public {
  background: #d1fae5;
  color: #065f46;
}

.visibility-badge.private {
  background: #e2e8f0;
  color: #475569;
}

.file-actions {
  display: flex;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity 0.2s;
}

.file-card:hover .file-actions {
  opacity: 1;
}

.file-card.list-view .file-actions {
  opacity: 1;
}

.action-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f3f4f6;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #4b5563;
  transition: background 0.2s;
}

.action-btn:hover {
  background: #e5e7eb;
}

.action-btn.danger:hover {
  background: #fee2e2;
  color: #dc2626;
}

.action-btn svg {
  width: 18px;
  height: 18px;
}
</style>
