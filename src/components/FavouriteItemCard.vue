<template>
    <div class="item-card" @click="goToItem">
        <div class="item-card-swiper">
            <div class="item-card-image">
                <img :src="image" :alt="title">
            </div>
        </div>
        <div class="card-tags-block">
            <p>{{ title }}</p>
            <p>{{ price.toLocaleString() }} ₽</p>
            <p :class="'status-' + availability">{{ $t(availability) }}</p>
        </div>
        <div class="favourite-action-buttons">
            <button @click.stop="remove" class="remove-btn">{{ $t('remove') }}</button>
            <button @click.stop="addToCart" class="cart-btn">{{ $t('inCart') }}</button>
        </div>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const props = defineProps({
    uniqueId: String,
    title: String,
    image: String,
    price: Number,
    isOnRequest: Boolean,
    availability: String,
    delay: Number,
});

const emit = defineEmits(['remove', 'addToCart']);
const router = useRouter();

const remove = () => emit('remove', props.uniqueId);
const addToCart = () => emit('addToCart', props.uniqueId);
const goToItem = () => router.push({ name: 'Item', params: { itemId: props.uniqueId } });
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/favourite-item';
</style>