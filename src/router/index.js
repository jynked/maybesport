import { createRouter, createWebHistory } from 'vue-router'
import MainPage from '../views/MainPage.vue';
import CatalogPage from '../views/CatalogPage.vue';
import ItemPage from '../views/ItemPage.vue';
import AdminProducts from '../views/AdminProducts.vue';

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
    path: '/admin/items',
    name: 'AdminProducts',
    component: AdminProducts,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  }
});

let loaderInstance = null;

export const setLoaderInstance = (instance) => {
  loaderInstance = instance;
};

router.beforeEach((to, from, next) => {
  if (from.name === null) {
    next();
    return;
  }
  
  if (loaderInstance) {
    const appearElements = document.querySelectorAll('[v-appear], [data-v-appear]');
    appearElements.forEach(el => {
      el.style.transition = 'all 0.3s ease';
      el.style.opacity = '0';
      el.style.transform = 'translateY(-20px)';
    });
    
    loaderInstance.startLoading();
    setTimeout(() => {
      next();
    }, 1000);
  } else {
    next();
  }
});


router.afterEach(() => {
});

router.onError((error) => {
  console.error('Router error:', error);
  if (loaderInstance) {
    loaderInstance.finishLoading();
  }
});

export default router;
