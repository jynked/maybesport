<template>
    <main>
        <div class="favourites-heading">
            <h1 v-appear="{ delay: 200 }">{{ $t('favouritesHeaderNav') }}</h1>
            <div class="favourites-action-buttons" ref="actionButtons" v-if="favouritesStore.favouriteItems.length > 0">
                <button @click="toggleViewMode" v-appear="{ delay: 200 }">
                    {{ viewMode === 'category' ? $t('inOrder') : $t('byCategory') }}
                </button>
                <template v-if="viewMode === 'category' && ((grouped.inStock.length > 0 && grouped.outOfStock.length > 0) ||
                    (grouped.onRequest.length > 0 && grouped.outOfStock.length > 0) ||
                    (grouped.onRequest.length > 0 && grouped.inStock.length > 0))">
                    <button v-appear="{ delay: 300 }" v-if="grouped.inStock.length"
                        @click="scrollToSection(inStockSection)">
                        {{ $t('available') }}
                    </button>
                    <button v-appear="{ delay: 400 }" v-if="grouped.outOfStock.length"
                        @click="scrollToSection(outOfStockSection)">
                        {{ $t('out_of_stock') }}
                    </button>
                    <button v-appear="{ delay: 500 }" v-if="grouped.onRequest.length"
                        @click="scrollToSection(onRequestSection)">
                        {{ $t('on_request') }}
                    </button>
                </template>
            </div>
        </div>

        <div v-if="favouritesStore.favouriteItems.length == 0" class="undefined-items-container"
            v-appear="{ delay: 400 }">
            <p>{{ $t('undefinedItems') }}</p>
            <img src="../assets/img/fail.png" :alt="$t('failAlt')">
        </div>

        <div v-else-if="viewMode === 'category'" class="favourites-sections">
            <div class="section in-stock-section" v-if="grouped.inStock.length" ref="inStockSection">
                <div class="section-header">
                    <h2 v-appear="{ delay: 200 }">{{ $t('available') }}</h2>
                </div>
                <div v-show="expanded.inStock" class="section-content">
                    <div class="favourites-grid">
                        <div v-for="(row, rowIndex) in chunkedItems(grouped.inStock)" :key="`inStock-row-${rowIndex}`"
                            class="items-row">
                            <FavouriteItemCard v-for="(item, colIndex) in row" :key="`${item.uniqueId}|${item.size}`"
                                v-memo="[item.uniqueId, item.size, item.availability, item.price]"
                                :uniqueId="item.uniqueId" :title="$i18n.locale === 'en' ? item.title.en : item.title.ru"
                                :image="item.image" :size="item.size" :price="item.price"
                                :isOnRequest="item.isOnRequest" :quantity="item.quantity"
                                :availability="item.availability" @remove="handleRemove"
                                @addToCart="() => openQuantityModal(item)" v-appear="{ delay: 200 + colIndex * 150 }"
                                :delay="100 + colIndex * 200" />
                        </div>
                    </div>
                </div>
            </div>

            <div class="section out-of-stock-section" v-if="grouped.outOfStock.length" ref="outOfStockSection">
                <div class="section-header">
                    <h2 v-appear="{ delay: 200 }">{{ $t('out_of_stock') }}</h2>
                </div>
                <div v-show="expanded.outOfStock" class="section-content">
                    <div class="favourites-grid">
                        <div v-for="(row, rowIndex) in chunkedItems(grouped.outOfStock)"
                            :key="`outOfStock-row-${rowIndex}`" class="items-row">
                            <FavouriteItemCard v-for="(item, colIndex) in row" :key="`${item.uniqueId}|${item.size}`"
                                :uniqueId="item.uniqueId" :title="$i18n.locale === 'en' ? item.title.en : item.title.ru"
                                :image="item.image" :size="item.size" :price="item.price"
                                :isOnRequest="item.isOnRequest" :quantity="item.quantity"
                                :availability="item.availability" @remove="handleRemove"
                                @addToCart="() => openQuantityModal(item)" v-appear="{ delay: 200 + colIndex * 150 }"
                                :delay="600 + colIndex * 200" />
                        </div>
                    </div>
                </div>
            </div>

            <div class="section on-request-section" v-if="grouped.onRequest.length" ref="onRequestSection">
                <div class="section-header">
                    <h2 v-appear="{ delay: 200 }">{{ $t('on_request') }}</h2>
                </div>
                <div v-show="expanded.onRequest" class="section-content">
                    <div class="favourites-grid">
                        <div v-for="(row, rowIndex) in chunkedItems(grouped.onRequest)"
                            :key="`onRequest-row-${rowIndex}`" class="items-row">
                            <FavouriteItemCard v-for="(item, colIndex) in row" :key="`${item.uniqueId}|${item.size}`"
                                :uniqueId="item.uniqueId" :title="$i18n.locale === 'en' ? item.title.en : item.title.ru"
                                :image="item.image" :size="item.size" :price="item.price"
                                :isOnRequest="item.isOnRequest" :quantity="item.quantity"
                                :availability="item.availability" @remove="handleRemove"
                                @addToCart="() => openQuantityModal(item)" v-appear="{ delay: 200 + colIndex * 150 }"
                                :delay="600 + colIndex * 200" />
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-else class="favourites-sections order-mode">
            <div class="section-content">
                <div class="favourites-grid">
                    <div v-for="(row, rowIndex) in chunkedItems(favouritesStore.favouriteItems)"
                        :key="`order-row-${rowIndex}`" class="items-row">
                        <FavouriteItemCard v-for="(item, colIndex) in row" :key="`${item.uniqueId}|${item.size}`"
                            :uniqueId="item.uniqueId" :title="$i18n.locale === 'en' ? item.title.en : item.title.ru"
                            :image="item.image" :size="item.size" :price="item.price" :isOnRequest="item.isOnRequest"
                            :quantity="item.quantity" :availability="item.availability" @remove="handleRemove"
                            @addToCart="() => openQuantityModal(item)" v-appear="{ delay: 200 + colIndex * 150 }"
                            :delay="100 + colIndex * 200" />
                    </div>
                </div>
            </div>
        </div>
        <QuantityModal :isOpen="isQuantityModalOpen" :initialQuantity="1" @close="closeQuantityModal"
            @confirm="addToCartWithQuantity" />
    </main>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { useFavouritesStore } from '../stores/favourites';
import { useCartStore } from '../stores/cart';
import { useAuthStore } from '../stores/auth';
import FavouriteItemCard from '../components/FavouriteItemCard.vue';
import QuantityModal from '../components/QuantityModal.vue';
import { useRouter } from 'vue-router';
import { useToastStore } from '../stores/toast';
import { useI18n } from 'vue-i18n';

const favouritesStore = useFavouritesStore();
const cartStore = useCartStore();
const authStore = useAuthStore();
const emit = defineEmits(['page-loaded']);

const router = useRouter();
const { t } = useI18n();

const actionButtons = ref(null);
let isFixedActive = false;
let ticking = false;

const expanded = ref({
    inStock: true,
    onRequest: true,
    outOfStock: true
});

const viewMode = ref('order');

const inStockSection = ref(null);
const outOfStockSection = ref(null);
const onRequestSection = ref(null);

const toggleViewMode = () => {
    viewMode.value = viewMode.value === 'category' ? 'order' : 'category';
    window.scrollTo(0, 0);
};

const getItemStatus = (item) => {
    if (item.quantity > 0 && !item.isOnRequest) return 'inStock';
    if (item.quantity > 0 && item.isOnRequest) return 'onRequest';
    return 'outOfStock';
};

const grouped = computed(() => {
    const groups = {
        inStock: [],
        onRequest: [],
        outOfStock: []
    };
    favouritesStore.favouriteItems.forEach(item => {
        const status = getItemStatus(item);
        groups[status].push(item);
    });
    return groups;
});

const chunkedItems = (items) => {
    const itemsPerRow = 4;
    const result = [];
    for (let i = 0; i < items.length; i += itemsPerRow) {
        result.push(items.slice(i, i + itemsPerRow));
    }
    return result;
};

const handleRemove = async (uniqueId, size) => {
    await favouritesStore.removeFromFavourites(uniqueId, size);
    useToastStore().success(t('removedFromFavourites'));
};

const handleScroll = () => {
    if (!ticking) {
        requestAnimationFrame(() => {
            const scrollY = window.scrollY;
            const threshold = 150;
            const shouldBeFixed = scrollY > threshold;
            if (shouldBeFixed && !isFixedActive) {
                actionButtons.value?.classList.add('favourites-action-buttons--fixed');
                isFixedActive = true;
            } else if (!shouldBeFixed && isFixedActive) {
                actionButtons.value?.classList.remove('favourites-action-buttons--fixed');
                isFixedActive = false;
            }
            ticking = false;
        });
        ticking = true;
    }
};

const scrollToSection = (sectionRef) => {
    if (sectionRef) {
        sectionRef.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
};

const isQuantityModalOpen = ref(false);
const pendingItem = ref(null);

function openQuantityModal(item) {
    pendingItem.value = {
        uniqueId: item.uniqueId,
        size: item.size
    };
    isQuantityModalOpen.value = true;
}

function closeQuantityModal() {
    isQuantityModalOpen.value = false;
    pendingItem.value = null;
}

async function addToCartWithQuantity(quantity) {
    if (!pendingItem.value) return;

    if (!authStore.isAuthenticated) {
        localStorage.setItem('pendingAction', JSON.stringify({
            action: 'cart',
            uniqueId: pendingItem.value.uniqueId,
            size: pendingItem.value.size,
            quantity: quantity
        }));
        router.push({ name: 'UserAuth', query: { redirect: router.currentRoute.value.fullPath } });
        closeQuantityModal();
        return;
    }

    try {
        await cartStore.addToCart(pendingItem.value.uniqueId, pendingItem.value.size, quantity);
        closeQuantityModal();
    } catch (error) {
        console.error('Ошибка добавления в корзину', error);
    }
}

onMounted(async () => {
    if (authStore.isAuthenticated) {
        await favouritesStore.loadFavourites();
    } else {
        favouritesStore.favouriteItems = [];
        favouritesStore.favouriteKeys = [];
    }
    emit('page-loaded', true);
    window.addEventListener('scroll', handleScroll);
});

onBeforeUnmount(() => {
    window.removeEventListener('scroll', handleScroll);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/favourites';
</style>