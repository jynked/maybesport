import axios from 'axios'
import { useAuthStore } from './stores/auth'

const API_BASE = 'http://localhost:5173/api'

const instance = axios.create({
  baseURL: API_BASE,
})

instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
    }
    return Promise.reject(error)
  }
)

export const http = instance

export const api = {
  getExchangeRate() {
    return http.get('/exchange-rate');
  },
  getItems(params) {
    return http.get('/items', { params })
  },
  getItem(uniqueId) {
    return http.get(`/items/${uniqueId}`)
  },
  getSimilar(uniqueId, limit = 4) {
    return http.get(`/items/${uniqueId}/similar`, { params: { limit } })
  },
  getMainPageNew() {
    return http.get('/main-page/new')
  },
  getAdminItems() {
    return http.get('/admin/items')
  },
  getAdminItem(id) {
    return http.get(`/admin/items/${id}`)
  },
  createMainItem(data) {
    return http.post('/admin/items', data)
  },
  updateMainItem(id, data) {
    return http.put(`/admin/items/${id}`, data)
  },
  deleteMainItem(id) {
    return http.delete(`/admin/items/${id}`)
  },
  createItem(itemData) {
    console.warn('createItem устарел, используйте createMainItem')
    return http.post('/admin/items', itemData)
  },
  updateItem(id, itemData) {
    console.warn('updateItem устарел, используйте updateMainItem')
    return http.patch(`/admin/items/${id}`, itemData)
  },
  deleteItem(id) {
    console.warn('deleteItem устарел, используйте deleteMainItem')
    return http.delete(`/admin/items/${id}`)
  },
  getFilters() {
    return http.get('/filters')
  },

  addToFavourites(uniqueId, size) {
    return http.post(`/user/favourites/${uniqueId}`, { size });
  },
  removeFromFavourites(uniqueId, size) {
    return http.delete(`/user/favourites/${uniqueId}?size=${encodeURIComponent(size)}`);
  },
  getFavouriteItems() {
    return http.get('/user/favourites/items');
  },

  getUserOrders() {
    return http.get('/user/orders');
  },
  getOrderDetails(orderId) {
    return http.get(`/user/orders/${orderId}`);
  },

  getCart() {
    return http.get('/user/cart');
  },
  addToCart(uniqueId, size, quantity) {
    return http.post(`/user/cart/${uniqueId}`, { size, quantity });
  },
  updateCartItem(uniqueId, size, quantity) {
    return http.put(`/user/cart/${uniqueId}`, { size, quantity });
  },
  removeFromCart(uniqueId, size) {
    return http.delete(`/user/cart/${uniqueId}?size=${encodeURIComponent(size)}`);
  },

  updateProfile(data) {
    return http.put('/user/profile', data);
  },
  deleteAccount() {
    return http.delete('/user/profile');
  },

  createOrder(items, deliveryAddress, comment) {
    return http.post('/user/orders', { items, deliveryAddress, comment });
  },

  getAdminOrders(params) {
    return http.get('/admin/orders', { params });
  },
  getAdminOrderDetails(orderId) {
    return http.get(`/admin/orders/${orderId}`);
  },
  updateOrderStatus(orderId, status, description) {
    return http.put(`/admin/orders/${orderId}/status`, { status, description });
  },
  deleteLastOrderStatus(orderId) {
    return http.delete(`/admin/orders/${orderId}/status/last`);
  },
}