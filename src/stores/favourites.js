import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useAuthStore } from './auth';
import { api } from '../api';

export const useFavouritesStore = defineStore('favourites', () => {
    const favouriteItems = ref([]);
    const favouriteKeys = ref([]);
    const loading = ref(false);

    let queue = Promise.resolve();

    function enqueue(fn) {
        const result = queue.then(() => fn());
        queue = result.catch(() => {});
        return result;
    }

    async function loadFavourites() {
        const authStore = useAuthStore();
        if (!authStore.isAuthenticated) {
            favouriteItems.value = [];
            favouriteKeys.value = [];
            return;
        }
        loading.value = true;
        try {
            const response = await api.getFavouriteItems();
            if (response.data !== null) {
                favouriteItems.value = response.data;
                favouriteKeys.value = response.data.map(item => item.uniqueId);
            }
        } catch (error) {
            console.error('Failed to load favourites', error);
        } finally {
            loading.value = false;
        }
    }

    async function addToFavourites(uniqueId) {
        return enqueue(async () => {
            if (favouriteKeys.value.includes(uniqueId)) return;
            const oldKeys = [...favouriteKeys.value];
            const oldItems = [...favouriteItems.value];

            favouriteKeys.value.push(uniqueId);
            try {
                await api.addToFavourites(uniqueId);
                await loadFavourites();
            } catch (error) {
                favouriteKeys.value = oldKeys;
                favouriteItems.value = oldItems;
                console.error('Add to favourites failed', error);
                throw error;
            }
        });
    }

    async function removeFromFavourites(uniqueId) {
        return enqueue(async () => {
            if (!favouriteKeys.value.includes(uniqueId)) return;
            const oldKeys = [...favouriteKeys.value];
            const oldItems = [...favouriteItems.value];

            favouriteKeys.value = favouriteKeys.value.filter(k => k !== uniqueId);
            favouriteItems.value = favouriteItems.value.filter(item => item.uniqueId !== uniqueId);

            try {
                await api.removeFromFavourites(uniqueId);
            } catch (error) {
                favouriteKeys.value = oldKeys;
                favouriteItems.value = oldItems;
                console.error('Remove from favourites failed', error);
                throw error;
            }
        });
    }

    function isFavourite(uniqueId) {
        return favouriteKeys.value.includes(uniqueId);
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