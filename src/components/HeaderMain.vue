<template>
    <header v-appear="{ delay: 1200 }">
        <div class="header">
            <div v-appear="{ delay: 1400 }">
                <router-link :to="{ name: 'Main' }" class="header-logo">
                    <img src="../assets/img/logo.png" :alt="$t('altLogo')" style="filter: invert(1)">
                </router-link>
            </div>
            <nav>
                <router-link :to="{ name: 'Catalog' }" v-appear="{ delay: 1500 }">
                    {{ $t('catalog') }}
                    <span></span>
                </router-link>
                <a href="#" v-appear="{ delay: 1550 }">
                    {{ $t('ordersHeaderNav') }}
                    <span></span>
                </a>
                <a href="#" v-appear="{ delay: 1600 }">
                    {{ $t('favouritesHeaderNav') }}
                    <span></span>
                </a>
                <a href="#" v-appear="{ delay: 1650 }">
                    {{ $t('contactsHeaderNav') }}
                    <span></span>
                </a>
            </nav>
            <div class="header-links">
                <router-link :to="{ name: 'User' }" v-appear="{ delay: 1850 }">
                    <img src="../assets/img/user.png" :alt="$t('altLK')" style="filter: invert(1);">
                </router-link>
                <a href="#" v-appear="{ delay: 1900 }">
                    <img src="../assets/img/cart.png" :alt="$t('altCart')" style="filter: invert(1);">
                </a>
                <button class="change-language"
                    @click="toggleLanguage" v-appear="{ delay: 1950 }">
                    <span :class="{ 'active': selectedLanguage === 'ru' }">RU</span>
                    <span :class="{ 'active': selectedLanguage === 'en' }">EN</span>
                </button>
            </div>
        </div>
    </header>
    <header class="hidden-translate-header" 
            :style="{ transform: heightScrolled > 150 ? 'translateY(0%) translateX(-50%)' : 'translateY(-110%) translateX(-50%)' }">
        <div class="header" v-appear.repeat="{ delay: 0 }">
            <div v-appear.repeat="{ delay: 50 }">
                <router-link :to="{ name: 'Main' }" class="header-logo">
                    <img src="../assets/img/logo.png" :alt="$t('altLogo')" style="filter: invert(1)">
                </router-link>
            </div>
            <nav>
                <router-link :to="{ name: 'Catalog' }" v-appear.repeat="{ delay: 250 }">
                    {{ $t('catalog') }}
                    <span></span>
                </router-link>
                <a href="#" v-appear.repeat="{ delay: 300 }">
                    {{ $t('ordersHeaderNav') }}
                    <span></span>
                </a>
                <a href="#" v-appear.repeat="{ delay: 350 }">
                    {{ $t('favouritesHeaderNav') }}
                    <span></span>
                </a>
                <a href="#" v-appear.repeat="{ delay: 400 }">
                    {{ $t('contactsHeaderNav') }}
                    <span></span>
                </a>
            </nav>
            <div class="header-links">
                <router-link :to="{ name: 'User' }" v-appear.repeat="{ delay: 600 }">
                    <img src="../assets/img/user.png" :alt="$t('altLK')" style="filter: invert(1);">
                </router-link>
                <a href="#" v-appear.repeat="{ delay: 700 }">
                    <img src="../assets/img/cart.png" :alt="$t('altCart')" style="filter: invert(1);">
                </a>
                <button class="change-language" v-appear.repeat="{ delay: 800 }"
                    @click="toggleLanguage">
                    <span :class="{ 'active': selectedLanguage === 'ru' }">RU</span>
                    <span :class="{ 'active': selectedLanguage === 'en' }">EN</span>
                </button>
            </div>
        </div>
    </header>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';

const { locale } = useI18n();

const selectedLanguage = ref('ru');
const heightScrolled = ref(window.scrollY);

const toggleLanguage = () => {
    selectedLanguage.value = selectedLanguage.value === 'ru' ? 'en' : 'ru';
    locale.value = selectedLanguage.value;
};

const updateScroll = () => {
    heightScrolled.value = window.scrollY;
};

onMounted(() => {
    window.addEventListener('scroll', updateScroll);
});

onBeforeUnmount(() => {
    window.removeEventListener('scroll', updateScroll);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/layout/header';
</style>