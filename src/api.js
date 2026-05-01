import axios from 'axios'
import { useAuthStore } from './stores/auth'
import { useToastStore } from './stores/toast'
import i18n from './i18n';

const API_BASE = import.meta.env.VITE_API_URL || '/api';

const instance = axios.create({
  baseURL: API_BASE,
  withCredentials: true,
})

instance.interceptors.request.use((config) => {
  if (config.method !== 'get') {
    const csrfToken = document.cookie.split('; ').find(row => row.startsWith('csrf_token='))?.split('=')[1];
    if (csrfToken) {
      config.headers['X-CSRF-Token'] = csrfToken;
    }
  }
  return config;
});

const $t = (key) => i18n.global.t(key)

instance.interceptors.response.use(
  (response) => response,
  (error) => {
    const toastStore = useToastStore();
    const status = error.response?.status;
    const data = error.response?.data;

    if (status === 401) {
      const authStore = useAuthStore();
      authStore.token = null;
      authStore.user = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      toastStore.error($t('sessionError'));
    }
    else if (status === 403) {
      toastStore.error($t('noIssues'));
    }
    else if (status === 500) {
      toastStore.error($t('serverError'));
    }
    else if (!error.response) {
      toastStore.error($t('connectError'));
    }
    else {
      const message = data?.error || data || $t('errorLoader');
      if (typeof message === 'string') toastStore.error(message);
      else toastStore.error($t('errorLoader'));
    }
    return Promise.reject(error);
  }
);

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

  createOrder(items, deliveryAddress, comment, captchaToken) {
    return http.post('/user/orders', { items, deliveryAddress, comment, captchaToken });
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

  getAdminOrderItems(orderId) {
    return http.get(`/admin/orders/${orderId}/items`);
  },
  updateOrderItemStatus(itemId, status, description) {
    return http.put(`/admin/order-items/${itemId}/status`, { status, description });
  },
}