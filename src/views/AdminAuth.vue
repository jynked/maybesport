<template>
  <main class="auth-page">
    <div class="auth-card">
      <h1>Вход в панель администратора</h1>
      <form @submit.prevent="handleSubmit">
        <div class="input-group">
          <input type="email" v-model="email" placeholder="Email" required>
        </div>
        <div class="input-group">
          <div class="password-wrapper">
            <input :type="showPassword ? 'text' : 'password'" v-model="password" placeholder="Пароль" required>
            <button type="button" class="toggle-password" @click="showPassword = !showPassword">
              <img src="../assets/img/eye-lock.png" v-if="showPassword" alt="">
              <img src="../assets/img/eye.png" v-else alt="">
            </button>
          </div>
        </div>
        <button type="submit" :disabled="isLoading">Войти</button>
        <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      </form>
    </div>
  </main>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const email = ref('')
const password = ref('')
const isLoading = ref(false)
const errorMessage = ref('')
const showPassword = ref(false)

const emit = defineEmits(['page-loaded'])

async function handleSubmit() {
  isLoading.value = true
  errorMessage.value = ''
  try {
    await authStore.login(email.value, password.value)
    if (!authStore.isAdmin) {
      errorMessage.value = 'У вас нет прав администратора'
      authStore.logout()
      return
    }
    router.push('/admin/items')
  } catch (err) {
    errorMessage.value = 'Ошибка входа. Проверьте email и пароль.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => emit('page-loaded', true));
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/auth-user';
</style>