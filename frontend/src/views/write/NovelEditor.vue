<template>
  <div class="novel-editor">
    <div class="editor-header">
      <div class="header-left">
        <el-input
          v-model="chapterTitle"
          placeholder="请输入章节标题"
          class="title-input"
        />
        <div class="word-count">
          <el-icon><Document /></el-icon>
          {{ wordCount }} 字
        </div>
      </div>
      <div class="header-right">
        <el-button class="ai-button" @click="openAIDialog">
          <el-icon><Brush /></el-icon>AI助手
        </el-button>
        <el-switch
          v-model="autoSave"
          active-text="自动保存"
          class="auto-save-switch"
        />
        <el-radio-group v-model="editorMode" size="large" class="mode-switch">
          <el-radio-button value="rich">富文本</el-radio-button>
          <el-radio-button value="markdown">Markdown</el-radio-button>
        </el-radio-group>
        <el-button type="primary" @click="saveChapter" :loading="saving">
          保存
        </el-button>
      </div>
    </div>
    
    <div class="editor-content" :class="{ 'markdown-mode': editorMode === 'markdown' }">
      <!-- 富文本编辑器 -->
      <div v-show="editorMode === 'rich'" class="rich-editor">
        <Toolbar
          :editor="editorRef"
          :defaultConfig="toolbarConfig"
          mode="default"
          class="editor-toolbar"
        />
        <Editor
          v-model="content"
          :defaultConfig="editorConfig"
          mode="default"
          @onCreated="handleCreated"
          @onChange="handleContentChange"
          class="editor-container"
        />
      </div>

      <!-- Markdown编辑器 -->
      <div v-show="editorMode === 'markdown'" class="markdown-editor">
        <div class="markdown-container">
          <div class="markdown-input">
            <el-input
              v-model="content"
              type="textarea"
              :rows="20"
              placeholder="请输入 Markdown 内容"
              @input="handleContentChange"
            />
          </div>
          <div class="markdown-preview" v-html="markdownHtml"></div>
        </div>
      </div>
    </div>

    <!-- 保存提示
    <transition name="fade">
      <div v-if="showSaveStatus" class="save-status" :class="{ 'success': saveSuccess }">
        {{ saveStatusText }}
      </div>
    </transition> -->

    <!-- AI助手对话框 -->
    <el-dialog 
      v-model="aiDialogVisible" 
      title="AI创作助手" 
      width="800px"
      class="ai-dialog custom-dialog"
      :close-on-click-modal="false"
    >
      <div class="ai-chat">
        <div class="chat-messages" ref="chatMessages">
          <div 
            v-for="(message, index) in aiMessages" 
            :key="index"
            :class="['message', message.role]"
          >
            <el-avatar 
              :size="40" 
              :src="message.role === 'assistant' ? '/ai-avatar.png' : '/user-avatar.png'"
            />
            <div class="message-content">
              <div class="message-text" v-html="formatMessage(message.content)"></div>
              <div class="message-actions" v-if="message.role === 'assistant' && message.suggestions">
                <el-button 
                  v-for="(text, idx) in message.suggestions" 
                  :key="idx"
                  size="small"
                  @click="insertText(text)"
                >
                  插入文本
                </el-button>
              </div>
              <div class="message-time">{{ formatTime(message.timestamp) }}</div>
            </div>
          </div>
        </div>
        <div class="quick-actions">
          <el-button-group>
            <el-button size="small" @click="requestSuggestion('plot')">
              <el-icon><Connection /></el-icon>情节发展
            </el-button>
            <el-button size="small" @click="requestSuggestion('character')">
              <el-icon><User /></el-icon>人物刻画
            </el-button>
            <el-button size="small" @click="requestSuggestion('scene')">
              <el-icon><Picture /></el-icon>场景描写
            </el-button>
            <el-button size="small" @click="requestSuggestion('dialogue')">
              <el-icon><ChatDotRound /></el-icon>对话生成
            </el-button>
          </el-button-group>
        </div>
        <div class="chat-input">
          <el-input
            v-model="aiInput"
            type="textarea"
            :rows="3"
            placeholder="描述你的需求，例如：'帮我构思下一段剧情'..."
            @keyup.enter.ctrl="sendMessage"
          />
          <div class="input-actions">
            <span class="input-tip">Ctrl + Enter 发送</span>
            <el-button type="primary" @click="sendMessage" :loading="isAiLoading">
              发送
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, shallowRef, onMounted, onBeforeUnmount, computed, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Document, Brush, User, Picture, ChatDotRound,
  Connection 
} from '@element-plus/icons-vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import '@wangeditor/editor/dist/css/style.css'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { novelApi } from '../../api/novel'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()
const novelId = route.params.id
const chapterId = route.query.chapterId

// 编辑器相关
const editorRef = shallowRef()
const content = ref('')
const chapterTitle = ref('')
const wordCount = ref(0)
const autoSave = ref(true)
const editorMode = ref('rich')
const saving = ref(false)
const saveSuccess = ref(true)
const saveStatusText = ref('')

let autoSaveInterval = null
const md = new MarkdownIt({
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch (__) {}
    }
    return ''
  }
})

// 富文本编辑器配置
const toolbarConfig = {
  excludeKeys: [
    'uploadImage',
    'uploadVideo',
    'insertTable',
    'group-video',
    'group-image'
  ]
}

const editorConfig = {
  placeholder: '请输入内容...',
  autoFocus: false
}

// Markdown 预览
const markdownHtml = computed(() => {
  if (editorMode.value === 'markdown') {
    return md.render(content.value || '')
  }
  return ''
})

// 计算字数
const calculateWordCount = (text) => {
  if (!text) return 0
  // 移除 HTML 标签和空白字符后计算字数
  return text.replace(/<[^>]+>/g, '').replace(/\s+/g, '').length
}

// 内容变化处理
const handleContentChange = () => {
  wordCount.value = calculateWordCount(content.value)
  showAutoSaveStatus()
}

// 编辑器创建完成
const handleCreated = (editor) => {
  editorRef.value = editor
}

// 显示保存状态
const showSaveStatus = (message, success = true) => {
  saveStatusText.value = message
  saveSuccess.value = success
  showSaveStatus.value = true
  setTimeout(() => {
    showSaveStatus.value = false
  }, 6000)
}

// 自动保存状态提示
const showAutoSaveStatus = () => {
  if (autoSave.value) {
    showSaveStatus('正在自动保存...')
  }
}

// 保存章节
const saveChapter = async () => {
  if (saving.value) return
  
  try {
    if (!chapterTitle.value.trim()) {
      ElMessage.warning('请输入章节标题')
      return
    }
    
    saving.value = true
    const data = {
      novelId: Number(novelId),
      title: chapterTitle.value,
      content: content.value,
      wordCount: wordCount.value
    }
    
    if (chapterId) {
      await novelApi.updateChapter(chapterId, data)
    } else {
      await novelApi.createChapter(data)
    }
    
    showSaveStatus('保存成功')
    
    // 如果是手动保存（不是自动保存），则跳转到章节列表
    router.push(`/write/chapters/${novelId}`)
  } catch (error) {
    console.error('保存失败:', error)
    showSaveStatus('保存失败', false)
  } finally {
    saving.value = false
  }
}

// 自动保存
const startAutoSave = () => {
  if (autoSaveInterval) {
    clearInterval(autoSaveInterval)
  }
  
  autoSaveInterval = setInterval(() => {
    if (autoSave.value && content.value?.trim()) {
      saveChapter()
    }
  }, 60000) // 改为60秒
}

// 监听自动保存开关变化
watch(autoSave, (newValue) => {
  if (newValue) {
    startAutoSave()
  } else {
    if (autoSaveInterval) {
      clearInterval(autoSaveInterval)
      autoSaveInterval = null
    }
  }
})




// 加载章节内容
const loadChapter = async () => {
  if (!chapterId) return
  
  try {
    const response = await novelApi.getChapter(chapterId)
    chapterTitle.value = response.title
    content.value = response.content
    handleContentChange()
  } catch (error) {
    console.error('加载章节失败:', error)
    ElMessage.error('加载章节失败')
  }
}

// AI 助手相关
const aiDialogVisible = ref(false)
const aiMessages = ref([])
const aiInput = ref('')
const isAiLoading = ref(false)
const chatMessages = ref(null)

// 打开AI对话框
const openAIDialog = () => {
  aiDialogVisible.value = true
  if (aiMessages.value.length === 0) {
    aiMessages.value.push({
      role: 'assistant',
      content: '你好！我是你的AI创作助手。我可以帮你：\n\n- 构思情节发展\n- 优化人物刻画\n- 完善场景描写\n- 生成自然对话\n\n请告诉我你需要什么帮助？',
      timestamp: Date.now()
    })
  }
}

// 发送消息
const sendMessage = async () => {
  if (!aiInput.value.trim() || isAiLoading.value) return
  
  const userMessage = aiInput.value.trim()
  aiMessages.value.push({
    role: 'user',
    content: userMessage,
    timestamp: Date.now()
  })
  aiInput.value = ''
  
  isAiLoading.value = true
  try {
    // 调用DeepSeek API
    const response = await fetch('https://api.deepseek.com/v1/chat/completions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${import.meta.env.VITE_DEEPSEEK_API_KEY}`
      },
      body: JSON.stringify({
        model: 'deepseek-chat',
        messages: [
          { 
            role: 'system', 
            content: '你是一个专业的小说创作助手，擅长提供具体的写作建议和示例文本。当用户请求帮助时，尽可能提供具体的文本建议，并在回复中包含可以直接使用的文本片段。' 
          },
          { role: 'user', content: userMessage }
        ]
      })
    })
    
    const data = await response.json()
    const assistantMessage = {
      role: 'assistant',
      content: data.choices[0].message.content,
      timestamp: Date.now(),
      suggestions: extractSuggestions(data.choices[0].message.content)
    }
    aiMessages.value.push(assistantMessage)
  } catch (error) {
    ElMessage.error('AI响应失败，请稍后重试')
  } finally {
    isAiLoading.value = false
    await nextTick()
    scrollToBottom()
  }
}

// 快捷请求建议
const requestSuggestion = async (type) => {
  const prompts = {
    plot: '请帮我构思接下来的情节发展，考虑到目前的内容：\n\n' + content.value,
    character: '请帮我加强人物刻画，基于当前场景：\n\n' + content.value,
    scene: '请帮我完善场景描写，丰富环境细节：\n\n' + content.value,
    dialogue: '请帮我生成自然的对话，体现人物性格：\n\n' + content.value
  }
  
  aiInput.value = prompts[type]
  await sendMessage()
}

// 从AI回复中提取可用的文本建议
const extractSuggestions = (content) => {
  const suggestions = []
  const regex = /```([\s\S]*?)```/g
  let match
  
  while ((match = regex.exec(content)) !== null) {
    suggestions.push(match[1].trim())
  }
  
  return suggestions
}

// 插入AI建议的文本
const insertText = (text) => {
  if (editorMode.value === 'rich' && editorRef.value) {
    editorRef.value.insertText(text)
  } else {
    content.value += '\n' + text
  }
  ElMessage.success('文本已插入')
}

const formatMessage = (content) => {
  return marked.parse(content)
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString()
}

const scrollToBottom = () => {
  if (chatMessages.value) {
    chatMessages.value.scrollTop = chatMessages.value.scrollHeight
  }
}

onMounted(() => {
  loadChapter()
  startAutoSave()
})

onBeforeUnmount(() => {
  if (autoSaveInterval) {
    clearInterval(autoSaveInterval)
  }
  if (editorRef.value) {
    editorRef.value.destroy()
  }
})
</script>

<style scoped>
.novel-editor {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: 100vh;
  background: #fff;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.title-input {
  width: 300px;
}

.title-input :deep(.el-input__inner) {
  font-size: 18px;
  border: none;
  border-radius: 6px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.03);
}

.word-count {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #666;
  font-size: 14px;
}

.mode-switch {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
}

.editor-content {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.rich-editor {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 200px);
}

.editor-toolbar {
  border-bottom: 1px solid #eee;
  background: #f9f9f9;
}

.editor-container {
  flex: 1;
  overflow-y: auto;
}

.markdown-editor {
  height: calc(100vh - 200px);
}

.markdown-container {
  display: flex;
  height: 100%;
  gap: 20px;
}

.markdown-input,
.markdown-preview {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.markdown-preview {
  background: #f9f9f9;
  border-radius: 8px;
}

.save-status {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  border-radius: 20px;
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
  font-size: 14px;
  transition: all 0.3s ease;
}

.save-status.success {
  background: rgba(52, 199, 89, 0.9);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

:deep(.w-e-toolbar) {
  border: none !important;
  padding: 10px !important;
}

:deep(.w-e-text-container) {
  border: none !important;
}

:deep(.markdown-preview) {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  line-height: 1.6;
}

:deep(.markdown-preview h1),
:deep(.markdown-preview h2),
:deep(.markdown-preview h3),
:deep(.markdown-preview h4),
:deep(.markdown-preview h5),
:deep(.markdown-preview h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

:deep(.markdown-preview code) {
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
  padding: 0.2em 0.4em;
  font-size: 85%;
  font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace;
}

:deep(.markdown-preview pre) {
  background-color: #f6f8fa;
  border-radius: 6px;
  padding: 16px;
  overflow: auto;
}

:deep(.markdown-preview blockquote) {
  margin: 0;
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
}

.ai-dialog {
  .ai-chat {
    height: 600px;
    display: flex;
    flex-direction: column;
  }

  .chat-messages {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    background: var(--system-grouped-secondary);
    border-radius: 12px;
    margin-bottom: 20px;
  }

  .message {
    display: flex;
    gap: 16px;
    margin-bottom: 24px;

    &.assistant {
      .message-content {
        background: var(--system-blue-tint);
      }
    }

    &.user {
      flex-direction: row-reverse;

      .message-content {
        background: var(--system-grouped-primary);
      }
    }
  }

  .message-content {
    flex: 1;
    padding: 16px;
    border-radius: 12px;
    max-width: 80%;
  }

  .message-text {
    color: var(--text-primary);
    line-height: 1.6;

    :deep(p) {
      margin: 0 0 12px;
      &:last-child {
        margin-bottom: 0;
      }
    }

    :deep(code) {
      background: var(--system-grouped-tertiary);
      padding: 2px 6px;
      border-radius: 4px;
    }
  }

  .message-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
  }

  .message-time {
    font-size: 12px;
    color: var(--text-tertiary);
    margin-top: 8px;
  }

  .quick-actions {
    margin-bottom: 16px;
    display: flex;
    justify-content: center;
  }

  .chat-input {
    background: var(--system-background);
    padding: 16px;
    border-radius: 12px;
  }

  .input-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
  }

  .input-tip {
    font-size: 12px;
    color: var(--text-tertiary);
  }
}

.ai-button {
  display: flex;
  align-items: center;
  gap: 6px;
}
</style> 