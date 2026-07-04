import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import { api } from '../api';
import { useAuthStore } from './auth';
import { useToastStore } from './toast';
import { i18n } from '../i18n';
import router from '../router';

const { t } = i18n.global;

export const useCartStore = defineStore('cart', () => {
  const cartItems = ref([]);
  const loading = ref(false);
  let queue = Promise.resolve();

  function enqueue(fn) {
    const result = queue.then(() => fn()).catch(err => {
      console.warn('Cart action failed:', err);
    });
    queue = result;
    return result;
  }

  async function loadCart() {
    const authStore = useAuthStore();
    if (!authStore.isAuthenticated) {
      cartItems.value = [];
      return;
    }
    loading.value = true;
    try {
      const response = await api.getCart();
      cartItems.value = response.data;
    } catch (error) {
      console.error('Failed to load cart', error);
    } finally {
      loading.value = false;
    }
  }

  async function addToCart(uniqueId, size, quantity) {
    const authStore = useAuthStore();
    if (!authStore.isAuthenticated) {
      localStorage.setItem('pendingAction', JSON.stringify({
        action: 'cart',
        uniqueId,
        size,
        quantity: quantity || 1
      }));
      router.push({ name: 'UserAuth', query: { redirect: router.currentRoute.value.fullPath } });
      throw new Error('Not authenticated');
    }
    return enqueue(async () => {
      if (!cartItems.value) cartItems.value = [];
      const existingIndex = cartItems.value.findIndex(
        item => item && item.uniqueId === uniqueId && String(item.size) === String(size)
      );
      let oldItems = [...cartItems.value];
      if (existingIndex !== -1) {
        cartItems.value[existingIndex].quantity = quantity;
      } else {
        cartItems.value.push({
          uniqueId,
          size,
          quantity,
        });
      }

      try {
        await api.addToCart(uniqueId, size, quantity);
        useToastStore().success(t('addedToCart'));
        await loadCart();
      } catch (error) {
        cartItems.value = oldItems;
        console.error('Add to cart failed', error);
        throw error;
      }
    });
  }

  async function updateQuantity(uniqueId, size, quantity) {
    return enqueue(async () => {
      const index = cartItems.value.findIndex(
        item => item.uniqueId === uniqueId && String(item.size) === String(size)
      );
      if (index === -1) return;
      const oldQuantity = cartItems.value[index].quantity;
      cartItems.value[index].quantity = quantity;

      try {
        await api.updateCartItem(uniqueId, size, quantity);
        if (quantity <= 0) {
          cartItems.value = cartItems.value.filter(
            item => !(item.uniqueId === uniqueId && String(item.size) === String(size))
          );
        }
        useToastStore().success(t('addedToCart'));
      } catch (error) {
        cartItems.value[index].quantity = oldQuantity;
        console.error('Update cart item failed', error);
        throw error;
      }
    });
  }

  async function removeFromCart(uniqueId, size) {
    return enqueue(async () => {
      const oldItems = [...cartItems.value];
      cartItems.value = cartItems.value.filter(
        item => !(item.uniqueId === uniqueId && String(item.size) === String(size))
      );
      try {
        await api.removeFromCart(uniqueId, size);
        useToastStore().success(t('removedFromCart'));
      } catch (error) {
        cartItems.value = oldItems;
        console.error('Remove from cart failed', error);
        throw error;
      }
    });
  }

  function getItemQuantity(uniqueId, size) {
    const item = cartItems.value?.find(
      i => i.uniqueId === uniqueId && String(i.size) === String(size)
    );
    return item ? item.quantity : 0;
  }

  const totalQuantity = computed(() => {
    return cartItems.value?.reduce((sum, item) => sum + (item?.quantity || 0), 0) > 99 ? '99+' : cartItems.value?.reduce((sum, item) => sum + (item?.quantity || 0), 0);
  });

  function getTotalQuantityByUniqueId(uniqueId) {
    return cartItems.value?.reduce((sum, item) => {
      if (item.uniqueId === uniqueId) return sum + item.quantity;
      return sum;
    }, 0) || 0;
  }

  return {
    cartItems,
    loading,
    loadCart,
    addToCart,
    updateQuantity,
    removeFromCart,
    getItemQuantity,
    totalQuantity,
    getTotalQuantityByUniqueId,
  };
});