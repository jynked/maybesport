import { createRouter, createWebHistory } from 'vue-router'
import MainPage from '../views/MainPage.vue';
import CatalogPage from '../views/CatalogPage.vue';
import ItemPage from '../views/ItemPage.vue';
import AdminProducts from '../views/AdminProducts.vue';
import AuthUser from '../views/AuthUser.vue';
import UserPage from '../views/UserPage.vue';
import { useAuthStore } from '../stores/auth';

const routes = [
  {
    path: '/',
    name: 'Main',
    component: MainPage
  },
  {
    path: '/catalog',
    name: 'Catalog',
    component: CatalogPage,
  },
  {
    path: '/catalog/:itemId',
    name: 'Item',
    component: ItemPage,
    props: true,
  },
  {
    path: '/user/auth',
    name: 'UserAuth',
    component: AuthUser,
  },
  {
    path: '/user',
    name: 'User',
    component: UserPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/items',
    name: 'AdminProducts',
    component: AdminProducts,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { top: 0 };
  }
});

let loaderInstance = null;
export const setLoaderInstance = (instance) => { loaderInstance = instance; };

router.beforeEach(async (to, from, next) => {
  if (from.name !== null && loaderInstance) {
    const appearElements = document.querySelectorAll('[v-appear], [data-v-appear]');
    appearElements.forEach(el => {
      el.style.transition = 'all 0.3s ease';
      el.style.opacity = '0';
      el.style.transform = 'translateY(-20px)';
    });
    loaderInstance.startLoading();
    await new Promise(resolve => setTimeout(resolve, 1000));
  }

  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  if (requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'UserAuth', query: { redirect: to.fullPath } });
  } else if (to.name === 'UserAuth' && authStore.isAuthenticated) {
    next({ name: 'User' });
  } else {
    next();
  }
});

router.onError((error) => {
  console.error('Router error:', error);
  if (loaderInstance) loaderInstance.finishLoading();
});

export default router;