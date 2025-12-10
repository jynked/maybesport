<template>
    <footer>
        <div class="footer">
            <div class="logo">
                <router-link :to="{ name: 'Main' }">
                    <img src="../assets/img/logo.png" :alt="$t('altLogo')" style="filter: invert(1)">
                </router-link>
            </div>
            <div class="link">
                <a href="#" target="_blank">
                    {{ $t('privacyPolicy') }}
                    <span></span>
                </a>
            </div>
            <div class="link">
                <a href="#" target="_blank">
                    {{ $t('userAgreement') }}
                    <span></span>
                </a>
            </div>
        </div>
    </footer>
</template>

<script setup>

</script>

<style lang="scss" scoped>
@use '../assets/styles/layout/footer';
</style>