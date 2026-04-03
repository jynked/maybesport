import { defineStore } from 'pinia'
import { http } from '../api'
import router from '../router'
import i18n from '../i18n'
import { useFavouritesStore } from './favourites'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    user: JSON.parse(localStorage.getItem('user')) || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    async register(email, password) {
      try {
        const response = await http.post('/auth/register', { email, password })
        const { token, data } = response.data

        this.token = token
        this.user = data
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(data))

        await this.executePendingAction()

        const redirect = router.currentRoute.value.query.redirect || '/user'
        router.push(redirect)
      } catch (error) {
        console.error('Registration error:', error.response?.data || error.message)
        throw error
      }
    },

    async login(email, password) {
      try {
        const response = await http.post('/auth/login', { email, password })
        const { token, data } = response.data

        this.token = token
        this.user = data
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(data))

        await this.executePendingAction()

        const redirect = router.currentRoute.value.query.redirect || '/user'
        router.push(redirect)
      } catch (error) {
        console.error('Login error:', error.response?.data || error.message)
        throw error
      }
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('catalog_filters_applied_knowledge')
      localStorage.removeItem('pendingAction')
      router.push('/user/auth')
    },

    async fetchUser() {
      if (!this.token) return
      try {
        const response = await http.get('/auth/me')
        this.user = response.data
        localStorage.setItem('user', JSON.stringify(response.data))
      } catch (error) {
        if (error.response?.status === 401) {
          this.logout()
        }
        console.error('Fetch user error:', error)
      }
    },

    async executePendingAction() {
      const pendingStr = localStorage.getItem('pendingAction');
      if (!pendingStr) return;

      const pending = JSON.parse(pendingStr);
      if (pending.action === 'favourite') {
        const favouritesStore = useFavouritesStore();
        await favouritesStore.addToFavourites(pending.uniqueId, pending.size);
      }
      localStorage.removeItem('pendingAction');
    }
  }
})