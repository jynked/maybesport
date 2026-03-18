<template>
    <main>
        <h1 v-appear="{ delay: 300 }">{{ $t('altLK') }}</h1>
        <section class="user-wrapper">
            <div class="user-container">
                <div class="user-avatar">
                    <label for="userAvatar">
                        <img src="../assets/img/user.png" :alt="$t('UserAvatarAlt')">
                    </label>
                    <input type="file" name="userInfo" id="userAvatar" accept=".jpg, .png, .jpeg, .webp">
                </div>
                <div class="user-info">
                    <div class="user-info-block">
                        <label for="userEmail">{{ $t('email') }}</label>
                        <input type="text" name="userInfo" id="userEmail">
                    </div>
                    <div class="user-info-block">
                        <label for="userPassword">{{ $t('password') }}</label>
                        <input type="text" name="userInfo" id="userPassword">
                    </div>
                    <div class="user-info-block">
                        <label for="userName">{{ $t('userName') }}</label>
                        <input type="text" name="userInfo" id="userName">
                    </div>
                </div>
                <div class="action-buttons">
                    <button class="btn btn-outline-danger">Delete account</button>
                    <button class="btn" :class="{ 'redacting': !isRedacting, 'btn-warning': !isRedacting, 'btn-outline-warning': isRedacting }">Accept redact</button>
                </div>
            </div>
        </section>
    </main>
</template>

<script setup>
import { onMounted, ref } from 'vue';

const emit = defineEmits(['page-loaded']);
const isRedacting = ref(false);

onMounted(() => {
  emit('page-loaded', true);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/user';
</style>