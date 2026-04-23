<template>
  <LoaderVue ref="loader" :isPageLoaded="isAppLoaded"/>
  <HeaderMain/>
  <router-view @page-loaded="handlePageLoaded"/>
  <FooterMain/>
  <ToastNotification />
</template>

<script setup>
import { ref, onMounted } from 'vue';
import HeaderMain from './components/HeaderMain.vue';
import FooterMain from './components/FooterMain.vue';
import LoaderVue from './components/LoaderVue.vue';
import ToastNotification from './components/ToastNotification.vue';
import { setLoaderInstance } from './router';
import { useAuthStore } from './stores/auth';
import { useI18n } from 'vue-i18n';
import { useToastStore } from './stores/toast';
import { useFavouritesStore } from './stores/favourites';

const loader = ref(null);
const isAppLoaded = ref(false);

const { t } = useI18n();

const handlePageLoaded = (loaded) => {
  if (loader.value) {
    if (loaded) {
      loader.value.finishLoading();
    } else {
      loader.value.finishLoading();
      useToastStore().error(t('pageLoadError'));
    }
  }
};

onMounted(() => {
  setLoaderInstance(loader.value);

  const authStore = useAuthStore();
  const favouritesStore = useFavouritesStore();
  authStore.fetchUser();
  favouritesStore.loadFavourites();

  setTimeout(() => {
    isAppLoaded.value = true;
  }, 1000);
});
</script>

<style>
@import url('./assets/styles/style.scss');
</style>