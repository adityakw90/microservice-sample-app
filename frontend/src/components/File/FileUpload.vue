<template>
  <div class="file-upload-modal" @click.self="$emit('close')">
    <div class="upload-container" @click.self="$emit('close')">
      <div class="upload-content" @click.self="$emit('close')">
        <div class="upload-header">
          <h2>Upload File</h2>
          <button @click="$emit('close')" class="close-btn">&times;</button>
        </div>

        <div
          class="upload-zone"
          :class="{ 'drag-over': isDragOver }"
          @dragover.prevent="isDragOver = true"
          @dragleave.prevent="isDragOver = false"
          @drop.prevent="handleDrop"
          @click="$refs.fileInput.click()"
        >
          <input
            ref="fileInput"
            type="file"
            @change="handleFileSelect"
            class="file-input"
          />
          <svg class="upload-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9 6L16 6a5 5 0 011 2 2v8a2 2 0 01-2 2 2 2 2-4 4 1 1 4 4-9.5 9.5z"/>
          </svg>
          <p class="upload-text">
            {{ isDragOver ? 'Drop your file here' : 'Drag & drop a file or click to browse' }}
          </p>
          <p class="upload-hint">Maximum file size: 10MB</p>
        </div>

        <div v-if="selectedFile" class="file-preview">
          <div class="file-info">
            <span class="file-name">{{ selectedFile.name }}</span>
            <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
          </div>
          <button @click="selectedFile = null" class="remove-btn">&times;</button>
        </div>

        <div v-if="uploadProgress > 0" class="upload-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
          </div>
          <p class="progress-text">{{ uploadProgress }}% uploaded</p>
        </div>

        <div class="upload-options">
          <label class="visibility-toggle">
            <input type="checkbox" v-model="makePublic" />
            <span>Make file public</span>
          </label>
        </div>

        <div class="upload-actions">
          <button @click="$emit('close')" class="cancel-btn">Cancel</button>
          <button
            @click="handleUpload"
            :disabled="!selectedFile || uploading"
            class="upload-btn"
          >
            {{ uploading ? 'Uploading...' : 'Upload' }}
          </button>
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useFilesStore } from '@/stores/files'

const emit = defineEmits<{
  close: []
  uploaded: []
}>()

const authStore = useAuthStore()
const filesStore = useFilesStore()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const isDragOver = ref(false)
const makePublic = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const error = ref<string | null>(null)

function handleDrop(e: DragEvent) {
  isDragOver.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    validateAndSelectFile(files[0])
  }
}

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    validateAndSelectFile(target.files[0])
  }
}

function validateAndSelectFile(file: File) {
  error.value = null

  // Check file size (10MB max)
  if (file.size > 10 * 1024 * 1024) {
    error.value = 'File size exceeds 10MB limit'
    return
  }

  selectedFile.value = file
}

async function handleUpload() {
  if (!selectedFile.value || !authStore.user?.uid) return

  uploading.value = true
  error.value = null
  uploadProgress.value = 0

  try {
    // Simulate upload progress
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += 10
      }
    }, 200)

    // Read file as base64
    const reader = new FileReader()
    reader.onload = async (e) => {
      const base64 = (e.target?.result as string).split(',')[1]

      const uid = await filesStore.uploadFile({
        user_uid: authStore.user!.uid,
        file_type: getFileType(selectedFile.value!.type),
        file_name: selectedFile.value!.name,
        file_path: '',  // Will be set by server
        mime_type: selectedFile.value!.type,
        size: selectedFile.value!.size,
        visibility: makePublic.value ? 1 : 0
      })

      clearInterval(progressInterval)
      uploadProgress.value = 100

      if (uid) {
        setTimeout(() => {
          emit('uploaded')
          resetForm()
        }, 500)
      }
    }
    reader.readAsDataURL(selectedFile.value)
  } catch (err: any) {
    error.value = err.message || 'Upload failed'
    uploading.value = false
  }
}

function getFileType(mimeType: string): string {
  if (mimeType.startsWith('image/')) return 'image'
  if (mimeType.startsWith('video/')) return 'video'
  if (mimeType.includes('pdf')) return 'pdf'
  return 'document'
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function resetForm() {
  selectedFile.value = null
  uploading.value = false
  uploadProgress.value = 0
  error.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}
</script>

<style scoped>
.file-upload-modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.upload-container {
  width: 100%;
  max-width: 500px;
  padding: 2rem;
}

.upload-content {
  background: white;
  border-radius: 12px;
  padding: 2rem;
}

.upload-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.upload-header h2 {
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #7f8c8d;
}

.upload-zone {
  border: 2px dashed #cbd5e0;
  border-radius: 12px;
  padding: 3rem 2rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-zone:hover,
.upload-zone.drag-over {
  border-color: #667eea;
  background: #f8f9ff;
}

.file-input {
  display: none;
}

.upload-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  color: #cbd5e0;
}

.upload-text {
  margin: 0 0 0.5rem 0;
  font-weight: 500;
  color: #2c3e50;
}

.upload-hint {
  margin: 0;
  font-size: 0.875rem;
  color: #7f8c8d;
}

.file-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  margin: 1rem 0;
}

.file-info {
  display: flex;
  flex-direction: column;
}

.file-name {
  font-weight: 500;
  color: #2c3e50;
}

.file-size {
  font-size: 0.875rem;
  color: #7f8c8d;
}

.remove-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #dc2626;
}

.upload-progress {
  margin: 1rem 0;
}

.progress-bar {
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s;
}

.progress-text {
  margin: 0.5rem 0 0;
  font-size: 0.875rem;
  color: #7f8c8d;
  text-align: center;
}

.upload-options {
  margin: 1rem 0;
}

.visibility-toggle {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.upload-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.cancel-btn,
.upload-btn {
  flex: 1;
  padding: 0.875rem;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

.cancel-btn {
  background: #e2e8f0;
  color: #475569;
}

.upload-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.upload-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #fee;
  color: #c33;
  border-radius: 8px;
  text-align: center;
}
</style>
