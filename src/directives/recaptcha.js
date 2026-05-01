import { VueReCaptcha } from 'vue-recaptcha-v3'

export default {
  install(app) {
    app.use(VueReCaptcha, {
      siteKey: import.meta.env.VITE_RECAPTCHA_SITE_KEY,
      loaderOptions: {
        autoHideBadge: true,
      },
    })
  },
}