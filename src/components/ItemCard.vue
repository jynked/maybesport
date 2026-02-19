<template>
    <a href="javascript:" @click="handleCardClick($event)" class="item-card"
        :style="{ opacity: 0, transform: 'translateY(20px)' }">
        <swiper direction="horizontal" slides-per-view="1" space-between="20" :modules="modules"
            :pagination="pagination" class="item-card-swiper" v-if="props.images.length > 1">
            <swiper-slide v-for="(image, index) in props.images" :key="index" class="item-card-image">
                <img :src="image" :alt="$t('itemImageAlt')">
            </swiper-slide>

            <div class="swiper-pagination"></div>
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
            <p>{{ $t('from') }} {{ Number(props.minPrice).toLocaleString('ru-RU') }} ₽</p>
            <p
                :style="{ 'color': props.availability == 'available' ? 'rgb(0 198 99)' : props.availability == 'out_of_stock' ? 'red' : '#5a5aff' }">
                {{ $t(`${props.availability}`) }}</p>
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
    availability: [String, Number],
    minPrice: [String, Number],
    tags: Array,
    uniqueId: [String, Number],
    delay: Number,
});

function handleCardClick(e) {
    if (isSwiping.value) {
        e.preventDefault();
        return
    }

    router.push({
        name: 'Item',
        params: { itemId: props.uniqueId }
    })
}
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/item-card';
</style>
