<template>
  <div class="novel-outline">
    <div class="header">
      <div class="title-section">
        <h1>大纲管理</h1>
        <div class="statistics">
          <el-tag type="info">总字数: {{ totalWords }}</el-tag>
          <el-tag type="success">完成进度: {{ completionRate }}%</el-tag>
        </div>
      </div>
      <div class="actions">
        <el-button type="primary" @click="saveOutline">保存大纲</el-button>
        <el-dropdown split-button type="primary" @command="handleAdd">
          添加
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="volume">添加卷</el-dropdown-item>
              <el-dropdown-item command="chapter">添加章节</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <div class="outline-content">
      <el-tabs v-model="activeTab" class="outline-tabs">
        <el-tab-pane label="章节大纲" name="chapters">
          <div class="outline-tree">
            <el-tree
              ref="treeRef"
              :data="outlineData"
              :props="defaultProps"
              node-key="id"
              default-expand-all
              draggable
              :allow-drop="allowDrop"
              :allow-drag="allowDrag"
              @node-drag-end="handleDragEnd"
            >
              <template #default="{ node, data }">
                <div class="custom-tree-node">
                  <div class="node-content">
                    <el-tag 
                      size="small" 
                      :type="data.type === 'volume' ? 'warning' : 'info'"
                      class="node-type"
                    >
                      {{ data.type === 'volume' ? '卷' : '章' }}
                    </el-tag>
                    <span class="node-title">{{ data.title }}</span>
                    <span v-if="data.type === 'chapter'" class="node-info">
                      <el-tag size="small" :type="getStatusType(data.status)">
                        {{ getStatusText(data.status) }}
                      </el-tag>
                      <span class="word-count">{{ data.wordCount || 0 }}字</span>
                    </span>
                  </div>
                  <div class="node-actions">
                    <el-button-group>
                      <el-button 
                        v-if="data.type === 'volume'"
                        link 
                        type="primary" 
                        @click.stop="addSection(data)"
                      >
                        添加章节
                      </el-button>
                      <el-button 
                        link 
                        type="primary" 
                        @click.stop="editNode(node, data)"
                      >
                        编辑
                      </el-button>
                      <el-button 
                        link 
                        type="danger" 
                        @click.stop="removeNode(node, data)"
                      >
                        删除
                      </el-button>
                    </el-button-group>
                  </div>
                </div>
              </template>
            </el-tree>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="世界观设定" name="worldbuilding">
          <div class="world-building">
            <el-collapse v-model="activeSettings">
              <el-collapse-item title="背景设定" name="background">
                <el-input
                  v-model="worldBuilding.background"
                  type="textarea"
                  :rows="4"
                  placeholder="描述故事的背景设定..."
                />
              </el-collapse-item>
              
              <el-collapse-item title="主要人物" name="characters">
                <div v-for="(char, index) in worldBuilding.characters" :key="index" class="character-item">
                  <div class="character-header">
                    <el-input v-model="char.name" placeholder="人物名称" class="char-name" />
                    <el-button link type="danger" @click="removeCharacter(index)">删除</el-button>
                  </div>
                  <el-input
                    v-model="char.description"
                    type="textarea"
                    :rows="3"
                    placeholder="人物描述..."
                  />
                </div>
                <el-button text @click="addCharacter">
                  <el-icon><Plus /></el-icon>添加人物
                </el-button>
              </el-collapse-item>
              
              <el-collapse-item title="地理环境" name="locations">
                <div v-for="(loc, index) in worldBuilding.locations" :key="index" class="location-item">
                  <div class="location-header">
                    <el-input v-model="loc.name" placeholder="地点名称" class="loc-name" />
                    <el-button link type="danger" @click="removeLocation(index)">删除</el-button>
                  </div>
                  <el-input
                    v-model="loc.description"
                    type="textarea"
                    :rows="3"
                    placeholder="地点描述..."
                  />
                </div>
                <el-button text @click="addLocation">
                  <el-icon><Plus /></el-icon>添加地点
                </el-button>
              </el-collapse-item>
            </el-collapse>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 编辑节点对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="editingNode.type === 'volume' ? '编辑卷' : '编辑章节'"
      width="50%"
    >
      <el-form :model="editingNode" label-width="80px">
        <el-form-item label="标题" required>
          <el-input v-model="editingNode.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="editingNode.description"
            type="textarea"
            :rows="4"
            placeholder="请输入描述"
          />
        </el-form-item>
        <template v-if="editingNode.type === 'chapter'">
          <el-form-item label="状态">
            <el-select v-model="editingNode.status">
              <el-option
                v-for="item in statusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="字数">
            <el-input-number 
              v-model="editingNode.wordCount" 
              :min="0"
              :step="100"
            />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveNode">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { novelApi } from '../../api/novel'

const route = useRoute()
const activeTab = ref('chapters')
const activeSettings = ref(['background'])
const treeRef = ref(null)
const editDialogVisible = ref(false)
const editingNode = ref({})
const isNewNode = ref(false)

// 大纲数据
const outlineData = ref([])
const defaultProps = {
  children: 'children',
  label: 'title'
}

// 世界观设定数据
const worldBuilding = ref({
  background: '',
  characters: [],
  locations: []
})

// 状态选项
const statusOptions = [
  { value: 0, label: '未开始' },
  { value: 1, label: '进行中' },
  { value: 2, label: '已完成' }
]

// 计算总字数
const totalWords = computed(() => {
  let total = 0
  const countWords = (nodes) => {
    for (const node of nodes) {
      if (node.wordCount) total += node.wordCount
      if (node.children?.length) countWords(node.children)
    }
  }
  countWords(outlineData.value)
  return total
})

// 计算完成进度
const completionRate = computed(() => {
  let total = 0
  let completed = 0
  const countChapters = (nodes) => {
    for (const node of nodes) {
      if (node.type === 'chapter') {
        total++
        if (node.status === 2) completed++
      }
      if (node.children?.length) countChapters(node.children)
    }
  }
  countChapters(outlineData.value)
  return total ? Math.round((completed / total) * 100) : 0
})

// 加载大纲数据
const loadOutline = async () => {
  try {
    const response = await novelApi.getNovelOutline(route.params.id)
    outlineData.value = response.data.outline || []
    worldBuilding.value = response.data.worldBuilding || {
      background: '',
      characters: [],
      locations: []
    }
  } catch (error) {
    console.error('加载大纲失败:', error)
    ElMessage.error('加载大纲失败')
  }
}

// 保存大纲
const saveOutline = async () => {
  try {
    await novelApi.updateNovelOutline(route.params.id, {
      outline: outlineData.value,
      worldBuilding: worldBuilding.value
    })
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  }
}

// 添加卷
const addVolume = () => {
  editingNode.value = {
    id: Date.now(),
    title: '',
    description: '',
    type: 'volume',
    children: []
  }
  isNewNode.value = true
  editDialogVisible.value = true
}

// 添加章节
const addSection = (parentNode = null) => {
  editingNode.value = {
    id: Date.now(),
    title: '',
    description: '',
    type: 'chapter',
    status: 0,
    wordCount: 0,
    children: []
  }
  isNewNode.value = true
  editDialogVisible.value = true
  editingNode.value.parentId = parentNode?.id
}

// 编辑节点
const editNode = (node, data) => {
  editingNode.value = { ...data }
  isNewNode.value = false
  editDialogVisible.value = true
}

// 保存节点
const saveNode = () => {
  if (!editingNode.value.title.trim()) {
    ElMessage.warning('标题不能为空')
    return
  }

  if (isNewNode.value) {
    // 添加新节点
    if (editingNode.value.parentId) {
      // 添加到指定父节点
      const parent = findNode(outlineData.value, editingNode.value.parentId)
      if (parent) {
        parent.children = parent.children || []
        parent.children.push({ ...editingNode.value })
      }
    } else {
      // 添加到根级别
      outlineData.value.push({ ...editingNode.value })
    }
  } else {
    // 更新现有节点
    const node = findNode(outlineData.value, editingNode.value.id)
    if (node) {
      Object.assign(node, editingNode.value)
    }
  }

  editDialogVisible.value = false
  saveOutline()
  ElMessage.success('保存成功')
}

// 删除节点
const removeNode = async (node, data) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除${data.type === 'volume' ? '这一卷' : '这一章'}吗？${
        data.children?.length ? '其下所有内容都会被删除！' : ''
      }`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const parent = node.parent
    const children = parent.data.children || outlineData.value
    const index = children.findIndex(d => d.id === data.id)
    children.splice(index, 1)
    saveOutline()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 拖拽相关方法
const allowDrop = (draggingNode, dropNode, type) => {
  // 卷只能放在根级
  if (draggingNode.data.type === 'volume' && dropNode.level !== 0) {
    return false
  }
  // 章节只能放在卷下或根级
  if (draggingNode.data.type === 'chapter') {
    if (type === 'inner' && dropNode.data.type !== 'volume') {
      return false
    }
  }
  return true
}

const allowDrag = (draggingNode) => {
  return true
}

const handleDragEnd = async (draggingNode, dropNode, dropType) => {
  // 更新排序
  await saveOutline()
  ElMessage.success('更新排序成功')
}

// 查找节点
const findNode = (nodes, id) => {
  for (const node of nodes) {
    if (node.id === id) return node
    if (node.children) {
      const found = findNode(node.children, id)
      if (found) return found
    }
  }
  return null
}

// 世界观设定相关方法
const addCharacter = () => {
  worldBuilding.value.characters.push({
    name: '',
    description: ''
  })
}

const removeCharacter = (index) => {
  worldBuilding.value.characters.splice(index, 1)
}

const addLocation = () => {
  worldBuilding.value.locations.push({
    name: '',
    description: ''
  })
}

const removeLocation = (index) => {
  worldBuilding.value.locations.splice(index, 1)
}

// 获取状态类型
const getStatusType = (status) => {
  const types = {
    0: 'info',    // 未开始
    1: 'warning', // 进行中
    2: 'success'  // 已完成
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    0: '未开始',
    1: '进行中',
    2: '已完成'
  }
  return texts[status] || '未知'
}

// 处理添加按钮的下拉菜单命令
const handleAdd = (command) => {
  switch (command) {
    case 'volume':
      addVolume()
      break
    case 'chapter':
      addSection()
      break
  }
}

onMounted(() => {
  loadOutline()
})
</script>

<style scoped>
.novel-outline {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.statistics {
  display: flex;
  gap: 10px;
}

.actions {
  display: flex;
  gap: 10px;
}

.outline-content {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px;
}

.node-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.node-type {
  min-width: 32px;
  text-align: center;
}

.node-title {
  font-size: 14px;
  font-weight: 500;
}

.node-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.word-count {
  color: #666;
  font-size: 12px;
}

.node-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.custom-tree-node:hover .node-actions {
  opacity: 1;
}

.world-building {
  margin-top: 20px;
}

.character-item,
.location-item {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.character-header,
.location-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.char-name,
.loc-name {
  width: 200px;
}

:deep(.el-tree-node__content) {
  height: auto;
  padding: 4px 0;
}

:deep(.el-tree-node__children) {
  padding-left: 24px;
}
</style> 