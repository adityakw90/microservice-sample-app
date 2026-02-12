<template>
  <div class="user-files-tab">
    <div class="files-header">
      <h3>Files</h3>
      <button @click="showUploadModal = true" class="upload-btn-small">
        Upload New File
      </button>
    </div>

    <div v-if="filesStore.loading && filesStore.files.length === 0" class="loading-state">
      <LoadingSpinner />
    </div>

    <div v-else-if="filesStore.error" class="error-state">
      <p>{{ filesStore.error }}</p>
      <button @click="filesStore.clearError(); loadFiles()" class="retry-btn">Retry</button>
    </div>

    <div v-else-if="filesStore.isEmpty" class="empty-state">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2 2z"/>
      </svg>
      <p>No files yet</p>
      <button @click="showUploadModal = true" class="upload-link">Upload your first file</button>
    </div>

    <div v-else class="files-content">
      <div class="files-controls">
        <select v-model="viewMode" class="view-toggle">
          <option value="grid">Grid</option>
          <option value="list">List</option>
        </select>
        <select v-model="filterType" class="filter-select" @change="loadFiles()">
          <option value="">All Types</option>
          <option value="image">Images</option>
          <option value="document">Documents</option>
          <option value="video">Videos</option>
        </select>
      </div>

      <div v-if="viewMode === 'grid'" class="files-grid">
        <FileCard
          v-for="file in filesStore.files"
          :key="file.uid"
          :file="file"
          @delete="confirmDelete(file)"
          @update="openUpdateModal(file)"
        />
      </div>

      <div v-else class="files-list">
        <div v-for="file in filesStore.files" :key="file.uid" class="file-row">
          <FileCard
            :file="file"
            :list-view="true"
            @delete="confirmDelete(file)"
            @update="openUpdateModal(file)"
          />
        </div>
      </div>

      <div v-if="filesStore.hasMore" class="load-more">
        <button @click="loadMore()" :disabled="filesStore.loading" class="load-more-btn">
          {{ filesStore.loading ? 'Loading...' : 'Load More' }}
        </button>
      </div>
    </div>

    <!-- Upload Modal -->
    <FileUpload
      v-if="showUploadModal"
      :user-uid="userUid"
      @close="showUploadModal = false"
      @uploaded="handleFileUploaded"
    />

    <!-- Update Modal -->
    <div v-if="showUpdateModal" class="modal-overlay" @click.self="showUpdateModal = false">
      <div class="modal-content">
        <h2>Update File</h2>
        <form @submit.prevent="handleUpdateFile">
          <div class="form-group">
            <label for="fileName">File Name</label>
            <input
              id="fileName"
              v-model="updateForm.file_name"
              type="text"
              class="form-input"
              required
            />
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="updateForm.public" />
              Make Public
            </label>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showUpdateModal = false" class="cancel-btn">Cancel</button>
            <button type="submit" class="submit-btn" :disabled="filesStore.loading">
              {{ filesStore.loading ? 'Updating...' : 'Update' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Confirm Delete Dialog -->
    <ConfirmDialog
      v-if="showDeleteConfirm"
      title="Delete File"
      :message="`Are you sure you want to delete '${fileToDelete?.file_name}'?`"
      @confirm="handleDeleteFile"
      @cancel="showDeleteConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useFilesStore } from '@/stores/files'
import FileCard from '@/components/File/FileCard.vue'
import FileUpload from '@/components/File/FileUpload.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import type { UserFile } from '@/services/api'

const props = defineProps<{
  userUid: string
}>()

const filesStore = useFilesStore()

const viewMode = ref<'grid' | 'list'>('grid')
const filterType = ref('')
const showUploadModal = ref(false)
const showUpdateModal = ref(false)
const showDeleteConfirm = ref(false)
const fileToDelete = ref<UserFile | null>(null)

const updateForm = ref({
  file_name: '',
  public: false
})

onMounted(() => {
  loadFiles()
})

async function loadFiles() {
  await filesStore.fetchFiles({
    user_uid: props.userUid,
    file_type: filterType.value || undefined
  })
}

async function loadMore() {
  const nextPage = filesStore.pagination.page + 1
  await filesStore.fetchFiles({
    page: nextPage,
    user_uid: props.userUid
  })
}

function handleFileUploaded() {
  showUploadModal.value = false
  loadFiles()
}

function openUpdateModal(file: UserFile) {
  updateForm.value = {
    file_name: file.file_name,
    public: file.visibility === 1
  }
  fileToDelete.value = file
  showUpdateModal.value = true
}

async function handleUpdateFile() {
  if (!fileToDelete.value) return

  const success = await filesStore.updateFile(fileToDelete.value.uid, {
    file_name: updateForm.value.file_name,
    visibility: updateForm.value.public ? 1 : 0
  })

  if (success) {
    showUpdateModal.value = false
  }
}

function confirmDelete(file: UserFile) {
  fileToDelete.value = file
  showDeleteConfirm.value = true
}

async function handleDeleteFile() {
  if (!fileToDelete.value) return

  const success = await filesStore.deleteFile(fileToDelete.value.uid)
  if (success) {
    showDeleteConfirm.value = false
    fileToDelete.value = null
  }
}
</script>

<style scoped>
.user-files-tab {
  padding: 1rem 0;
}

.files-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.files-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #2c3e50;
}

.upload-btn-small {
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.upload-btn-small:hover {
  transform: translateY(-1px);
}

.files-controls {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.view-toggle,
.filter-select {
  padding: 0.4rem 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: white;
  font-size: 0.875rem;
}

.files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.files-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.load-more {
  text-align: center;
  margin-top: 1.5rem;
}

.load-more-btn {
  padding: 0.5rem 1.5rem;
  background: white;
  border: 1px solid #ddd;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
}

.load-more-btn:hover:not(:disabled) {
  background: #f8f9fa;
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
}

.empty-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  color: #cbd5e0;
}

.upload-link {
  color: #667eea;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.875rem;
}

.retry-btn {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 1.5rem;
  border-radius: 10px;
  width: 100%;
  max-width: 400px;
}

.modal-content h2 {
  margin: 0 0 1rem 0;
  font-size: 1.25rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  font-size: 0.875rem;
}

.form-input {
  width: 100%;
  padding: 0.6rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 1rem;
}

.cancel-btn,
.submit-btn {
  padding: 0.6rem 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
}

.cancel-btn {
  background: #e2e8f0;
  color: #475569;
}

.submit-btn {
  background: #667eea;
  color: white;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
