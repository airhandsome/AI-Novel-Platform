<template>
  <div class="novel-outline">
    <el-tabs>
      <el-tab-pane label="章节大纲">
        <div class="outline-header">
          <h2>章节大纲</h2>
          <el-button type="primary" @click="addOutlineItem">添加大纲条目</el-button>
        </div>
        
        <el-tree
          v-if="outlineData.outline && outlineData.outline.length > 0"
          :data="outlineData.outline"
          node-key="id"
          default-expand-all
          draggable
          :allow-drop="allowDrop"
          :allow-drag="allowDrag"
          @node-drag-end="handleDragEnd"
        >
          <template #default="{ node, data }">
            <div class="outline-node">
              <div class="node-content">
                <span class="node-title">{{ data.title }}</span>
                <div class="node-actions">
                  <el-button link type="primary" @click="editItem(data)">编辑</el-button>
                  <el-button link type="danger" @click="deleteItem(node, data)">删除</el-button>
                </div>
              </div>
              <div v-if="data.description" class="node-description">
                {{ data.description }}
              </div>
            </div>
          </template>
        </el-tree>
        
        <el-empty v-else description="暂无大纲内容" />
      </el-tab-pane>

      <el-tab-pane label="世界观设定">
        <div class="world-building">
          <h3>背景设定</h3>
          <el-input
            v-model="outlineData.worldBuilding.background"
            type="textarea"
            rows="4"
            placeholder="请输入世界观背景"
            @change="saveOutline"
          />

          <h3 class="mt-4">角色设定</h3>
          <div class="character-list">
            <el-card v-for="(char, index) in outlineData.worldBuilding.characters" :key="index" class="mb-3">
              <div class="character-item">
                <div class="character-header">
                  <el-input v-model="char.name" placeholder="角色名称" @change="saveOutline" />
                  <el-button type="danger" link @click="deleteCharacter(index)">删除</el-button>
                </div>
                <el-input
                  v-model="char.description"
                  type="textarea"
                  rows="2"
                  placeholder="角色描述"
                  @change="saveOutline"
                />
              </div>
            </el-card>
            <el-button type="primary" @click="addCharacter">添加角色</el-button>
          </div>

          <h3 class="mt-4">地点设定</h3>
          <div class="location-list">
            <el-card v-for="(loc, index) in outlineData.worldBuilding.locations" :key="index" class="mb-3">
              <div class="location-item">
                <div class="location-header">
                  <el-input v-model="loc.name" placeholder="地点名称" @change="saveOutline" />
                  <el-button type="danger" link @click="deleteLocation(index)">删除</el-button>
                </div>
                <el-input
                  v-model="loc.description"
                  type="textarea"
                  rows="2"
                  placeholder="地点描述"
                  @change="saveOutline"
                />
              </div>
            </el-card>
            <el-button type="primary" @click="addLocation">添加地点</el-button>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 编辑大纲条目对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingItem ? '编辑大纲' : '添加大纲'"
      width="50%"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            rows="4"
            placeholder="请输入描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveOutlineItem">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { novelApi } from '@/api/novel'

const route = useRoute()
const novelId = route.params.id

const outlineData = ref({
  outline: [],
  worldBuilding: {
    background: '',
    characters: [],
    locations: []
  }
})
const dialogVisible = ref(false)
const editingItem = ref(null)
const form = ref({
  title: '',
  description: '',
})

// 加载大纲数据
const loadOutline = async () => {
  try {
    const response = await novelApi.getNovelOutline(novelId)
    if (response.data) {
      outlineData.value = response.data
    }
  } catch (error) {
    console.error('加载大纲失败:', error)
    ElMessage.error('加载大纲失败')
  }
}

// 添加大纲条目
const addOutlineItem = () => {
  editingItem.value = null
  form.value = {
    title: '',
    description: '',
  }
  dialogVisible.value = true
}

// 编辑大纲条目
const editItem = (item) => {
  editingItem.value = item
  form.value = {
    title: item.title,
    description: item.description,
  }
  dialogVisible.value = true
}

// 删除大纲条目
const deleteItem = async (node, data) => {
  try {
    const parent = node.parent
    const children = parent.data.children || outlineData.value.outline
    const index = children.findIndex(item => item.id === data.id)
    if (index > -1) {
      children.splice(index, 1)
      await saveOutline()
      ElMessage.success('删除成功')
    }
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  }
}

// 保存大纲条目
const saveOutlineItem = async () => {
  if (!form.value.title) {
    ElMessage.warning('请输入标题')
    return
  }

  try {
    if (editingItem.value) {
      // 更新现有条目
      editingItem.value.title = form.value.title
      editingItem.value.description = form.value.description
    } else {
      // 添加新条目
      const newItem = {
        id: Date.now(), // 临时ID
        title: form.value.title,
        description: form.value.description,
        children: [],
      }
      outlineData.value.outline.push(newItem)
    }

    await saveOutline()
    dialogVisible.value = false
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  }
}

// 添加角色
const addCharacter = () => {
  outlineData.value.worldBuilding.characters.push({
    name: '',
    description: ''
  })
}

// 删除角色
const deleteCharacter = (index) => {
  outlineData.value.worldBuilding.characters.splice(index, 1)
  saveOutline()
}

// 添加地点
const addLocation = () => {
  outlineData.value.worldBuilding.locations.push({
    name: '',
    description: ''
  })
}

// 删除地点
const deleteLocation = (index) => {
  outlineData.value.worldBuilding.locations.splice(index, 1)
  saveOutline()
}

// 保存整个大纲到服务器
const saveOutline = async () => {
  try {
    await novelApi.updateNovelOutline(novelId, outlineData.value)
  } catch (error) {
    console.error('保存大纲失败:', error)
    throw error
  }
}

// 拖拽相关方法
const allowDrop = (draggingNode, dropNode, type) => {
  return true // 允许所有拖拽
}

const allowDrag = (draggingNode) => {
  return true // 允许所有节点拖拽
}

const handleDragEnd = async (draggingNode, dropNode, dropType, ev) => {
  try {
    await saveOutline()
    ElMessage.success('更新排序成功')
  } catch (error) {
    console.error('更新排序失败:', error)
    ElMessage.error('更新排序失败')
    await loadOutline() // 重新加载以恢复原始顺序
  }
}

onMounted(() => {
  loadOutline()
})
</script>

<style scoped>
.novel-outline {
  padding: 20px;
}

.outline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.outline-node {
  flex: 1;
  padding: 8px 0;
}

.node-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.node-title {
  font-weight: bold;
}

.node-description {
  margin-top: 4px;
  color: #666;
  font-size: 14px;
}

.node-actions {
  opacity: 0;
  transition: opacity 0.3s;
}

.outline-node:hover .node-actions {
  opacity: 1;
}

.el-tree-node__content {
  height: auto !important;
  padding: 4px 0;
}

.world-building {
  padding: 20px;
}

.character-list,
.location-list {
  margin-top: 16px;
}

.character-item,
.location-item {
  margin-bottom: 16px;
}

.character-header,
.location-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.mt-4 {
  margin-top: 24px;
}

.mb-3 {
  margin-bottom: 16px;
}
</style> 