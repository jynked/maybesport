<template>
    <main>
        <h1 v-appear="{ delay: 200 }">{{ $t('favouritesHeaderNav') }}</h1>
        <div v-if="favouritesStore.loading" class="loader">Loading...</div>
        <div v-else-if="favouritesStore.favouriteItems.length === 0" class="undefined-items-container"
            v-appear="{ delay: 400 }">
            <p>{{ $t('undefinedItems') }}</p>
            <img src="../assets/img/fail.png" :alt="$t('failAlt')">
        </div>

        <div v-else class="favourites-sections">
            <div class="section in-stock-section" v-if="grouped.inStock.length">
                <button class="section-header" @click="toggleSection('inStock')">
                    <h2>{{ $t('available') }} ({{ grouped.inStock.length }})</h2>
                    <img src="../assets/img/down.png" :alt="$t('Down')" :class="{ 'rotated': expanded.inStock }">
                </button>
                <transition name="collapse">
                    <div v-show="expanded.inStock" class="section-content">
                        <div class="favourites-grid">
                            <div v-for="(row, rowIndex) in chunkedItems(grouped.inStock)"
                                :key="`inStock-row-${rowIndex}`" class="items-row">
                                <FavouriteItemCard v-for="(item, colIndex) in row"
                                    :key="`${item.uniqueId}|${item.size}`" :uniqueId="item.uniqueId"
                                    :title="$i18n.locale === 'en' ? item.title.en : item.title.ru" :image="item.image"
                                    :size="item.size" :price="item.price" :isOnRequest="item.isOnRequest"
                                    :quantity="item.quantity" :availability="item.availability" @remove="handleRemove"
                                    v-appear.repeat="{ delay: 200 + colIndex * 150 }" 
                                    :delay="100 + colIndex * 200"
                                    class="in-stock-favourite-item"/>
                            </div>
                        </div>
                    </div>
                </transition>
            </div>

            <div class="section on-request-section" v-if="grouped.onRequest.length">
                <button class="section-header" @click="toggleSection('onRequest')">
                    <h2>{{ $t('on_request') }} ({{ grouped.onRequest.length }})</h2>
                    <img src="../assets/img/down.png" :alt="$t('Down')" :class="{ 'rotated': expanded.onRequest }">
                </button>
                <transition name="collapse">
                    <div v-show="expanded.onRequest" class="section-content">
                        <div class="favourites-grid">
                            <div v-for="(row, rowIndex) in chunkedItems(grouped.onRequest)"
                                :key="`onRequest-row-${rowIndex}`" class="items-row">
                                <FavouriteItemCard v-for="(item, colIndex) in row"
                                    :key="`${item.uniqueId}|${item.size}`" :uniqueId="item.uniqueId"
                                    :title="$i18n.locale === 'en' ? item.title.en : item.title.ru" :image="item.image"
                                    :size="item.size" :price="item.price" :isOnRequest="item.isOnRequest"
                                    :quantity="item.quantity" :availability="item.availability" @remove="handleRemove"
                                    v-appear.repeat="{ delay: 200 + colIndex * 150 }" 
                                    :delay="600 + colIndex * 200"
                                    class="on-request-favourite-item"/>
                            </div>
                        </div>
                    </div>
                </transition>
            </div>

            <div class="section out-of-stock-section" v-if="grouped.outOfStock.length">
                <button class="section-header" @click="toggleSection('outOfStock')">
                    <h2>{{ $t('out_of_stock') }} ({{ grouped.outOfStock.length }})</h2>
                    <img src="../assets/img/down.png" :alt="$t('Down')" :class="{ 'rotated': expanded.outOfStock }">
                </button>
                <transition name="collapse">
                    <div v-show="expanded.outOfStock" class="section-content">
                        <div class="favourites-grid">
                            <div v-for="(row, rowIndex) in chunkedItems(grouped.outOfStock)"
                                :key="`outOfStock-row-${rowIndex}`" class="items-row">
                                <FavouriteItemCard v-for="(item, colIndex) in row"
                                    :key="`${item.uniqueId}|${item.size}`" :uniqueId="item.uniqueId"
                                    :title="$i18n.locale === 'en' ? item.title.en : item.title.ru" :image="item.image"
                                    :size="item.size" :price="item.price" :isOnRequest="item.isOnRequest"
                                    :quantity="item.quantity" :availability="item.availability" @remove="handleRemove"
                                    v-appear.repeat="{ delay: 200 + colIndex * 150 }" 
                                    :delay="600 + colIndex * 200"
                                    class="out-of-stock-favourite-item"/>
                            </div>
                        </div>
                    </div>
                </transition>
            </div>
        </div>
    </main>
</template>

<script setup>
import { useFavouritesStore } from '../stores/favourites';
import { computed, ref, onMounted } from 'vue';
import FavouriteItemCard from '../components/FavouriteItemCard.vue';

const favouritesStore = useFavouritesStore();
const emit = defineEmits(['page-loaded']);

const expanded = ref({
    inStock: true,
    onRequest: true,
    outOfStock: true
});

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

const toggleSection = (section) => {
    expanded.value[section] = !expanded.value[section];
};

const handleRemove = async (uniqueId, size) => {
    await favouritesStore.removeFromFavourites(uniqueId, size);
};

onMounted(async () => {
    await favouritesStore.loadFavourites();
    emit('page-loaded', true);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/favourites';
</style>