import { defineStore } from 'pinia';
import { api } from '../api';

export const useOrdersStore = defineStore('orders', {
  state: () => ({
    orders: [],
    currentOrder: null,
    loading: false,
  }),
  actions: {
    async fetchOrders() {
      this.loading = true;
      try {
        const response = await api.getUserOrders();
        this.orders = Array.isArray(response.data) ? response.data : [];
      } catch (error) {
        console.error('Failed to fetch orders:', error);
        this.orders = [];
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async fetchOrderDetails(orderId) {
      this.loading = true;
      try {
        const response = await api.getOrderDetails(orderId);
        this.currentOrder = response.data;
        return response.data;
      } catch (error) {
        console.error('Failed to fetch order details:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },
});