import axios from 'axios'

const API_BASE = 'http://localhost:8080/api'

export const api = {
  // Существующие методы
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

  createItem(itemData) {
    return axios.post(`${API_BASE}/admin/items`, itemData)
  },

  updateItem(id, itemData) {
    return axios.patch(`${API_BASE}/admin/items/${id}`, itemData)
  },

  deleteItem(id) {
    return axios.delete(`${API_BASE}/admin/items/${id}`)
  }
}