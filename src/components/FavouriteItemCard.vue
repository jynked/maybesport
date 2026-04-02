<template>
    <button class="item-card" @click="goToItem">
        <div class="item-card-swiper">
            <div class="item-card-image">
                <img :src="image" :alt="title">
            </div>
        </div>

        <div class="card-tags-block">
            <p v-appear.repeat="{ delay: props.delay + 200 }">{{ title }}</p>
            <p v-appear.repeat="{ delay: props.delay + 250 }">{{ $t('size') }}: {{ size }}</p>
            <p v-appear.repeat="{ delay: props.delay + 300 }">{{ price.toLocaleString() }} ₽</p>
        </div>

        <button @click.stop="remove" class="remove-btn">
            <img src="../assets/img/favourite.png" :alt="$t('favouriteAlt')"/>
        </button>
    </button>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const props = defineProps({
    uniqueId: String,
    title: String,
    image: String,
    size: [String, Number],
    price: Number,
    isOnRequest: Boolean,
    quantity: Number,
    availability: String,
    delay: Number,
});

const emit = defineEmits(['remove']);
const router = useRouter();

const remove = () => {
    emit('remove', props.uniqueId, props.size);
};

const goToItem = () => {
    router.push({ name: 'Item', params: { itemId: props.uniqueId } });
};

const statusText = computed(() => {
    if (props.quantity > 0 && !props.isOnRequest) return 'В наличии';
    if (props.quantity > 0 && props.isOnRequest) return 'Под заказ';
    return 'Нет в наличии';
});

const statusClass = computed(() => {
    if (props.quantity > 0 && !props.isOnRequest) return 'available';
    if (props.quantity > 0 && props.isOnRequest) return 'on-request';
    return 'out-of-stock';
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/favourite-item';
</style>