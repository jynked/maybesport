<template>
    <a href="javascript:" @click="handleCardClick($event)" class="item-card" :class="{ 'out-of-stock-card': props.availability == 'out_of_stock' }">
        <swiper
            direction="horizontal" 
            slides-per-view="1" 
            space-between="20"
            :modules="modules"
            :pagination="pagination"
            class="item-card-swiper"
        >
            <swiper-slide v-for="(image, index) in props.images" :key="index" class="item-card-image">
                <img :src="image" :alt="$t('itemImageAlt')">
            </swiper-slide>

            <div class="swiper-pagination"></div>
        </swiper>

        <div class="card-info">
            <h3>{{ $i18n.locale === 'en' ? props.title.en : props.title.ru }}</h3>
            <p>{{ $t('from') }} {{ Number(props.minPrice).toLocaleString('ru-RU') }} {{ props.sizes[0]?.currency }}</p>
            <p :style="{ 'color': props.availability == 'available' ? 'rgb(0 198 99)' : props.availability == 'out_of_stock' ? 'red' : '#5a5aff' }">{{ $t(`${props.availability}`) }}</p>
        </div>
    </a>
</template>

<script setup>
import { Swiper, SwiperSlide } from 'swiper/vue';
import { Pagination } from 'swiper/modules';
import 'swiper/css';
import 'swiper/css/pagination';
import { useRouter } from 'vue-router';
import { ref } from 'vue';

const modules = [Pagination];
const router = useRouter();
const isSwiping = ref(false);

const pagination = {
    el: '.swiper-pagination',
    clickable: true,
    bulletClass: 'custom-bullet',
    bulletActiveClass: 'custom-bullet-active'
};

const props = defineProps({
    id: Number,
    title: Object,
    images: Array,
    color: Object,
    sizes: Array,
    availability: String,
    minPrice: Number,
});

function handleSwiperClick() {
    isSwiping.value = true;
    
    setTimeout(() => {
        isSwiping.value = false;
    }, 100);
}

function handleCardClick(e) {
    if (props.availability == 'out_of_stock') {
        e.preventDefault();
        return;
    }
    
    if (isSwiping.value) {
        e.preventDefault();
        console.log(isSwiping);
        return;
    }
    
    // router.push({ name: 'Item', params: { id: props.id } });
}
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/item-card';
</style>
