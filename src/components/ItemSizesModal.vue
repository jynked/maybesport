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
                                    <span class="cart-quantity" v-if="getQuantityInCart(sizeItem.size)">
                                        {{ getQuantityInCart(sizeItem.size) > 99 ? '99+' :
                                        getQuantityInCart(sizeItem.size) }}
                                    </span>
                                </button>
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
    <QuantityModal :isOpen="isQuantityModalOpen" :initialQuantity="initialQty" :maxQuantity="maxQty"
        @close="closeQuantityModal" @confirm="(qty) => addToCartWithQuantity(selectedSizeForCart, qty)" />
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useFavouritesStore } from '../stores/favourites'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import QuantityModal from './QuantityModal.vue'
import { useToastStore } from '../stores/toast'
import { useI18n } from 'vue-i18n'

const cartStore = useCartStore();
const { t } = useI18n();
const maxQty = ref(Infinity);

const getQuantityInCart = (size) => {
    return cartStore.getItemQuantity(props.uniqueId, size);
};

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

const initialQty = ref(1);

const sortedSizes = computed(() => {
    if (!props.sizes.length) return []
    const getPriority = (sizeItem) => {
        if (sizeItem.quantity > 0 && !sizeItem.isOnRequest) return 1
        if (sizeItem.quantity === 0 && !sizeItem.isOnRequest) return 2
        return 3
    }
    return [...props.sizes].sort((a, b) => getPriority(a) - getPriority(b))
})

const toString = (v) => {
    if (v === undefined || v === null) return '';
    return String(v);
};

function getSizeStatusClass(sizeItem) {
    if (sizeItem.isOnRequest) return 'on-request';
    if (sizeItem.quantity > 0 && !props.isOnRequest) return 'available';
    return 'out-of-stock';
}

function getAvailabilityStatus(sizes) {
    const onRequest = sizes.some(s => s.isOnRequest);
    const available = sizes.some(s => s.quantity > 0 && !s.isOnRequest);
    if (onRequest) return 'on_request';
    if (available) return 'available';
    return 'out_of_stock';
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
            await favouritesStore.removeFromFavourites(props.uniqueId, size);
            useToastStore().success(t('removedFromFavourites'));
        } else {
            await favouritesStore.addToFavourites(props.uniqueId, size);
            useToastStore().success(t('addedToFavourites'));
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

function openQuantityModal(size) {
    const sizeObj = props.sizes.find(s => toString(s.size) === toString(size));
    if (!sizeObj) return;
    selectedSizeForCart.value = size;
    const maxQ = sizeObj.isOnRequest ? Infinity : (sizeObj.quantity > 0 ? sizeObj.quantity : Infinity);
    maxQty.value = maxQ;
    initialQty.value = cartStore.getItemQuantity(props.uniqueId, size) || 1;
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

onMounted(async () => {
  await favouritesStore.loadFavourites();
});

watch(() => props.isOpen, async (newVal) => {
  if (newVal) {
    await favouritesStore.loadFavourites();
  }
});
</script>

<style scoped lang="scss">
@use '../assets/styles/pages/item';
</style>