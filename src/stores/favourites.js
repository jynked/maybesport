import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from '../api';

export const useFavouritesStore = defineStore('favourites', () => {
    const favouriteIds = ref([]);
    const favouriteItems = ref([]);
    const loading = ref(false);

    async function loadFavourites() {
        loading.value = true;
        try {
            const response = await api.getFavouriteItems();
            favouriteItems.value = response.data;
            favouriteIds.value = favouriteItems.value.map(item => item.uniqueId);
        } catch (error) {
            console.error('Failed to load favourites', error);
        } finally {
            loading.value = false;
        }
    }

    async function addToFavourites(uniqueId) {
        try {
            await api.addToFavourites(uniqueId);
            if (!favouriteIds.value.includes(uniqueId)) {
                await loadFavourites();
            }
        } catch (error) {
            console.error('Add to favourites failed', error);
        }
    }

    async function removeFromFavourites(uniqueId) {
        try {
            await api.removeFromFavourites(uniqueId);
            favouriteIds.value = favouriteIds.value.filter(id => id !== uniqueId);
            favouriteItems.value = favouriteItems.value.filter(item => item.uniqueId !== uniqueId);
        } catch (error) {
            console.error('Remove from favourites failed', error);
        }
    }

    function isFavourite(uniqueId) {
        return favouriteIds.value.includes(uniqueId);
    }

    return {
        favouriteIds,
        favouriteItems,
        isFavourite,

        loading,
        loadFavourites,
        addToFavourites,
        removeFromFavourites
    };
});