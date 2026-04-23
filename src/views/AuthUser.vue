<template>
  <main class="auth-page">
    <div class="auth-card" v-appear="{ delay: 200 }">
      <div class="auth-header">
        <h1>{{ currentAuth === 'auth' ? $t('auth') : $t('register') }}</h1>
        <button @click="changeAuth" class="switch-btn">
          {{ currentAuth === 'auth' ? $t('register') : $t('auth') }}
        </button>
      </div>

      <form @submit.prevent="handleSubmit">
        <div class="input-group">
          <input type="email" id="email" v-model="userForm.email" required :placeholder="$t('email')">
        </div>
        <div class="input-group">
          <div class="password-wrapper">
            <input :type="showPassword ? 'text' : 'password'" id="password" v-model="userForm.password" required
              :placeholder="$t('password')">
            <button type="button" class="toggle-password" @click="togglePasswordVisibility">
              <img src="../assets/img/eye-lock.png" alt="" v-if="showPassword">
              <img src="../assets/img/eye.png" alt="" v-else>
            </button>
          </div>
        </div>
        <div v-if="currentAuth === 'register'" class="password-strength">
          <progress :value="passwordStrength" max="100"></progress>
          <span>{{ strengthText }}</span>
        </div>

        <button type="submit" class="submit-btn" :disabled="!isFormValid || isLoading">
          <span v-if="isLoading" class="spinner"></span>
          <span v-else>{{ currentAuth === 'auth' ? $t('send') : $t('register') }}</span>
        </button>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>
      </form>
    </div>
  </main>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useI18n } from 'vue-i18n';

const emit = defineEmits(['page-loaded']);
const authStore = useAuthStore();

const { t } = useI18n();

const userForm = ref({ email: '', password: '' });
const currentAuth = ref('auth');
const errorMessage = ref('');
const isLoading = ref(false);
const showPassword = ref(false);

function togglePasswordVisibility() {
  showPassword.value = !showPassword.value;
}

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
    if (error.response && error.response.data) {
      errorMessage.value = error.response.data;
    } else {
      errorMessage.value = currentAuth.value === 'auth' ? t('errorLog') : t('errorReg');
    }
  } finally {
    isLoading.value = false;
  }
}

const isEmailValid = (email) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
const isFormValid = computed(() => isEmailValid(userForm.value.email) && userForm.value.password.trim() !== '');

const passwordStrength = computed(() => {
  const pwd = userForm.value.password;
  let score = 0;
  if (pwd.length >= 8) score += 25;
  if (/[A-Z]/.test(pwd)) score += 25;
  if (/[a-z]/.test(pwd)) score += 25;
  if (/[0-9]/.test(pwd)) score += 25;
  return score;
});

const strengthText = computed(() => {
  if (passwordStrength.value < 30) return t('weakPassword');
  if (passwordStrength.value < 60) return t('mediumPassword');
  return t('strongPassword');
});

onMounted(() => emit('page-loaded', true));
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/auth-user';
</style>