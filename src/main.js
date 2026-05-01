import App from './App.vue';
import { createApp } from 'vue';
import { createPinia } from 'pinia';
import 'perfect-scrollbar/css/perfect-scrollbar.css';
import router from './router';
import i18n from './i18n';
import appearDirective from './directives/appear';
import recaptchaPlugin from './directives/recaptcha'
import PerfectScrollbarDirective from './directives/perfectScrollbar';

const pinia = createPinia();

createApp(App)
    .use(i18n)
    .use(pinia)
    .directive('appear', appearDirective)
    .directive('perfect-scrollbar', PerfectScrollbarDirective)
    .use(recaptchaPlugin)
    .use(router)
    .mount('#app');