<template>
  <div class="song-management">
    <div class="management-header">
      <div class="header-left">
        <h2 class="page-title">歌曲管理</h2>
        <p class="page-subtitle">管理英语学习歌曲内容</p>
      </div>
      <div class="header-right">
        <el-button 
          type="primary" 
          :icon="Plus" 
          @click="handleAdd"
          size="large"
        >
          新增歌曲
        </el-button>
        <el-button 
          :icon="Upload" 
          @click="handleBatchImport"
          size="large"
        >
          批量导入
        </el-button>
      </div>
    </div>

    <!-- 搜索过滤器 -->
    <div class="search-filters">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索歌曲标题或描述"
            :prefix-icon="Search"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-col>
        <el-col :span="4">
          <el-select v-model="searchForm.category_id" placeholder="选择分类" clearable>
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name_cn"
              :value="category.id"
            />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="searchForm.difficulty" placeholder="难度等级" clearable>
            <el-option label="入门 (1星)" :value="1" />
            <el-option label="初级 (2星)" :value="2" />
            <el-option label="中级 (3星)" :value="3" />
            <el-option label="高级 (4星)" :value="4" />
            <el-option label="专家 (5星)" :value="5" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="searchForm.is_published" placeholder="发布状态" clearable>
            <el-option label="已发布" :value="true" />
            <el-option label="草稿" :value="false" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="RefreshRight" @click="handleReset">重置</el-button>
          <el-button 
            type="danger" 
            :icon="Delete" 
            :disabled="selectedSongs.length === 0"
            @click="handleBatchDelete"
          >
            批量删除
          </el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 歌曲列表表格 -->
    <div class="song-table">
      <el-table
        v-loading="loading"
        :data="songs"
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column type="index" label="#" width="60" />
        
        <el-table-column label="歌曲信息" min-width="300">
          <template #default="scope">
            <div class="song-info">
              <div class="song-cover">
                <img 
                  :src="scope.row.cover_image || defaultCover" 
                  :alt="scope.row.title"
                  class="cover-img"
                />
              </div>
              <div class="song-details">
                <div class="song-title">{{ scope.row.title }}</div>
                <div class="song-title-cn">{{ scope.row.title_cn }}</div>
                <div class="song-meta">
                  <el-tag 
                    v-if="scope.row.category" 
                    :color="scope.row.category.color"
                    size="small"
                    effect="light"
                  >
                    {{ scope.row.category.name_cn }}
                  </el-tag>
                  <span class="duration">{{ formatDuration(scope.row.duration) }}</span>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="难度" width="120" sortable="custom" prop="difficulty">
          <template #default="scope">
            <div class="difficulty-stars">
              <el-icon 
                v-for="i in 5" 
                :key="i" 
                :class="['star', { active: i <= scope.row.difficulty }]"
              >
                <Star />
              </el-icon>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="统计数据" width="150">
          <template #default="scope">
            <div class="stats">
              <div class="stat-item">
                <el-icon><View /></el-icon>
                <span>{{ formatNumber(scope.row.view_count) }}</span>
              </div>
              <div class="stat-item">
                <el-icon><Star /></el-icon>
                <span>{{ formatNumber(scope.row.like_count) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_published ? 'success' : 'info'">
              {{ scope.row.is_published ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="更新时间" width="160" sortable="custom" prop="updated_at">
          <template #default="scope">
            {{ formatDate(scope.row.updated_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="240" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button 
                size="small" 
                :icon="Edit" 
                @click="handleEdit(scope.row)"
              >
                编辑
              </el-button>
              <el-button 
                size="small" 
                :icon="View" 
                @click="handlePreview(scope.row)"
              >
                预览
              </el-button>
              <el-button 
                size="small" 
                type="danger" 
                :icon="Delete" 
                @click="handleDelete(scope.row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 歌曲编辑对话框 -->
    <song-edit-dialog
      v-model="editDialogVisible"
      :song="currentSong"
      :categories="categories"
      @success="handleEditSuccess"
    />

    <!-- 歌曲预览对话框 -->
    <song-preview-dialog
      v-model="previewDialogVisible"
      :song="currentSong"
      @edit="handlePreviewEdit"
    />

    <!-- 批量导入对话框 -->
    <batch-import-dialog
      v-model="batchImportVisible"
      :categories="categories"
      @success="handleImportSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Upload, Search, RefreshRight, Delete, Edit, View,
  Star
} from '@element-plus/icons-vue'
import { useEnglishLearningStore } from '../../stores/englishLearningStore'
import SongEditDialog from './SongEditDialog.vue'
import SongPreviewDialog from './SongPreviewDialog.vue'
import BatchImportDialog from './BatchImportDialog.vue'

const store = useEnglishLearningStore()

// 响应式数据
const loading = ref(false)
const editDialogVisible = ref(false)
const previewDialogVisible = ref(false)
const batchImportVisible = ref(false)
const currentSong = ref(null)
const selectedSongs = ref([])

// 搜索表单
const searchForm = ref({
  keyword: '',
  category_id: null,
  difficulty: null,
  is_published: null
})

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)

// 计算属性
const songs = computed(() => store.songs)
const categories = computed(() => store.categories)
const total = computed(() => store.songsTotal)

const defaultCover = 'https://via.placeholder.com/150x150/667eea/ffffff?text=♪'

// 生命周期
onMounted(async () => {
  await Promise.all([
    fetchSongs(),
    fetchCategories()
  ])
})

// 监听分页变化
watch([currentPage, pageSize], () => {
  fetchSongs()
})

// 方法
const fetchSongs = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      ...searchForm.value
    }
    // 移除空值
    Object.keys(params).forEach(key => {
      if (params[key] === null || params[key] === '') {
        delete params[key]
      }
    })
    
    await store.fetchSongs(params)
  } catch (error) {
    ElMessage.error('获取歌曲列表失败')
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    await store.fetchCategories({ limit: 100 })
  } catch (error) {
    console.warn('获取分类失败:', error)
  }
}

// 搜索和过滤
const handleSearch = () => {
  currentPage.value = 1
  fetchSongs()
}

const handleReset = () => {
  searchForm.value = {
    keyword: '',
    category_id: null,
    difficulty: null,
    is_published: null
  }
  currentPage.value = 1
  fetchSongs()
}

const handleSortChange = ({ prop, order }) => {
  const sortOrder = order === 'ascending' ? 'asc' : 'desc'
  searchForm.value.sort = `${prop}:${sortOrder}`
  fetchSongs()
}

// CRUD操作
const handleAdd = () => {
  currentSong.value = null
  editDialogVisible.value = true
}

const handleEdit = (song) => {
  currentSong.value = { ...song }
  editDialogVisible.value = true
}

const handleDelete = async (song) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除歌曲"${song.title}"吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        buttonSize: 'default'
      }
    )
    
    await store.deleteSong(song.id)
    ElMessage.success('删除成功')
    fetchSongs()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handlePreview = (song) => {
  currentSong.value = song
  previewDialogVisible.value = true
}

// 从预览对话框触发编辑
const handlePreviewEdit = (song) => {
  previewDialogVisible.value = false
  currentSong.value = { ...song }
  editDialogVisible.value = true
}

// 批量操作
const handleSelectionChange = (selection) => {
  selectedSongs.value = selection
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedSongs.value.length} 首歌曲吗？此操作不可恢复。`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const deletePromises = selectedSongs.value.map(song => store.deleteSong(song.id))
    await Promise.all(deletePromises)
    
    ElMessage.success(`成功删除 ${selectedSongs.value.length} 首歌曲`)
    selectedSongs.value = []
    fetchSongs()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const handleBatchImport = () => {
  batchImportVisible.value = true
}

// 分页
const handlePageChange = (page) => {
  currentPage.value = page
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

// 成功回调
const handleEditSuccess = () => {
  editDialogVisible.value = false
  fetchSongs()
}

const handleImportSuccess = () => {
  batchImportVisible.value = false
  fetchSongs()
}

// 工具方法
const formatDuration = (seconds) => {
  if (!seconds) return '--'
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const formatDate = (dateString) => {
  if (!dateString) return '--'
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.song-management {
  padding: 24px;
  background: #f8fafc;
  min-height: calc(100vh - 60px);
}

.management-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
}

.search-filters {
  background: white;
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.song-table {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.song-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.song-cover .cover-img {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  object-fit: cover;
}

.song-details {
  flex: 1;
  min-width: 0;
}

.song-title {
  font-weight: 600;
  color: #1f2937;
  font-size: 14px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-title-cn {
  color: #6b7280;
  font-size: 12px;
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.duration {
  font-size: 12px;
  color: #9ca3af;
}

.difficulty-stars {
  display: flex;
  gap: 2px;
}

.star {
  color: #d1d5db;
  font-size: 14px;
}

.star.active {
  color: #fbbf24;
}

.stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
}

.action-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: flex-start;
  flex-wrap: nowrap;
}

.action-buttons .el-button {
  min-width: auto;
  padding: 6px 12px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .song-management {
    padding: 16px;
  }
  
  .management-header {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }
  
  .search-filters {
    padding: 16px;
  }
  
  .song-table {
    padding: 16px;
  }
}
</style>