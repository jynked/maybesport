<template>
  <LoaderVue ref="loader" :isPageLoaded="isAppLoaded"/>
  <HeaderMain/>
  <router-view @page-loaded="handlePageLoaded"/>
  <FooterMain/>
</template>

<script>
import HeaderMain from './components/HeaderMain.vue'
import FooterMain from './components/FooterMain.vue'
import LoaderVue from './components/LoaderVue.vue'
import { setLoaderInstance } from './router'

export default {
  components: {
    HeaderMain,
    FooterMain,
    LoaderVue
  },
  data() {
    return {
      isAppLoaded: false
    }
  },
  mounted() {
    setLoaderInstance(this.$refs.loader);
    
    setTimeout(() => {
      this.isAppLoaded = true;
    }, 1000);
  },
  methods: {
    handlePageLoaded(loaded) {
      if (this.$refs.loader) {
        if (loaded) {
          this.$refs.loader.finishLoading();
        } else {
          alert('Упс! Произошла ошибка! Повторите попытку немного позднее');
          console.log(this.$refs.loader);
        }
      }
    }
  }
}
</script>

<style>
@import url('./assets/styles/style.scss');
</style>