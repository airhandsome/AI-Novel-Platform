// src/api/novel.js
import request from './request'

export const novelApi = {
  // 获取小说列表
  getNovels(params) {
    return request.get('/novels', { params })
  },
  getNovelList(params) {
    return request.get('/v1/novels', { params })
  },
  // 获取小说详情
  getNovelDetail(id) {
    return request.get(`/v1/novels/${id}`)
  },

  // 获取作者的小说列表
  getAuthorNovels(params) {
    return request.get('/v1/novels/author/me', { params })
  },

  // 创建新小说
  createNovel(data) {
    return request.post('/v1/novels', data)
  },

  // 更新小说信息
  updateNovel(id, data) {
    return request.put(`/v1/novels/${id}`, data)
  },

  // 删除小说
  deleteNovel(id) {
    return request.delete(`/v1/novels/${id}`)
  },

  // 获取小说章节列表
  getNovelChapters(params) {
    return request.get(`/v1/chapters/novel/${params.novelId}`, { params })
  },

  getNovel(novelId) {
    return request.get(`/v1/novels/${novelId}`)
  },

  // 获取章节详情
  getChapterDetail(chapterId) {
    return request.get(`/v1/chapters/${chapterId}`)
  },

  // 创建新章节
  createChapter(data) {
    return request.post('/v1/chapters', data)
  },

  // 更新章节内容
  updateChapter(id, data) {
    return request.put(`/v1/chapters/${id}`, data)
  },

  // 删除章节
  deleteChapter(id) {
    return request.delete(`/v1/chapters/${id}`)
  },

  // 收藏小说
  favoriteNovel: (novelId) => {
    return request({
      url: `/v1/novels/favorite/${novelId}`,
      method: 'POST'
    })
  },

  // 取消收藏
  unfavoriteNovel: (novelId) => {
    return request({
      url: `/v1/novels/favorite/${novelId}`,
      method: 'DELETE'
    })
  },

  // 获取收藏列表
  getFavorites: (params) => {
    return request({
      url: '/v1/novels/favorites',
      method: 'GET',
      params
    })
  },
  checkFavorite: (novelId) => {    
    return request.get(`/v1/novels/favorite/${novelId}`)
  },

  // 更新阅读进度
  updateReadProgress: async (data) => {
    return request({
      url: '/v1/reading-progress',
      method: 'post',
      data
    })
  },

  // 获取阅读进度
  getReadProgress: async (novelId) => {
    return request({
      url: `/v1/reading-progress/novel/${novelId}`,
      method: 'get'
    })
  },

  // AI 辅助写作
  getAIAssistance(prompt) {
    return request.post('/ai/assist', { prompt })
  },

  // 上传小说封面
  uploadCover(novelId, formData) {
    return request.post(`/novels/${novelId}/cover`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取章节
  getChapter(id) {
    return request.get(`/v1/chapters/${id}`)
  },

  // 移动章节
  moveChapter(id, direction) {
    return request.put(`/v1/chapters/${id}/move`, null, {
      params: { direction }
    })
  },

  // 更新章节状态
  updateChapterStatus(id, status) {
    return request.put(`/v1/chapters/${id}/status`, { status })
  },

  // 获取作者统计数据
  getAuthorStats() {
    return request.get('/v1/novels/author/stats')
  },

  // 获取小说大纲
  getNovelOutline(novelId) {
    return request.get(`/v1/novels/${novelId}/outline`)
  },

  // 更新小说大纲
  updateNovelOutline(novelId, data) {
    return request.put(`/v1/novels/${novelId}/outline`, data)
  },

  // 更新小说状态
  updateNovelStatus: async (novelId, status) => {
    const response = await request.put(`/v1/novels/${novelId}/status`, { status })
    return response.data
  }
}