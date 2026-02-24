<template>
    <section class="page-loader" ref="loaderRef" :class="{ 'loading': isLoading, 'loaded': isLoaded }">
        <div class="loader-block"></div>
        <div class="loader-block">
            <p class="loader-text">{{ loader }}</p>
        </div>
        <div class="loader-block"></div>
    </section>
</template>

<script setup>
import { ref, watch, onMounted, defineExpose } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps({
    isPageLoaded: {
        required: false,
        type: Boolean,
        default: null,
    }
})

const loaderRef = ref(null);
const loader = ref('Maybesport');
const isLoading = ref(false);
const isLoaded = ref(false);

const startLoading = () => {
    loaderRef.value.style.display = 'grid';
    loaderRef.value.style.opacity = '0';
    isLoading.value = true;
    isLoaded.value = false;
    document.body.style.overflow = 'hidden';

    loaderRef.value.style.transition = 'opacity 0.3s ease';
    loaderRef.value.style.opacity = '1';

    loaderRef.value.querySelectorAll('.loader-block').forEach(element => {
        element.style.width = '100%';
    });

    loaderRef.value.querySelector('.loader-text').style.opacity = '1';
};

const finishLoading = () => {
    isLoaded.value = true;
    loaderRef.value.querySelectorAll('.loader-block').forEach(element => {
        element.style.width = '0%';
    });

    loaderRef.value.querySelector('.loader-text').style.opacity = '0';

    setTimeout(() => {
        document.body.style.overflow = 'auto';
        if (loaderRef.value) {
            loaderRef.value.style.opacity = '0';
            setTimeout(() => {
                isLoading.value = false;
            }, 300);
        }
    }, 100);
};

defineExpose({
    startLoading,
    finishLoading
});

watch(() => props.isPageLoaded, (newVal) => {
    if (newVal === true) {
        finishLoading();
    } else if (newVal === false) {
        loader.value = t('errorLoader');
    }
});

onMounted(() => {
    document.body.style.overflow = 'hidden';
    if (loaderRef.value) {
        loaderRef.value.style.display = 'none';
    }
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/layout/loader';
</style>