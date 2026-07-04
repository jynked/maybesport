<template>
  <main v-if="isPageLoaded">
    <section class="heading-container">
      <div class="heading">
        <h1 v-appear="{ delay: 600 }">
          {{ $i18n.locale == 'en' ? item.title.en : item.title.ru }}
        </h1>
      </div>
      <p v-appear="{ delay: 650 }">
        {{ $i18n.locale == 'en' ? item.category.en : item.category.ru }},
        {{ $i18n.locale == 'en' ? item.sport?.en : item.sport?.ru }}
      </p>
    </section>

    <section class="item-wrapper">
      <div class="item-images-container">
        <button class="item-image" v-appear="{ delay: 700 }"
          :style="{ 'width': item.images.length > 1 ? '75%' : '100%' }" @click="openModal(0)">
          <img :src="item.images[0]" :alt="$t('itemImageAlt')">
        </button>
        <div class="another-images" v-if="item.images.length > 1" v-appear="{ delay: 750 }">
          <button v-for="(image, index) in item.images.slice(1)" :key="index" v-appear="{ delay: 50 * (index % 4) }"
            @click="openModal(index + 1)">
            <img :src="image" :alt="$t('itemImageAlt')">
          </button>
        </div>
      </div>

      <span v-appear="{ delay: 600 }"></span>

      <div class="item-wrapper-right-block">
        <div class="action-item-buttons">
          <button @click.stop="toggleFavourite" v-appear="{ delay: 400 }">
            <svg class="favourite-icon" :class="{ 'favourite-active': isFav }" viewBox="0 0 24 24" width="24"
              height="24" fill="none" stroke="currentColor" stroke-width="2">
              <path
                d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
            </svg>
          </button>
          <button @click="openMainItemSizesModal" v-appear="{ delay: 500 }">
            <div class="item-badge" v-if="cartStore.getTotalQuantityByUniqueId(currentUniqueId)">
              {{ cartStore.getTotalQuantityByUniqueId(currentUniqueId) > 99 ? '99+' :
                cartStore.getTotalQuantityByUniqueId(currentUniqueId) }}
            </div>
            <img src="../assets/img/characteristics.png" :alt="$t('characteristicsAlt')">
          </button>
          <button @click="openCharacteristicsModal" v-appear="{ delay: 600 }">
            <img src="../assets/img/info.png" :alt="$t('info')">
          </button>
        </div>
        <p class="structure-item-in-page" v-appear="{ delay: 800 }" v-if="item.description.en || item.description.ru">
          <strong>{{ $t('description') }}:</strong>
          {{ $i18n.locale == 'en' ? item.description.en : item.description.ru }}
        </p>
        <p v-appear="{ delay: 850 }" v-if="item.siblingItems && item.siblingItems.length > 0">{{ $t(randomSubtitleKey)
        }}</p>
        <swiper :free-mode="true" :slides-per-view="2.4" direction="horizontal" space-between="5" class="cards-subitems"
          v-if="item.siblingItems && item.siblingItems.length > 0" v-appear="{ delay: 900 }">
          <swiper-slide v-for="(sibling, index) in item.siblingItems" :key="sibling.uniqueId">
            <router-link :to="{ name: 'Item', params: { itemId: sibling.uniqueId } }" class="card-subitem">
              <img :src="sibling.image" :alt="$t('subItem') + ' ' + (index + 1)">
            </router-link>
          </swiper-slide>
        </swiper>
      </div>
    </section>

    <section class="simple-offers-wrapper">
      <h4 v-appear="{ delay: 500 }">{{ $t('similarOffers') }}</h4>
      <div class="similar-products">
        <ItemCard v-for="(similarItem, index) in similarProducts" :key="similarItem.uniqueId"
          v-memo="[similarItem.uniqueId, similarItem.availability, similarItem.minPrice]" :id="similarItem.id"
          :title="similarItem.title" :images="similarItem.images" :color="similarItem.color" :sizes="similarItem.sizes"
          :availability="similarItem.availability" :minPrice="similarItem.minPrice" :tags="similarItem.tags"
          :uniqueId="similarItem.uniqueId"
          :delay="windowWidth > 1440 ? 600 + index * 200 : windowWidth > 480 ? (600 + index % 2 * 200) : 600"
          v-appear="{ delay: windowWidth > 1440 ? 600 + index * 200 : (600 + index % 2 * 200) }"
          @openSizeModal="openSimilarItemSizesModal" />
      </div>
    </section>

    <Transition name="modal">
      <div v-if="isModalOpen" class="item-photos-modal" @click="closeModal">
        <div class="modal-content">
          <div class="view-photo">
            <img :src="currentImage" :alt="$t('mainImage')">
          </div>
          <div class="pagination-cards" @click.stop>
            <swiper :free-mode="true" :slides-per-view="'auto'" direction="horizontal" space-between="5"
              :initial-slide="startFromFirst ? 0 : currentIndex" @click.stop>
              <swiper-slide v-for="(img, index) in item.images" :key="index">
                <button :class="{ active: currentIndex === index }" @click="setCurrentImage(index)">
                  <img :src="img" :alt="$t('itemImageAlt') + ' ' + index">
                </button>
              </swiper-slide>
            </swiper>
          </div>
        </div>
      </div>
    </Transition>

    <Transition name="modal">
      <div v-if="isCharacteristicsModalOpen" class="characteristics-modal" @click="closeCharacteristicsModal">
        <div class="modal-content-characteristics" @click.stop>
          <div class="modal-header">
            <h2>{{ $t('characteristics').toUpperCase() }}</h2>
            <button class="close-button" @click="closeCharacteristicsModal">×</button>
          </div>
          <div class="characteristics-content">
            <section class="characteristics-container">
              <p v-if="item.tags && item.tags.length > 0" class="tags">
                <span v-for="tag in item.tags" :key="tag.en" class="tag">
                  {{ $i18n.locale == 'en' ? tag.en : tag.ru }}
                </span>
              </p>
              <div class="subitem-characteristics">
                <div class="characteristic-group">
                  <p class="structure-item">
                    <strong>{{ $t('sport') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.sport?.en : item.sport?.ru }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('category') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.category.en : item.category.ru }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('color') }}: </strong>
                    <span v-for="(col, idx) in item.color" :key="idx">
                      {{ $i18n.locale == 'en' ? col.en : col.ru }}{{ idx < item.color.length - 1 ? ', ' : '' }} </span>
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('type') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.type.en : item.type.ru }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('brand') }}:</strong>
                    {{ item.brand }}
                  </p>
                  <div class="structure">
                    <div class="structure-items">
                      <span v-for="material in item.structure" :key="material.name.en" class="structure-item">
                        {{ $i18n.locale == 'en' ? material.name.en : material.name.ru }}: {{ material.percent }}%
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </section>
          </div>
        </div>
      </div>
    </Transition>
    <ItemSizesModal :uniqueId="currentModalUniqueId" :sizes="currentModalSizes" :isOpen="isSizesModalOpen"
      @close="closeSizesModal" @addToCart="handleAddToCart" />
  </main>
</template>

<script setup>
import { Swiper, SwiperSlide } from 'swiper/vue';
import 'swiper/css';
import { ref, watch, computed, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ItemCard from '../components/ItemCard.vue';
import ItemSizesModal from '../components/ItemSizesModal.vue';
import { api } from '../api';
import { useFavouritesStore } from '../stores/favourites';
import { useCartStore } from '../stores/cart';
import { useToastStore } from '../stores/toast';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '../stores/auth';

const cartStore = useCartStore();
const favouritesStore = useFavouritesStore();

const emit = defineEmits(['page-loaded']);

const props = defineProps({ itemId: String });

const { t, locale } = useI18n();


const route = useRoute();
const router = useRouter();

const item = ref(null);
const similarProducts = ref([]);
const isPageLoaded = ref(false);
let loadingPromise = null;

const isSizesModalOpen = ref(false);
const currentModalUniqueId = ref('');
const currentModalSizes = ref([]);

const isModalOpen = ref(false);
const currentIndex = ref(0);
const startFromFirst = ref(false);
const currentImage = ref('');
const isCharacteristicsModalOpen = ref(false);
const randomIndex = ref(0);

const windowWidth = ref(window.innerWidth);

const currentUniqueId = computed(() => {
  return item.value?.uniqueId || route.params.itemId;
});

const randomSubtitleKey = computed(() => `subItems${randomIndex.value}`);

const isFav = computed(() => favouritesStore.isFavourite(currentUniqueId.value));
const totalInCart = computed(() => cartStore.getTotalQuantityByUniqueId(currentUniqueId.value));

async function loadItemData(itemId) {
  if (loadingPromise) return loadingPromise;

  isPageLoaded.value = false;

  loadingPromise = (async () => {
    try {
      const [itemResp, similarResp] = await Promise.all([
        api.getItem(itemId),
        api.getSimilar(itemId, 4)
      ]);

      item.value = itemResp.data;
      similarProducts.value = similarResp.data;
      isPageLoaded.value = true;
      emit('page-loaded', true);
    } catch (error) {
      const status = error.response?.status;
      if (status === 404) {
        router.push({ name: 'Error' });
      } else {
        console.error(error);
      }
      isPageLoaded.value = false;
      emit('page-loaded', false);
    } finally {
      loadingPromise = null;
    }
  })();

  return loadingPromise;
}

function openModal(index, fromShowMore = false) {
  currentIndex.value = fromShowMore ? 0 : index;
  startFromFirst.value = fromShowMore;
  currentImage.value = item.value.images[fromShowMore ? 0 : index];
  isModalOpen.value = true;
}

function closeModal() {
  isModalOpen.value = false;
}

function setCurrentImage(index) {
  currentIndex.value = index;
  currentImage.value = item.value.images[index];
}

async function toggleFavourite() {
  if (!useAuthStore().isAuthenticated) {
    localStorage.setItem('pendingAction', JSON.stringify({
      action: 'favourite',
      uniqueId: currentUniqueId.value
    }));
    router.push({ name: 'UserAuth', query: { redirect: router.currentRoute.value.fullPath } });
    return;
  }
  if (isFav.value) {
    await favouritesStore.removeFromFavourites(currentUniqueId.value);
    useToastStore().success(t('removedFromFavourites'));
  } else {
    await favouritesStore.addToFavourites(currentUniqueId.value);
    useToastStore().success(t('addedToFavourites'));
  }
}

function openMainItemSizesModal() {
  currentModalUniqueId.value = currentUniqueId.value;
  currentModalSizes.value = item.value?.sizes || [];
  isSizesModalOpen.value = true;
  document.body.style.overflow = 'hidden';
}

function openSimilarItemSizesModal({ uniqueId, sizes }) {
  currentModalUniqueId.value = uniqueId;
  currentModalSizes.value = sizes;
  isSizesModalOpen.value = true;
  document.body.style.overflow = 'hidden';
}

function closeSizesModal() {
  isSizesModalOpen.value = false;
  document.body.style.overflow = 'auto';
  currentModalUniqueId.value = '';
  currentModalSizes.value = [];
}

function handleAddToCart({ uniqueId, size }) {

}

function calcWidth() {
  windowWidth.value = window.innerWidth;
}

function openCharacteristicsModal() {
  isCharacteristicsModalOpen.value = true;
  document.body.style.overflow = 'hidden';
}

function closeCharacteristicsModal() {
  isCharacteristicsModalOpen.value = false;
  document.body.style.overflow = 'auto';
}

function redirectToAuthWithAction(action, uniqueId, size) {
  const pending = { action, uniqueId, size };
  localStorage.setItem('pendingAction', JSON.stringify(pending));
  router.push({ name: 'UserAuth', query: { redirect: '/favourites' } });
}

watch(() => route.params.itemId, async (newId) => {
  if (newId) {
    await loadItemData(newId);
    randomIndex.value = Math.floor(Math.random() * 5);
  };
}, { immediate: true });

onMounted(() => {
  window.addEventListener('resize', calcWidth);
})

onUnmounted(() => {
  window.removeEventListener('resize', calcWidth);
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/item';
</style>