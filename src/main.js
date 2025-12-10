import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import i18n from './i18n';
import appearDirective from './directives/appear';

createApp(App).use(i18n).directive('appear', appearDirective).use(router).mount('#app')