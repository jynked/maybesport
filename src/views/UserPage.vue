<template>
    <main>
        <h1 v-appear="{ delay: 300 }">{{ $t('altLK') }}</h1>
    </main>
</template>

<script setup>
import { onMounted } from 'vue';

const emit = defineEmits(['page-loaded']);

onMounted(() => {
  emit('page-loaded', true);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/user';
</style>