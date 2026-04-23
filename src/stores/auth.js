import { defineStore } from 'pinia'
import { http, api } from '../api'
import router from '../router'
import { useFavouritesStore } from './favourites'
import { useCartStore } from './cart'
import { i18n } from '../i18n';

const { t } = i18n.global;

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    user: (() => {
      const userStr = localStorage.getItem('user');
      if (userStr) {
        try {
          const parsed = JSON.parse(userStr);
          if (parsed.is_admin === undefined) parsed.is_admin = false;
          return parsed;
        } catch (e) {
          console.error('Failed to parse user from localStorage', e);
          return null;
        }
      }
      return null;
    })(),
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.is_admin === true,
  },

  actions: {
    async register(email, password, name = '') {
      try {
        const response = await http.post('/auth/register', { email, password, name })
        const { token, user } = response.data
        this.token = token
        this.user = user
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
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
        const { token, user } = response.data
        this.token = token
        this.user = user
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
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
          this.token = null
          this.user = null
          localStorage.removeItem('token')
          localStorage.removeItem('user')
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
      } else if (pending.action === 'cart') {
        await api.addToCart(pending.uniqueId, pending.size, pending.quantity || 1);
      }
      localStorage.removeItem('pendingAction');
    }
  }
})