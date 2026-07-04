import { defineStore } from 'pinia'
import { http } from '../api'
import router from '../router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    loading: false,
  }),
  getters: {
    isAuthenticated: (state) => !!state.user,
    isAdmin: (state) => state.user?.is_admin === true,
  },
  actions: {
    async register(email, password, name = '', captchaToken) {
      try {
        const response = await http.post('/auth/register', { email, password, name, captchaToken })
        this.user = response.data.user
        await this.executePendingAction()
        router.push(router.currentRoute.value.query.redirect || '/user')
      } catch (error) {
        console.error('Registration error:', error)
        throw error
      }
    },
    async login(email, password) {
      try {
        const response = await http.post('/auth/login', { email, password })
        this.user = response.data.user
        await this.executePendingAction()
        router.push(router.currentRoute.value.query.redirect || '/user')
      } catch (error) {
        console.error('Login error:', error)
        throw error
      }
    },
    async logout() {
      this.user = null
      localStorage.removeItem('pendingAction')
      router.push('/user/auth')
    },
    async fetchUser() {
      if (this.user) return
      this.loading = true
      try {
        const response = await http.get('/auth/me')
        this.user = response.data
      } catch (error) {
        this.user = null
        if (error.response?.status !== 401) {
          console.error('Fetch user error:', error)
        }
      } finally {
        this.loading = false
      }
    },
    async executePendingAction() {
      const pendingStr = localStorage.getItem('pendingAction')
      if (!pendingStr) return
      const pending = JSON.parse(pendingStr)
      if (pending.action === 'favourite') {
        const { useFavouritesStore } = await import('./favourites')
        const favouritesStore = useFavouritesStore()
        await favouritesStore.addToFavourites(pending.uniqueId, pending.size)
      } else if (pending.action === 'cart') {
        const { useCartStore } = await import('./cart')
        const cartStore = useCartStore()
        await cartStore.addToCart(pending.uniqueId, pending.size, pending.quantity || 1)
      }
      localStorage.removeItem('pendingAction')
    }
  }
})