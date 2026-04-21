<template>
  <main>
    <h1 v-appear="{ delay: 300 }">{{ $t('altLK') }}</h1>

    <div class="profile-grid">
      <div class="profile-card" v-appear="{ delay: 400 }">
        <div class="avatar-section">
          <label for="avatarUpload" class="avatar-label" :class="{ 'disabled': !isEditing }">
            <img :src="avatarPreview || defaultAvatar" alt="Avatar">
            <div class="avatar-overlay" v-if="isEditing">
              <span>✎</span>
            </div>
          </label>
          <input type="file" id="avatarUpload" accept="image/*" @change="handleAvatarUpload" style="display: none">
        </div>

        <div class="profile-fields">
          <div class="field">
            <label for="userName">{{ $t('userName') }}</label>
            <input type="text" id="userName" v-model="editForm.name" :disabled="!isEditing"
              :class="{ 'editing': isEditing }">
          </div>
          <div class="field">
            <label for="userEmail">{{ $t('email') }}</label>
            <input type="email" id="userEmail" v-model="editForm.email" :disabled="!isEditing"
              :class="{ 'editing': isEditing }">
          </div>
          <div class="field">
            <label for="userPassword">{{ $t('password') }}</label>
            <div class="password-wrapper">
              <input :type="showPassword ? 'text' : 'password'" id="userPassword" v-model="editForm.password"
                :disabled="!isEditing" :class="{ 'editing': isEditing }" placeholder="••••••••">
              <button type="button" class="toggle-password" @click="togglePasswordVisibility" v-if="isEditing">
                <img src="../assets/img/eye-lock.png" alt="" v-if="showPassword">
                <img src="../assets/img/eye.png" alt="" v-else>
              </button>
            </div>
          </div>
        </div>

        <div class="profile-actions">
          <button v-if="!isEditing" class="btn btn-outline" @click="startEditing">
            {{ $t('editProfile') }}
          </button>
          <button v-else class="btn btn-primary" @click="saveChanges">
            {{ $t('saveChanges') }}
          </button>
          <button class="btn btn-danger" @click="openDeleteModal">
            {{ $t('deleteAccount') }}
          </button>
        </div>
      </div>

      <div class="widgets-grid" v-appear="{ delay: 500 }">
        <router-link :to="{ name: 'Favourites' }" class="widget-card">
          <div class="widget-icon">❤️</div>
          <h3>{{ $t('favouritesHeaderNav') }}</h3>
          <p v-if="favouritesCount !== null">{{ favouritesCount }} {{ favouriteWord }}</p>
          <p v-else>{{ $t('loading') }}...</p>
          <span class="arrow">→</span>
        </router-link>

        <router-link :to="{ name: 'UserOrders' }" class="widget-card">
          <div class="widget-icon">📦</div>
          <h3>{{ $t('ordersHeaderNav') }}</h3>
          <p>{{ $t('lastOrders') }}</p>
          <span class="arrow">→</span>
        </router-link>

        <div class="widget-card" @click="openCartModal">
          <div class="widget-icon">🛒</div>
          <h3>{{ $t('cartHeaderNav') }}</h3>
          <p>{{ $t('viewCart') }}</p>
          <span class="arrow">→</span>
        </div>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="showDeleteModal" class="modal-overlay" @click.self="closeDeleteModal">
        <div class="modal-confirm">
          <h3>{{ $t('confirmDeleteTitle') }}</h3>
          <p>{{ $t('confirmDeleteMessage') }}</p>
          <div class="modal-actions">
            <button class="btn btn-secondary" @click="closeDeleteModal">{{ $t('cancel') }}</button>
            <button class="btn btn-danger" @click="deleteAccount">{{ $t('delete') }}</button>
          </div>
        </div>
      </div>
    </Transition>
    <CartModal :isOpen="isCartModalOpen" @close="closeCartModal" />
  </main>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useFavouritesStore } from '../stores/favourites';
import { useRouter } from 'vue-router';
import { api } from '../api';
import defaultAvatar from '/src/assets/img/user.png';
import { useI18n } from 'vue-i18n';
import CartModal from '../components/CartModal.vue';

const { locale } = useI18n();

const emit = defineEmits(['page-loaded']);
const authStore = useAuthStore();
const favouritesStore = useFavouritesStore();
const router = useRouter();

const avatarPreview = ref(null);
const isEditing = ref(false);
const favouritesCount = ref(null);
const showPassword = ref(false);
const showDeleteModal = ref(false);

const editForm = ref({
  name: authStore.user?.name || '',
  email: authStore.user?.email || '',
  password: ''
});

const favouriteWord = computed(() => {
  const count = favouritesCount.value;
  if (count === null) return '';
  const lang = locale.value || 'ru';
  if (lang === 'ru') {
    if (count % 10 === 1 && count % 100 !== 11) return 'товар';
    if (count % 10 >= 2 && count % 10 <= 4 && (count % 100 < 10 || count % 100 >= 20)) return 'товара';
    return 'товаров';
  } else {
    return count === 1 ? 'item' : 'items';
  }
});

const isCartModalOpen = ref(false);

function openCartModal() {
  isCartModalOpen.value = true;
}
function closeCartModal() {
  isCartModalOpen.value = false;
}

function togglePasswordVisibility() {
  showPassword.value = !showPassword.value;
}

function handleAvatarUpload(event) {
  const file = event.target.files[0];
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => { avatarPreview.value = e.target.result; };
    reader.readAsDataURL(file);
  }
}

function startEditing() {
  isEditing.value = true;
  editForm.value = {
    name: authStore.user?.name || '',
    email: authStore.user?.email || '',
    password: ''
  };
}

async function saveChanges() {
  try {
    const updateData = {};
    if (editForm.value.name !== authStore.user?.name) updateData.name = editForm.value.name;
    if (editForm.value.email !== authStore.user?.email) updateData.email = editForm.value.email;
    if (editForm.value.password) updateData.password = editForm.value.password;

    if (Object.keys(updateData).length > 0) {
      const response = await api.updateProfile(updateData);
      authStore.user = response.data;
      localStorage.setItem('user', JSON.stringify(response.data));
    }
    isEditing.value = false;
  } catch (error) {
    console.error('Update profile error', error);
    alert('Ошибка при обновлении профиля');
  }
}

function openDeleteModal() {
  showDeleteModal.value = true;
}

function closeDeleteModal() {
  showDeleteModal.value = false;
}

async function deleteAccount() {
  try {
    await api.deleteAccount();
    authStore.logout();
    router.push('/');
  } catch (error) {
    console.error('Delete account error', error);
    alert('Ошибка при удалении аккаунта');
  } finally {
    closeDeleteModal();
  }
}

async function loadFavouritesCount() {
  if (authStore.isAuthenticated) {
    try {
      await favouritesStore.loadFavourites();
      favouritesCount.value = favouritesStore.favouriteItems.length;
    } catch (e) {
      favouritesCount.value = 0;
    }
  } else {
    favouritesCount.value = 0;
  }
}

onMounted(() => {
  loadFavouritesCount();
  emit('page-loaded', true);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/user';
</style>