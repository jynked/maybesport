<template>
  <HeaderMain/>
  <router-view/>
  <FooterMain/>
</template>

<script>
import HeaderMain from './components/HeaderMain.vue'
import FooterMain from './components/FooterMain.vue'

export default {
  components: {
    HeaderMain,
    FooterMain,
  }
}

</script>

<style>
@import url('./assets/styles/style.scss');
</style>