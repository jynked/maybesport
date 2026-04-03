import axios from 'axios'

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
      localStorage.removeItem('token')
      localStorage.removeItem('user')
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

  // updateProfile(data) {
  //   return http.put('/user/profile', data);
  // },
  // deleteAccount() {
  //   return http.delete('/user/profile');
  // }
}