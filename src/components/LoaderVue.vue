<template>
    <section class="page-loader" ref="loaderRef">
        <div class="loader-block"></div>
        <div class="loader-block">
            <p class="loader-text">{{ loader }}</p>
        </div>
        <div class="loader-block"></div>
    </section>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';

const props = defineProps({
    isPageLoaded: {
        required: false,
        type: Boolean,
        default: null,
    }
})

const loaderRef = ref(null);
const loader = ref('Maybesport');
document.body.style.overflow = 'hidden';

watch(() => props.isPageLoaded, (newVal) => {
    if (newVal === true) {
        setTimeout(() => {
            document.body.style.overflow = 'auto';

            loaderRef.value.querySelectorAll('.loader-block').forEach(element => {
                element.style.width = '0%';
            });

            loaderRef.value.querySelector('.loader-text').style.opacity = '0';
        }, 400);
        
        setTimeout(() => {
            if (loaderRef.value) {
                loaderRef.value.style.display = 'none';
            }
        }, 1000);
    } else {
        loader.value = 'Что то пошло не так! Перезагрузите страницу и попробуйте позже.';
    }
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/layout/loader';
</style>
