import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from '../api';

export const useFavouritesStore = defineStore('favourites', () => {
    const favouriteItems = ref([]);
    const favouriteKeys = ref([]);
    const loading = ref(false);

    async function loadFavourites() {
        loading.value = true;
        try {
            const response = await api.getFavouriteItems();
            console.log(response.data);
            favouriteItems.value = response.data;
            favouriteKeys.value = response.data.map(item => `${item.uniqueId}|${item.size}`);
        } catch (error) {
            console.error('Failed to load favourites', error);
        } finally {
            loading.value = false;
        }
    }

    async function addToFavourites(uniqueId, size) {
        const key = `${uniqueId}|${size}`;
        if (favouriteKeys.value.includes(key)) return;

        const oldKeys = [...favouriteKeys.value];
        favouriteKeys.value.push(key);

        try {
            await api.addToFavourites(uniqueId, size);
            await loadFavourites();
        } catch (error) {
            favouriteKeys.value = oldKeys;
            console.error('Add to favourites failed', error);
            throw error;
        }
    }

    async function removeFromFavourites(uniqueId, size) {
        const key = `${uniqueId}|${size}`;
        if (!favouriteKeys.value.includes(key)) return;

        const oldItems = [...favouriteItems.value];
        const oldKeys = [...favouriteKeys.value];

        favouriteItems.value = favouriteItems.value.filter(
            item => !(item.uniqueId === uniqueId && String(item.size) === String(size))
        );
        favouriteKeys.value = favouriteKeys.value.filter(k => k !== key);

        try {
            await api.removeFromFavourites(uniqueId, size);
        } catch (error) {
            favouriteItems.value = oldItems;
            favouriteKeys.value = oldKeys;
            console.error('Remove from favourites failed', error);
            throw error;
        }
    }

    function isFavourite(uniqueId, size) {
        return favouriteKeys.value.includes(`${uniqueId}|${size}`);
    }

    return {
        favouriteItems,
        favouriteKeys,
        loading,
        loadFavourites,
        addToFavourites,
        removeFromFavourites,
        isFavourite
    };
});