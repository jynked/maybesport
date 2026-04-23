<template>
    <div class="toast-container">
        <transition-group name="toast" tag="div" class="toast-list">
            <div v-for="msg in toastStore.messages" :key="msg.id" :class="['toast', msg.type]">
                <span>{{ msg.text }}</span>
            </div>
        </transition-group>
    </div>
</template>

<script setup>
import { useToastStore } from '../stores/toast';
const toastStore = useToastStore();
</script>

<style scoped lang="scss">
@use '../assets/styles/layout/toast';
</style>