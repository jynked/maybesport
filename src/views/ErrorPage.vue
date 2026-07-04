<template>
    <main>
        <div class="error-image">
            <img src="../assets/img/404.png" :alt="$t('lampAlt')" :class="{
                flicker: isFlickering,
                burnt: isBurnt,
                clickable: isBurnt
            }" @click="handleBulbClick" v-appear="{ delay: 400 }" />
        </div>
        <div class="error-text-container">
            <h1 v-appear="{ delay: 500 }">404</h1>
            <h4 v-appear="{ delay: 600 }">{{ $t('errorMessage') }}</h4>
            <p v-appear="{ delay: 700 }">
                {{ $t('errorDescription') }}
                <router-link :to="{ name: 'Catalog' }">{{ $t('catalog').toLowerCase() }}</router-link>
            </p>
        </div>
    </main>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';

const emit = defineEmits(['page-loaded']);

const isFlickering = ref(false);
const isBurnt = ref(false);

let flickerInterval = null;
let mainTimer = null;

const startBulbCycle = () => {
    isFlickering.value = false;
    isBurnt.value = false;

    if (mainTimer) clearTimeout(mainTimer);
    if (flickerInterval) clearInterval(flickerInterval);

    mainTimer = setTimeout(() => {
        startFlickering();
    }, 4000);
};

const startFlickering = () => {
    isFlickering.value = true;

    setTimeout(() => {
        isFlickering.value = false;
        burnOut();
    }, 500);
};

const burnOut = () => {
    isBurnt.value = true;
};

const handleBulbClick = () => {
    if (!isBurnt.value) return;
    startBulbCycle();
};

onBeforeUnmount(() => {
    if (mainTimer) clearTimeout(mainTimer);
    if (flickerInterval) clearInterval(flickerInterval);
});

onMounted(() => {
    emit('page-loaded', true);
    startBulbCycle();
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/error';
</style>