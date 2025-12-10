<template>
    <header v-appear="{ delay: 1200 }">
        <div class="header">
            <div v-appear="{ delay: 1400 }">
                <router-link :to="{ name: 'Main' }" class="header-logo">
                    <img src="../assets/img/logo.png" :alt="$t('altLogo')" style="filter: invert(1)">
                </router-link>
            </div>
            <nav>
                <a href="#" v-appear="{ delay: 1500 }">
                    {{ $t('catalogHeaderNav') }}
                    <span></span>
                </a>
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
                <a href="#" v-appear="{ delay: 1850 }">
                    <img src="../assets/img/user.png" :alt="$t('altLK')" style="filter: invert(1);">
                </a>
                <a href="#" v-appear="{ delay: 1900 }">
                    <img src="../assets/img/cart.png" :alt="$t('altCart')" style="filter: invert(1);">
                </a>
                <button class="change-language"
                    @click="selectedLanguage == 'ru' ? selectedLanguage = 'en' : selectedLanguage = 'ru'; changeLanguage()" v-appear="{ delay: 1950 }">
                    <span :class="{ 'active': selectedLanguage == 'ru' }">RU</span>
                    <span :class="{ 'active': selectedLanguage == 'en' }">EN</span>
                </button>
            </div>
        </div>
    </header>
    <header class="hidden-translate-header" :style="{ 'transform': heightScrolled > 150 ? 'translateY(0%) translateX(-50%) !important' : 'translateY(-110%) translateX(-50%) !important' }">
        <div class="header" v-appear.repeat="{ delay: 0 }">
            <div v-appear.repeat="{ delay: 50 }">
                <router-link :to="{ name: 'Main' }" class="header-logo">
                    <img src="../assets/img/logo.png" :alt="$t('altLogo')" style="filter: invert(1)">
                </router-link>
            </div>
            <nav>
                <a href="#" v-appear.repeat="{ delay: 250 }">
                    {{ $t('catalogHeaderNav') }}
                    <span></span>
                </a>
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
                <a href="#" v-appear.repeat="{ delay: 600 }">
                    <img src="../assets/img/user.png" :alt="$t('altLK')" style="filter: invert(1);">
                </a>
                <a href="#" v-appear.repeat="{ delay: 700 }">
                    <img src="../assets/img/cart.png" :alt="$t('altCart')" style="filter: invert(1);">
                </a>
                <button class="change-language" v-appear.repeat="{ delay: 800 }"
                    @click="selectedLanguage == 'ru' ? selectedLanguage = 'en' : selectedLanguage = 'ru'; changeLanguage()">
                    <span :class="{ 'active': selectedLanguage == 'ru' }">RU</span>
                    <span :class="{ 'active': selectedLanguage == 'en' }">EN</span>
                </button>
            </div>
        </div>
    </header>
</template>

<script>

export default {
    data() {
        return {
            selectedLanguage: 'ru',
            heightScrolled: window.scrollY,
        }
    },
    methods: {
        changeLanguage() {
            this.$i18n.locale = this.selectedLanguage;
        },
        calcScroll() {
            this.heightScrolled = window.scrollY;
        }
    },
    mounted() {
        window.addEventListener('scroll', this.calcScroll);
    },
    beforeUnmount() {
        window.removeEventListener('scroll', this.calcScroll);
    }
}


</script>

<style lang="scss" scoped>
@use '../assets/styles/layout/header';
</style>