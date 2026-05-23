<template>
    <div class="item-card" @click="goToItem">
        <button class="card-action-btn" @click.stop="toggleFavourite">
            <svg class="favourite-icon" :class="{ 'favourite-active': isFav }" viewBox="0 0 24 24" width="22"
                height="22" fill="none" stroke="currentColor" stroke-width="2">
                <path
                    d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
            </svg>
        </button>

        <swiper direction="horizontal" slides-per-view="1" space-between="20" :modules="modules"
            :pagination="pagination" class="item-card-swiper" v-if="props.images.length > 1">
            <swiper-slide v-for="(image, index) in props.images" :key="index" class="item-card-image">
                <img :src="image" :alt="$t('itemImageAlt')">
            </swiper-slide>
            <div class="swiper-pagination" @click.stop></div>
        </swiper>
        <div class="item-card-swiper" v-else>
            <div class="item-card-image">
                <img :src="props.images[0]" :alt="$t('itemImageAlt')">
            </div>
        </div>

        <div class="card-tags-block">
            <p v-for="(tag, index) in props.tags" :key="tag" v-appear="{ delay: props.delay + 200 * index }">
                {{ $i18n.locale == 'en' ? tag.en : tag.ru }}
            </p>
        </div>

        <div class="card-info">
            <h3>{{ $i18n.locale == 'en' ? props.title.en : props.title.ru }}</h3>
            <p>
                {{ $t('from') }} {{ Number(props.minPrice).toLocaleString('ru-RU') }} ₽
            </p>
            <button class="buy-button" @click.stop="openModal">
                {{ $t('inCart') }}
                <span class="cart-badge" v-if="totalInCart > 0">{{ totalInCart > 99 ? '99+' : totalInCart }}</span>
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useFavouritesStore } from '../stores/favourites';
import { useCartStore } from '../stores/cart';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { Pagination } from 'swiper/modules';
import 'swiper/css';
import 'swiper/css/pagination';

const modules = [Pagination];
const router = useRouter();
const { t } = useI18n();

const props = defineProps({
    id: Number,
    title: Object,
    images: Array,
    color: Object,
    sizes: Array,
    availability: [String, Number],
    minPrice: [String, Number],
    tags: Array,
    uniqueId: [String, Number],
    delay: Number,
});

const emit = defineEmits(['openSizeModal']);

const favouritesStore = useFavouritesStore();
const cartStore = useCartStore();
const authStore = useAuthStore();
const toast = useToastStore();

const isFav = computed(() => favouritesStore.isFavourite(props.uniqueId));
const totalInCart = computed(() => cartStore.getTotalQuantityByUniqueId(props.uniqueId));

async function toggleFavourite() {
    if (!authStore.isAuthenticated) {
        localStorage.setItem('pendingAction', JSON.stringify({
            action: 'favourite',
            uniqueId: props.uniqueId
        }));
        router.push({ name: 'UserAuth', query: { redirect: router.currentRoute.value.fullPath } });
        return;
    }
    if (isFav.value) {
        await favouritesStore.removeFromFavourites(props.uniqueId);
        toast.success(t('removedFromFavourites'));
    } else {
        await favouritesStore.addToFavourites(props.uniqueId);
        toast.success(t('addedToFavourites'));
    }
}

function openModal() {
    emit('openSizeModal', {
        uniqueId: props.uniqueId,
        sizes: props.sizes
    });
}

function goToItem() {
    router.push({ name: 'Item', params: { itemId: props.uniqueId } });
}

const pagination = {
    el: '.swiper-pagination',
    clickable: true,
    bulletClass: 'custom-bullet',
    bulletActiveClass: 'custom-bullet-active'
};
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/item-card';
</style>