<template>
    <main>
        <h1 v-appear="{ delay: 200 }">{{ $t('favouritesHeaderNav') }}</h1>
        <div v-if="favouritesStore.loading" class="loader">Loading...</div>
        <div v-else-if="favouritesStore.favouriteItems.length === 0" class="undefined-items-container" v-appear="{ delay: 400 }">
            <p>{{ $t('undefinedItems') }}</p>
            <img src="../assets/img/fail.png" :alt="$t('failAlt')">
        </div>
        <div class="catalog-items" v-else>
            <div v-for="(row, rowIndex) in chunkedItems" :key="`row-${rowIndex}`" class="items-row">
                <ItemCard v-for="(item, colIndex) in row" :key="item.uniqueId"
                    v-appear="{ delay: 100 + colIndex * 100 }" :id="item.id" :title="item.title" :images="item.images"
                    :color="item.color" :sizes="item.sizes" :availability="item.availability" :uniqueId="item.uniqueId"
                    :minPrice="item.minPrice" :tags="item.tags" :delay="600 + colIndex * 200" />
            </div>
        </div>
    </main>
</template>

<script setup>
import { useFavouritesStore } from '../stores/favourites';
import { onMounted, computed } from 'vue';
import ItemCard from '../components/ItemCard.vue';

const favouritesStore = useFavouritesStore();
const emit = defineEmits(['page-loaded']);

const chunkedItems = computed(() => {
    console.log(favouritesStore.favouriteItems);
    const items = favouritesStore.favouriteItems;
    const chunkSize = 3;
    const rows = [];
    for (let i = 0; i < items.length; i += chunkSize) {
        rows.push(items.slice(i, i + chunkSize));
    }
    return rows;
});

onMounted(async () => {
    await favouritesStore.loadFavourites();
    emit('page-loaded', true);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/favourites';
</style>