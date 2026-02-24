import axios from 'axios'

const API_BASE = 'http://localhost:8080/api'

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
  getItems(params) {
    return axios.get(`${API_BASE}/items`, { params })
  },
  getItem(uniqueId) {
    return axios.get(`${API_BASE}/items/${uniqueId}`)
  },
  getSimilar(uniqueId, limit = 4) {
    return axios.get(`${API_BASE}/items/${uniqueId}/similar`, { params: { limit } })
  },
  getMainPageNew() {
    return axios.get(`${API_BASE}/main-page/new`)
  },
  getAdminItems() {
    return axios.get(`${API_BASE}/admin/items`)
  },
  getAdminItem(id) {
    return axios.get(`${API_BASE}/admin/items/${id}`)
  },
  createMainItem(data) {
    return axios.post(`${API_BASE}/admin/items`, data)
  },
  updateMainItem(id, data) {
    return axios.put(`${API_BASE}/admin/items/${id}`, data)
  },
  deleteMainItem(id) {
    return axios.delete(`${API_BASE}/admin/items/${id}`)
  },
  createItem(itemData) {
    console.warn('createItem устарел, используйте createMainItem')
    return axios.post(`${API_BASE}/admin/items`, itemData)
  },
  updateItem(id, itemData) {
    console.warn('updateItem устарел, используйте updateMainItem')
    return axios.patch(`${API_BASE}/admin/items/${id}`, itemData)
  },
  deleteItem(id) {
    console.warn('deleteItem устарел, используйте deleteMainItem')
    return axios.delete(`${API_BASE}/admin/items/${id}`)
  }
}