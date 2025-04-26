<template>
  <div class="novel-settings">
    <div class="header">
      <h1>作品设置</h1>
      <el-button type="primary" @click="saveSettings">保存设置</el-button>
    </div>

    <el-form
      ref="formRef"
      :model="novelForm"
      :rules="rules"
      label-width="100px"
      class="settings-form"
    >
      <el-form-item label="作品标题" prop="title">
        <el-input v-model="novelForm.title" placeholder="请输入作品标题" />
      </el-form-item>

      <el-form-item label="作品分类" prop="category">
        <el-select v-model="novelForm.category" placeholder="请选择作品分类">
          <el-option
            v-for="item in categories"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="作品状态" prop="status">
        <el-select v-model="novelForm.status" placeholder="请选择作品状态">
          <el-option
            v-for="item in statusOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="作品简介" prop="description">
        <el-input
          v-model="novelForm.description"
          type="textarea"
          :rows="4"
          placeholder="请输入作品简介"
        />
      </el-form-item>

      <el-form-item label="作品标签">
        <el-tag
          v-for="tag in novelForm.tags"
          :key="tag"
          class="tag-item"
          closable
          @close="handleRemoveTag(tag)"
        >
          {{ tag }}
        </el-tag>
        <el-input
          v-if="inputTagVisible"
          ref="tagInputRef"
          v-model="inputTagValue"
          class="tag-input"
          size="small"
          @keyup.enter="handleInputConfirm"
          @blur="handleInputConfirm"
        />
        <el-button v-else class="button-new-tag" size="small" @click="showTagInput">
          + 新增标签
        </el-button>
      </el-form-item>

      <el-form-item label="封面图片">
        <el-upload
          class="cover-uploader"
          :show-file-list="false"
          :before-upload="beforeCoverUpload"
          :http-request="uploadCover"
        >
          <img v-if="novelForm.coverUrl" :src="novelForm.coverUrl" class="cover-image" />
          <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
        </el-upload>
        <div class="upload-tip">建议尺寸: 300x400px，大小不超过2MB</div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { novelApi } from '../../api/novel'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const tagInputRef = ref(null)
const inputTagVisible = ref(false)
const inputTagValue = ref('')

const novelForm = ref({
  title: '',
  category: '',
  status: 0,
  description: '',
  tags: [],
  coverUrl: ''
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入作品标题', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择作品分类', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入作品简介', trigger: 'blur' },
    { min: 10, max: 500, message: '长度在 10 到 500 个字符', trigger: 'blur' }
  ]
}

// 分类选项
const categories = [
  { value: 'fantasy', label: '奇幻' },
  { value: 'scifi', label: '科幻' },
  { value: 'romance', label: '言情' },
  { value: 'mystery', label: '悬疑' },
  { value: 'history', label: '历史' },
  { value: 'other', label: '其他' }
]

// 状态选项
const statusOptions = [
  { value: 0, label: '连载中' },
  { value: 1, label: '已完结' },
  { value: 2, label: '暂停更新' }
]

// 加载小说信息
const loadNovel = async () => {
  try {
    const response = await novelApi.getNovel(route.params.id)
    const novel = response
    novelForm.value = {
      title: novel.title || '',
      category: novel.category || '',
      status: novel.status || 0,
      description: novel.description || '',
      tags: novel.tags || [],
      coverUrl: novel.coverUrl || ''
    }
  } catch (error) {
    console.error('加载小说信息失败:', error)
    ElMessage.error('加载小说信息失败')
  }
}

// 保存设置
const saveSettings = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await novelApi.updateNovel(route.params.id, novelForm.value)
        ElMessage.success('保存成功')
        router.push('/write')
      } catch (error) {
        console.error('保存失败:', error)
        ElMessage.error('保存失败')
      }
    }
  })
}

// 标签相关方法
const handleRemoveTag = (tag) => {
  novelForm.value.tags = novelForm.value.tags.filter(t => t !== tag)
}

const showTagInput = () => {
  inputTagVisible.value = true
  nextTick(() => {
    tagInputRef.value?.focus()
  })
}

const handleInputConfirm = () => {
  if (inputTagValue.value) {
    if (!novelForm.value.tags.includes(inputTagValue.value)) {
      novelForm.value.tags.push(inputTagValue.value)
    }
  }
  inputTagVisible.value = false
  inputTagValue.value = ''
}

// 封面上传相关方法
const beforeCoverUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

const uploadCover = async (options) => {
  try {
    const formData = new FormData()
    formData.append('cover', options.file)
    const response = await novelApi.uploadCover(route.params.id, formData)
    novelForm.value.coverUrl = response.data.coverUrl
    ElMessage.success('封面上传成功')
  } catch (error) {
    console.error('封面上传失败:', error)
    ElMessage.error('封面上传失败')
  }
}

onMounted(() => {
  loadNovel()
})
</script>

<style scoped>
.novel-settings {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.settings-form {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.tag-item {
  margin-right: 10px;
  margin-bottom: 10px;
}

.tag-input {
  width: 100px;
  margin-right: 10px;
  vertical-align: bottom;
}

.button-new-tag {
  margin-bottom: 10px;
}

.cover-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 300px;
  height: 400px;
}

.cover-uploader:hover {
  border-color: #409eff;
}

.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 300px;
  height: 400px;
  line-height: 400px;
  text-align: center;
}

.cover-image {
  width: 300px;
  height: 400px;
  object-fit: cover;
  display: block;
}

.upload-tip {
  font-size: 12px;
  color: #666;
  margin-top: 10px;
}
</style> 