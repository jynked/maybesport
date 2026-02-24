<template>
    <main class="auth-user-wrapper">
        <div class="title">
            <h1 v-appear="{ delay: 200 }">{{ currentAuth == 'auth' ? $t('auth') : $t('register') }}</h1>
            <button v-appear="{ delay: 300 }" @click="changeAuth">{{ currentAuth == 'auth' ? $t('register') : $t('auth')
                }}</button>
        </div>
        <form id="authUser" @submit.prevent="handleSubmit" v-appear="{ delay: 250 }">
            <div class="row mb-3">
                <div class="form-floating col-md-6" v-appear="{ delay: 250 }">
                    <input type="email" class="form-control" id="floatingInput" :placeholder="$t('email')"
                        v-model="userForm.email">
                    <label for="floatingInput">{{ $t('email') }}</label>
                </div>
                <div class="form-floating col-md-6" v-appear="{ delay: 350 }">
                    <input type="password" class="form-control" id="floatingPassword" :placeholder="$t('password')"
                        v-model="userForm.password">
                    <label for="floatingPassword">{{ $t('password') }}</label>
                </div>
            </div>
            <button class="decorated-button" :disabled="!isFormValid || isLoading">
                <span v-if="isLoading" class="spinner-border spinner-border-sm me-1"></span>
                {{ currentAuth === 'auth' ? $t('send') : $t('register') }}
            </button>
        </form>
    </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useAuthStore } from '../stores/auth';

const emit = defineEmits(['page-loaded']);

const authStore = useAuthStore();

const userForm = ref({
    email: '',
    password: '',
});

const currentAuth = ref('auth');
const errorMessage = ref('');
const isLoading = ref(false);

function changeAuth() {
    userForm.value = { email: '', password: '' };
    errorMessage.value = '';
    currentAuth.value = currentAuth.value === 'auth' ? 'register' : 'auth';
}

async function handleSubmit() {
    errorMessage.value = '';
    isLoading.value = true;

    try {
        if (currentAuth.value === 'auth') {
            await authStore.login(userForm.value.email, userForm.value.password);
        } else {
            await authStore.register(userForm.value.email, userForm.value.password);
        }
    } catch (error) {
        errorMessage.value = error.response?.data?.message;
    } finally {
        isLoading.value = false;
    }
}

const isEmailValid = (email) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

const isFormValid = computed(() => {
    return isEmailValid(userForm.value.email) && userForm.value.password.trim() !== '';
});

onMounted(() => {
  emit('page-loaded', true);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/auth-user';
</style>