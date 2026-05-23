<template>
    <header v-appear="{ delay: 1200 }">
        <div class="header">
            <div v-appear="{ delay: 1300 }">
                <router-link :to="{ name: 'Main' }" class="header-logo">
                    <img src="../assets/img/logo.png" :alt="$t('altLogo')" style="filter: invert(1)">
                </router-link>
            </div>

            <nav class="desktop-nav">
                <router-link :to="{ name: 'Catalog' }" v-appear="{ delay: 1500 }">
                    {{ $t('catalog') }}<span></span>
                </router-link>
                <router-link :to="{ name: 'UserOrders' }" v-appear="{ delay: 1550 }">
                    {{ $t('ordersHeaderNav') }}<span></span>
                </router-link>
                <router-link :to="{ name: 'Favourites' }" v-appear="{ delay: 1600 }">
                    {{ $t('favouritesHeaderNav') }}<span></span>
                </router-link>
                <a href="#" v-appear="{ delay: 1650 }">
                    {{ $t('contactsHeaderNav') }}<span></span>
                </a>
            </nav>

            <button class="burger-btn" @click="toggleMobileMenu" :class="{ 'active': isMobileMenuOpen }" v-appear="{ delay: 1600 }">
                <span></span><span></span><span></span>
            </button>

            <div class="header-links">
                <router-link :to="{ name: 'User' }" v-appear="{ delay: 1850 }">
                    <img src="../assets/img/user.png" :alt="$t('altLK')" style="filter: invert(1);">
                </router-link>
                <button @click.prevent="openCartModal" class="cart-link" v-appear="{ delay: 1900 }">
                    <img src="../assets/img/cart.png" :alt="$t('altCart')" style="filter: invert(1);">
                    <span v-if="cartStore.totalQuantity" class="cart-badge">{{ cartStore.totalQuantity }}</span>
                </button>
                <button class="change-language" @click="toggleLanguage" v-appear="{ delay: windowWidth > 1024 ? 1950 : 1500 }">
                    <span :class="{ 'active': selectedLanguage === 'ru' }">RU</span>
                    <span :class="{ 'active': selectedLanguage === 'en' }">EN</span>
                </button>
            </div>
        </div>

    </header>
    
    <transition name="mobile-menu-fade">
        <div class="mobile-menu" v-if="isMobileMenuOpen" @click.self="closeMobileMenu">
            <div class="mobile-menu-inner">
                <nav class="mobile-nav">
                    <router-link :to="{ name: 'Catalog' }" @click="closeMobileMenu">
                        {{ $t('catalog') }}
                    </router-link>
                    <router-link :to="{ name: 'UserOrders' }" @click="closeMobileMenu">
                        {{ $t('ordersHeaderNav') }}
                    </router-link>
                    <router-link :to="{ name: 'Favourites' }" @click="closeMobileMenu">
                        {{ $t('favouritesHeaderNav') }}
                    </router-link>
                    <a href="#" @click="closeMobileMenu">
                        {{ $t('contactsHeaderNav') }}
                    </a>
                    <router-link :to="{ name: 'User' }">
                        {{ $t('altLK') }}
                    </router-link>
                </nav>
                <div class="menu-action-links header-links">
                    <button @click.prevent="openCartModal" class="cart-link">
                        <img src="../assets/img/cart.png" :alt="$t('altCart')" style="filter: invert(1);">
                        <span v-if="cartStore.totalQuantity" class="cart-badge">{{ cartStore.totalQuantity }}</span>
                    </button>
                </div>
            </div>
        </div>
    </transition>
    <header class="hidden-translate-header"
        :style="{ transform: heightScrolled > 150 ? 'translateY(0%) translateX(-50%)' : 'translateY(-110%) translateX(-50%)' }">
        <div class="header" v-appear.repeat="{ delay: 0 }">
            <div v-appear.repeat="{ delay: 150 }">
                <router-link :to="{ name: 'Main' }" class="header-logo">
                    <img src="../assets/img/logo.png" :alt="$t('altLogo')" style="filter: invert(1)">
                </router-link>
            </div>

            <nav class="desktop-nav">
                <router-link :to="{ name: 'Catalog' }" v-appear.repeat="{ delay: 250 }">
                    {{ $t('catalog') }}<span></span>
                </router-link>
                <router-link :to="{ name: 'UserOrders' }" v-appear.repeat="{ delay: 300 }">
                    {{ $t('ordersHeaderNav') }}<span></span>
                </router-link>
                <router-link :to="{ name: 'Favourites' }" v-appear.repeat="{ delay: 350 }">
                    {{ $t('favouritesHeaderNav') }}<span></span>
                </router-link>
                <a href="#" v-appear.repeat="{ delay: 400 }">
                    {{ $t('contactsHeaderNav') }}<span></span>
                </a>
            </nav>

            <button class="burger-btn" @click="toggleMobileMenu" :class="{ 'active': isMobileMenuOpen }" v-appear="{ delay: 500 }">
                <span></span><span></span><span></span>
            </button>

            <div class="header-links">
                <router-link :to="{ name: 'User' }" v-appear.repeat="{ delay: 600 }">
                    <img src="../assets/img/user.png" :alt="$t('altLK')" style="filter: invert(1);">
                </router-link>
                <button @click.prevent="openCartModal" class="cart-link" v-appear.repeat="{ delay: 700 }">
                    <img src="../assets/img/cart.png" :alt="$t('altCart')" style="filter: invert(1);">
                    <span v-if="cartStore.totalQuantity" class="cart-badge">{{ cartStore.totalQuantity }}</span>
                </button>
                <button class="change-language" v-appear.repeat="{ delay: windowWidth > 1024 ? 800 : 200 }" @click="toggleLanguage">
                    <span :class="{ 'active': selectedLanguage === 'ru' }">RU</span>
                    <span :class="{ 'active': selectedLanguage === 'en' }">EN</span>
                </button>
            </div>
        </div>
    </header>

    <CartModal :isOpen="isCartModalOpen" @close="closeCartModal" />
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useCartStore } from '../stores/cart';
import CartModal from './CartModal.vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const router = useRouter();
router.afterEach(() => {
    closeMobileMenu();
});

const { locale } = useI18n();
const selectedLanguage = ref('ru');
const heightScrolled = ref(window.scrollY);
const isMobileMenuOpen = ref(false);
const windowWidth = ref(window.innerWidth);

const cartStore = useCartStore();
const authStore = useAuthStore();
const isCartModalOpen = ref(false);
const openCartModal = () => { isCartModalOpen.value = true; };
const closeCartModal = () => { isCartModalOpen.value = false; };

const toggleLanguage = () => {
    selectedLanguage.value = selectedLanguage.value === 'ru' ? 'en' : 'ru';
    locale.value = selectedLanguage.value;
};

const toggleMobileMenu = () => {
    isMobileMenuOpen.value = !isMobileMenuOpen.value;
    if (isMobileMenuOpen.value) {
        document.body.style.overflow = 'hidden';
    } else {
        document.body.style.overflow = '';
    }
};

const closeMobileMenu = () => {
    if (isMobileMenuOpen.value) {
        isMobileMenuOpen.value = false;
        document.body.style.overflow = '';
    }
};

const updateScroll = () => {
    heightScrolled.value = window.scrollY;
};

function calcWidth() {
    windowWidth.value = window.innerWidth;
    if (windowWidth.value > 1024) {
        closeMobileMenu();
    }
}

onMounted(() => {
    window.addEventListener('scroll', updateScroll);
    if (authStore.isAuthenticated) {
        cartStore.loadCart();
    }
    window.addEventListener('resize', calcWidth);
});

onBeforeUnmount(() => {
    window.removeEventListener('scroll', updateScroll);
    document.body.style.overflow = '';
    window.removeEventListener('resize', calcWidth);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/layout/header';
</style>