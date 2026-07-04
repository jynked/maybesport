<template>
    <footer v-appear="{ delay: 1200 }">
        <div class="footer">
            <div class="link" v-appear="{ delay: 1600 }">
                <a href="#" target="_blank">
                    {{ $t('privacyPolicy') }}
                    <span></span>
                </a>
            </div>
            <div class="link" v-appear="{ delay: 1800 }">
                <a href="#" target="_blank">
                    {{ $t('userAgreement') }}
                    <span></span>
                </a>
            </div>
            <div class="captcha-message">
                {{ $t('captchaMessage') }}
            </div>
        </div>
    </footer>
</template>

<script setup>

</script>

<style lang="scss" scoped>
@use '../assets/styles/layout/footer';
</style>