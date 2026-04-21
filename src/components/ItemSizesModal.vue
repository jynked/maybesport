<template>
    <Transition name="modal">
        <div v-if="isOpen" class="sizes-modal" @click="emit('close')">
            <div class="modal-content-sizes" @click.stop>
                <div class="modal-header">
                    <h2>{{ $t('sizesAndPrices').toUpperCase() }}</h2>
                    <button class="close-button" @click="emit('close')">×</button>
                </div>
                <div class="sizes-content">
                    <div class="size-items">
                        <div v-for="sizeItem in sortedSizes" :key="sizeItem.size"
                            :class="['size-item', getSizeStatusClass(sizeItem), { selected: selectedSize === sizeItem.size }]"
                            @click="selectSize(sizeItem.size)">
                            <span class="size">{{ sizeItem.size }}</span>
                            <span class="price">{{ sizeItem.price.toLocaleString() }} ₽</span>
                            <span class="status">{{ $t(getAvailabilityStatus([sizeItem])) }}</span>
                            <span class="quantity" v-if="sizeItem.quantity > 0 && !sizeItem.isOnRequest">
                                ({{ sizeItem.quantity }} {{ $t('pieces') }})
                            </span>
                            <span class="quantity">&nbsp;</span>
                            <span class="item-actions">
                                <button @click.stop="toggleFavourite(sizeItem.size)"
                                    :disabled="togglingSize === sizeItem.size">
                                    <img src="../assets/img/favourite.png" :alt="$t('favouriteAlt')"
                                        :style="{ filter: isFavouriteForSize(sizeItem.size) ? '' : 'sepia(1)' }" />
                                </button>
                                <button class="item-cart" @click.stop="openQuantityModal(sizeItem.size)">
                                    <img src="../assets/img/cart.png" :alt="$t('cartAlt')" />
                                </button>
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
    <QuantityModal :isOpen="isQuantityModalOpen" :initialQuantity="1" @close="closeQuantityModal"
        @confirm="(qty) => addToCartWithQuantity(selectedSizeForCart, qty)" />
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useFavouritesStore } from '../stores/favourites'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import QuantityModal from './QuantityModal.vue'

const props = defineProps({
    uniqueId: { type: String, required: true },
    sizes: { type: Array, required: true },
    isOpen: { type: Boolean, default: false },
})

const emit = defineEmits(['close', 'addToCart'])

const favouritesStore = useFavouritesStore()
const router = useRouter()
const selectedSize = ref(null)
const togglingSize = ref(null)

const isQuantityModalOpen = ref(false);
const selectedSizeForCart = ref(null);

const sortedSizes = computed(() => {
    if (!props.sizes.length) return []
    const getPriority = (sizeItem) => {
        if (sizeItem.quantity > 0 && !sizeItem.isOnRequest) return 1
        if (sizeItem.quantity === 0 && !sizeItem.isOnRequest) return 2
        return 3
    }
    return [...props.sizes].sort((a, b) => getPriority(a) - getPriority(b))
})

function getSizeStatusClass(sizeItem) {
    if (sizeItem.quantity > 0 && !sizeItem.isOnRequest) return 'available'
    if (sizeItem.quantity > 0 && sizeItem.isOnRequest) return 'on-request'
    return 'out-of-stock'
}

function getAvailabilityStatus(sizes) {
    const available = sizes.some(s => s.quantity > 0 && !s.isOnRequest)
    const onRequest = sizes.some(s => s.quantity > 0 && s.isOnRequest)
    if (available) return 'available'
    if (onRequest) return 'on_request'
    return 'out_of_stock'
}

function selectSize(size) {
    selectedSize.value = size
}

function isFavouriteForSize(size) {
    return favouritesStore.isFavourite(props.uniqueId, size)
}

function redirectToAuthWithAction(action, uniqueId, size) {
    const pending = { action, uniqueId, size }
    localStorage.setItem('pendingAction', JSON.stringify(pending))
    router.push({ name: 'UserAuth', query: { redirect: router.currentRoute.value.fullPath } })
}

async function toggleFavourite(size) {
    if (!props.uniqueId || !size) return
    if (togglingSize.value === size) return

    const authStore = useAuthStore()
    if (!authStore.isAuthenticated) {
        redirectToAuthWithAction('favourite', props.uniqueId, size)
        return
    }

    togglingSize.value = size
    try {
        if (favouritesStore.isFavourite(props.uniqueId, size)) {
            await favouritesStore.removeFromFavourites(props.uniqueId, size)
        } else {
            await favouritesStore.addToFavourites(props.uniqueId, size)
        }
    } catch (error) {
        console.error('Error toggling favourite:', error)
        if (error.response?.status === 401) {
            redirectToAuthWithAction('favourite', props.uniqueId, size)
        }
    } finally {
        togglingSize.value = null
    }
}

function addToCart(size) {
    if (!props.uniqueId || !size) return
    const authStore = useAuthStore()
    if (!authStore.isAuthenticated) {
        redirectToAuthWithAction('cart', props.uniqueId, size)
        return
    }
    console.log('Добавление в корзину:', props.uniqueId, size)
    emit('addToCart', { uniqueId: props.uniqueId, size })
    emit('close')
}

function openQuantityModal(size) {
  selectedSizeForCart.value = size;
  isQuantityModalOpen.value = true;
}

function closeQuantityModal() {
  isQuantityModalOpen.value = false;
  selectedSizeForCart.value = null;
}

async function addToCartWithQuantity(size, quantity) {
  const authStore = useAuthStore();
  if (!authStore.isAuthenticated) {
    redirectToAuthWithAction('cart', props.uniqueId, size);
    return;
  }
  const cartStore = useCartStore();
  try {
    await cartStore.addToCart(props.uniqueId, size, quantity);
    emit('close');
  } catch (error) {
    console.error('Ошибка добавления в корзину', error);
  }
}
</script>

<style scoped lang="scss">
@use '../assets/styles/pages/item';
</style>