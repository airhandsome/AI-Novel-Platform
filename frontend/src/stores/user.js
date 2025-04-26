import { defineStore } from 'pinia'
import { userApi } from '../api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token')
  }),
  
  actions: {
    async login(credentials) {
      try {
        const { data } = await userApi.login(credentials)
        this.token = data.token
        this.user = data.user
        localStorage.setItem('token', data.token)
        return data
      } catch (error) {
        throw error
      }
    },
    
    async register(userData) {
      try {
        const { data } = await userApi.register(userData)
        return data
      } catch (error) {
        throw error
      }
    },
    
    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
    },

    updateUser(userData) {
      this.user = {
        ...this.user,
        ...userData
      }
    }
  }
}) 