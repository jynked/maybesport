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
import { ref, watch, onMounted } from 'vue';
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

const isMobileLayout = () => window.innerWidth < 768;

const startLoading = () => {
    loaderRef.value.style.display = 'grid';
    loaderRef.value.style.opacity = '0';
    isLoading.value = true;
    isLoaded.value = false;
    document.body.style.overflow = 'hidden';

    loaderRef.value.style.transition = 'opacity 0.3s ease';
    loaderRef.value.style.opacity = '1';

    const blocks = loaderRef.value.querySelectorAll('.loader-block');
    const textEl = loaderRef.value.querySelector('.loader-text');
    
    if (isMobileLayout()) {
        blocks.forEach(element => {
            element.style.height = '100%';
            element.style.width = '100%';
        });
    } else {
        blocks.forEach(element => {
            element.style.width = '100%';
        });
    }
    textEl.style.opacity = '1';
};

const finishLoading = () => {
    isLoaded.value = true;
    const blocks = loaderRef.value.querySelectorAll('.loader-block');
    const textEl = loaderRef.value.querySelector('.loader-text');
    
    if (isMobileLayout()) {
        blocks.forEach(element => {
            element.style.height = '0%';
        });
    } else {
        blocks.forEach(element => {
            element.style.width = '0%';
        });
    }
    textEl.style.opacity = '0';

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