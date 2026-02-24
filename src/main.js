import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import i18n from './i18n';
import appearDirective from './directives/appear';

const pinia = createPinia();

createApp(App)
  .use(i18n)
  .use(pinia)
  .directive('appear', appearDirective)
  .use(router)
  .mount('#app');