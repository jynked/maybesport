<template>
  <Transition name="modal">
    <div v-if="isOpen" class="cart-modal-overlay" @click.self="closeModal">
      <div class="cart-modal">
        <div class="cart-header">
          <h2>{{ $t('cart') }}</h2>
          <button class="close-btn" @click="closeModal">×</button>
        </div>

        <div v-if="loading" class="cart-loading">
          <Loader />
        </div>

        <div v-else-if="cartItems.length === 0" class="cart-empty">
          <p>{{ $t('cartEmpty') }}</p>
        </div>

        <div v-else class="cart-content">
          <div class="cart-select-all">
            <label class="checkbox-label">
              <input type="checkbox" v-model="selectAll" />
              <span>{{ $t('selectAll') }}</span>
            </label>
          </div>

          <div v-if="availableItems.length">
            <h4 class="section-title">Доступные для заказа</h4>
            <div v-for="item in availableItems" :key="`${item.uniqueId}_${item.size}`" class="cart-item">
              <label class="item-checkbox">
                <input type="checkbox" :value="`${item.uniqueId}_${item.size}`" v-model="selectedKeys" />
              </label>
              <div class="item-image">
                <img :src="item.image" :alt="item.title.ru" />
              </div>
              <div class="item-details">
                <div class="item-title">{{ $i18n.locale === 'en' ? item.title.en : item.title.ru }}</div>
                <div class="item-size">{{ $t('size') }}: {{ item.size }}</div>
                <div class="item-price">{{ item.price.toLocaleString() }} ₽</div>
              </div>
              <div class="item-quantity">
                <div class="quantity-control">
                  <button @click="decrementQuantity(item)" :disabled="item.quantity == 1">-</button>
                  <input type="text" :value="localQuantities[`${item.uniqueId}_${item.size}`] ?? item.quantity"
                    @input="handleQuantityInput(item, $event)" @blur="validateQuantity(item)" />
                  <button @click="incrementQuantity(item)">+</button>
                </div>
              </div>
              <div class="item-actions">
                <div class="dropdown">
                  <button class="dropdown-trigger" @click="toggleDropdown(item.uniqueId, item.size, $event)">
                    ⋯
                  </button>
                  <div v-if="activeDropdown === `${item.uniqueId}_${item.size}`" class="dropdown-menu-custom"
                    @click.stop>
                    <button @click="copyLink(item.uniqueId)">{{ $t('copyLink') }}</button>
                    <button @click="shareItem(item)">{{ $t('share') }}</button>
                    <button v-if="!isItemFavourite(item)" @click="addToFavourites(item)">
                      {{ $t('addToFavourites') }}
                    </button>
                    <button v-else @click="removeFromFavourites(item)">
                      {{ $t('removeFromFavourites') }}
                    </button>
                  </div>
                </div>
                <button class="remove-btn" @click="removeItem(item)">
                  🗑️
                </button>
              </div>
            </div>
          </div>

          <div v-if="unavailableItems.length">
            <h4 class="section-title unavailable">Недоступные (нет в наличии)</h4>
            <div v-for="item in unavailableItems" :key="`${item.uniqueId}_${item.size}`" class="cart-item">
              <div class="item-checkbox disabled"></div>
              <div class="item-image">
                <img :src="item.image" :alt="item.title.ru" />
              </div>
              <div class="item-details">
                <div class="item-title">{{ $i18n.locale === 'en' ? item.title.en : item.title.ru }}</div>
                <div class="item-size">{{ $t('size') }}: {{ item.size }}</div>
                <div class="item-price">{{ item.price.toLocaleString() }} ₽</div>
              </div>
              <div class="item-quantity">
                <div class="quantity-control disabled">
                  <button disabled>-</button>
                  <input type="text" :value="localQuantities[`${item.uniqueId}_${item.size}`] ?? item.quantity"
                    disabled />
                  <button disabled>+</button>
                </div>
              </div>
              <div class="item-actions">
                <div class="dropdown">
                  <button class="dropdown-trigger" @click="toggleDropdown(item.uniqueId, item.size, $event)">⋯</button>
                  <div v-if="activeDropdown === `${item.uniqueId}_${item.size}`" class="dropdown-menu-custom"
                    @click.stop>
                    <button @click="copyLink(item.uniqueId)">{{ $t('copyLink') }}</button>
                    <button @click="shareItem(item)">{{ $t('share') }}</button>
                    <button v-if="!isItemFavourite(item)" @click="addToFavourites(item)">
                      {{ $t('addToFavourites') }}
                    </button>
                    <button v-else @click="removeFromFavourites(item)">
                      {{ $t('removeFromFavourites') }}
                    </button>
                  </div>
                </div>
                <button class="remove-btn" @click="removeItem(item)">🗑️</button>
              </div>
            </div>
          </div>
          <div class="cart-footer">
            <div class="total-price">
              {{ $t('totalSelected') }}: <strong>{{ totalSelectedPrice.toLocaleString() }} ₽</strong>
            </div>
            <button class="checkout-btn" @click="openCheckoutModal" :disabled="selectedItems.length === 0">
              {{ $t('checkout') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>

  <Transition name="modal">
    <div v-if="isCheckoutModalOpen" class="checkout-modal-overlay" @click.self="closeCheckoutModal">
      <div class="checkout-modal" @click.stop>
        <div class="checkout-header">
          <h2>{{ $t('checkoutOrder') }}</h2>
          <button class="close-btn" @click="closeCheckoutModal">×</button>
        </div>
        <div class="checkout-body">
          <div class="form-group">
            <label>{{ $t('deliveryAddress') }}</label>
            <input type="text" v-model="deliveryAddress" :placeholder="$t('enterAddress')" />
          </div>
          <div class="form-group">
            <label>{{ $t('comment') }} ({{ $t('optional') }})</label>
            <textarea v-model="comment" rows="3" :placeholder="$t('commentPlaceholder')"></textarea>
          </div>
          <div class="checkout-summary">
            <span>{{ $t('totalAmount') }}: </span>
            <strong>{{ totalSelectedPrice.toLocaleString() }} ₽</strong>
          </div>
          <button class="submit-order" @click="submitOrder" :disabled="orderLoading">
            <span v-if="orderLoading">{{ $t('processing') }}...</span>
            <span v-else>{{ $t('placeOrder') }}</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { useCartStore } from '../stores/cart';
import { useFavouritesStore } from '../stores/favourites';
import { useAuthStore } from '../stores/auth';
import { api } from '../api';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import Loader from './LoaderVue.vue';
import { debounce } from 'lodash';

const props = defineProps({
  isOpen: Boolean,
});

const emit = defineEmits(['close']);

const cartStore = useCartStore();
const favouritesStore = useFavouritesStore();
const authStore = useAuthStore();
const router = useRouter();
const { t } = useI18n();

const loading = ref(false);
const cartItems = ref([]);
const selectedKeys = ref([]);
const activeDropdown = ref(null);
const isCheckoutModalOpen = ref(false);
const deliveryAddress = ref('');
const comment = ref('');
const orderLoading = ref(false);

const localQuantities = ref({});

async function loadCart() {
  loading.value = true;
  try {
    await cartStore.loadCart();
    cartItems.value = cartStore.cartItems?.map(item => ({ ...item })) || [];
    selectedKeys.value = availableItems.value.map(item => `${item.uniqueId}_${item.size}`);
    initLocalQuantities();
  } catch (err) {
    console.error('Failed to load cart', err);
  } finally {
    loading.value = false;
  }
}

const totalSelectedPrice = computed(() => {
  return selectedItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0);
});

const availableItems = computed(() => {
  return cartItems.value.filter(item => item.isOnRequest || item.stock > 0);
});

const unavailableItems = computed(() => {
  return cartItems.value.filter(item => !item.isOnRequest && item.stock <= 0);
});

const selectedItems = computed(() => {
  return availableItems.value.filter(item => {
    const key = `${item.uniqueId}_${item.size}`;
    return selectedKeys.value.includes(key);
  });
});

const selectAll = computed({
  get() {
    return availableItems.value.length > 0 && selectedItems.value.length === availableItems.value.length;
  },
  set(val) {
    if (val) {
      selectedKeys.value = availableItems.value.map(item => `${item.uniqueId}_${item.size}`);
    } else {
      selectedKeys.value = [];
    }
  }
});

function isItemFavourite(item) {
  return favouritesStore.isFavourite(item.uniqueId, item.size);
}

async function updateQuantity(item, newQuantity) {
  if (newQuantity < 1) {
    await cartStore.removeFromCart(item.uniqueId, item.size);
  } else {
    await cartStore.updateQuantity(item.uniqueId, item.size, newQuantity);
  }
  await loadCart();
}

async function removeItem(item) {
  const key = `${item.uniqueId}_${item.size}`;
  const oldCartItems = [...cartItems.value];
  const oldSelectedKeys = [...selectedKeys.value];
  const oldLocalQuantities = { ...localQuantities.value };

  cartItems.value = cartItems.value.filter(i => !(i.uniqueId === item.uniqueId && String(i.size) === String(item.size)));
  selectedKeys.value = selectedKeys.value.filter(k => k !== key);
  delete localQuantities.value[key];

  try {
    await cartStore.removeFromCart(item.uniqueId, item.size);
  } catch (error) {
    cartItems.value = oldCartItems;
    selectedKeys.value = oldSelectedKeys;
    localQuantities.value = oldLocalQuantities;
    console.error('Failed to remove item', error);
    alert(t('removeFailed'));
  }
}

function toggleDropdown(uniqueId, size, event) {
  event.stopPropagation();
  const key = `${uniqueId}_${size}`;
  if (activeDropdown.value === key) {
    activeDropdown.value = null;
  } else {
    activeDropdown.value = key;
  }
}

async function copyLink(uniqueId) {
  const url = `${window.location.origin}/catalog/${uniqueId}`;
  await navigator.clipboard.writeText(url);
  alert(t('linkCopied'));
  activeDropdown.value = null;
}

async function shareItem(item) {
  const shareData = {
    title: item.title.ru,
    text: `${item.title.ru} - ${item.price} ₽`,
    url: `${window.location.origin}/catalog/${item.uniqueId}`,
  };
  if (navigator.share) {
    await navigator.share(shareData);
  } else {
    await copyLink(item.uniqueId);
  }
  activeDropdown.value = null;
}

async function addToFavourites(item) {
  if (!authStore.isAuthenticated) {
    localStorage.setItem('pendingAction', JSON.stringify({
      action: 'favourite',
      uniqueId: item.uniqueId,
      size: item.size
    }));
    router.push({ name: 'UserAuth', query: { redirect: '/user/favourites' } });
    closeModal();
    return;
  }
  await favouritesStore.addToFavourites(item.uniqueId, item.size);
  activeDropdown.value = null;
}

async function removeFromFavourites(item) {
  if (!authStore.isAuthenticated) {
    router.push({ name: 'UserAuth', query: { redirect: '/user/favourites' } });
    return;
  }
  await favouritesStore.removeFromFavourites(item.uniqueId, item.size);
  activeDropdown.value = null;
}

function openCheckoutModal() {
  if (selectedItems.value.length === 0) return;
  isCheckoutModalOpen.value = true;
}

function closeCheckoutModal() {
  isCheckoutModalOpen.value = false;
  deliveryAddress.value = '';
  comment.value = '';
}

async function submitOrder() {
  if (!deliveryAddress.value.trim()) {
    alert(t('addressRequired'));
    return;
  }
  if (orderLoading.value) return;
  orderLoading.value = true;
  try {
    const orderItems = selectedItems.value.map(item => ({
      uniqueId: item.uniqueId,
      size: item.size,
      quantity: item.quantity
    }));
    await api.createOrder(orderItems, deliveryAddress.value, comment.value);
    alert(t('orderSuccess'));
    await loadCart();
    closeCheckoutModal();
    closeModal();
    router.push({ name: 'UserOrders' });
  } catch (err) {
    console.error('Order creation failed', err);
    if (err.response?.status === 409) {
      alert(t('errorCart'));
      await loadCart();
    } else {
      alert(t('orderFailed'));
    }
  } finally {
    orderLoading.value = false;
  }
}

function closeModal() {
  emit('close');
  document.body.style.overflowY = 'auto'
  activeDropdown.value = null;
  isCheckoutModalOpen.value = false;
}

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    loadCart();
  } else {
    activeDropdown.value = null;
  }
});

function handleClickOutside(event) {
  if (activeDropdown.value && !event.target.closest('.dropdown')) {
    activeDropdown.value = null;
  }
}

function initLocalQuantities() {
  const newQuantities = {};
  cartItems.value.forEach(item => {
    const key = `${item.uniqueId}_${item.size}`;
    newQuantities[key] = item.quantity;
  });
  localQuantities.value = newQuantities;
}

const debouncedUpdateQuantity = debounce(async (uniqueId, size, newQuantity) => {
  try {
    await cartStore.updateQuantity(uniqueId, size, newQuantity);
  } catch (error) {
    const key = `${uniqueId}_${size}`;
    const originalItem = cartItems.value.find(i => i.uniqueId === uniqueId && String(i.size) === String(size));
    if (originalItem) {
      localQuantities.value[key] = originalItem.quantity;
      originalItem.quantity = originalItem.quantity;
    }
    console.error('Failed to update quantity', error);
  }
}, 500);

function updateQuantityOptimistic(item, newQuantity) {
  if (item.stock <= 0 && !item.isOnRequest) {
    return;
  }
  const maxQ = item.isOnRequest ? Infinity : item.stock;
  if (newQuantity > maxQ) {
    const key = `${item.uniqueId}_${item.size}`;
    localQuantities.value[key] = item.quantity;
    return;
  }
  const key = `${item.uniqueId}_${item.size}`;
  if (newQuantity < 1) {
    removeItem(item);
    return;
  }
  localQuantities.value[key] = newQuantity;
  item.quantity = newQuantity;
  debouncedUpdateQuantity(item.uniqueId, item.size, newQuantity);
}

function incrementQuantity(item) {
  const key = `${item.uniqueId}_${item.size}`;
  const current = localQuantities.value[key] ?? item.quantity;
  updateQuantityOptimistic(item, current + 1);
}

function decrementQuantity(item) {
  const key = `${item.uniqueId}_${item.size}`;
  const current = localQuantities.value[key] ?? item.quantity;
  if (current <= 1) {
    removeItem(item);
  } else {
    updateQuantityOptimistic(item, current - 1);
  }
}

function handleQuantityInput(item, event) {
  let val = event.target.value.replace(/[^0-9]/g, '');
  if (val === '' || val === '0') val = '1';
  val = val.replace(/^0+/, '');
  let num = parseInt(val, 10);
  if (isNaN(num) || num < 1) num = 1;
  updateQuantityOptimistic(item, num);
}

function validateQuantity(item) {
  const key = `${item.uniqueId}_${item.size}`;
  let val = localQuantities.value[key];
  if (val < 1) {
    removeItem(item);
    return;
  }
  const maxQ = item.isOnRequest ? Infinity : item.stock;
  if (val > maxQ) {
    updateQuantityOptimistic(item, maxQ);
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/cart-modal';
</style>