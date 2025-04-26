import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1'
})

// 请求拦截器添加token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const userApi = {
  register(data) {
    return api.post('/auth/register', data)
  },
  
  login(data) {
    return api.post('/auth/login', data)
  },
  
  updateProfile(data) {
    return api.put('/user/profile', data)
  },
  
  updatePassword(data) {
    return api.put('/user/password', data)
  },
  
  getUserInfo() {
    return api.get('/user/profile')
  },
  
  uploadAvatar(formData) {
    return api.post('/user/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
} 