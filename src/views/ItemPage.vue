<template>
  <main v-if="isPageLoaded">
    <section class="heading-container">
      <div class="heading">
        <h1 v-appear="{ delay: 600 }">
          {{ $i18n.locale == 'en' ? item.title.en : item.title.ru }}
        </h1>
        <div class="action-item-buttons">
          <button v-appear="{ delay: 900 }" @click="openCharacteristicsModal">
            <img src="../assets/img/characteristics.png" :alt="$t('characteristicsAlt')">
          </button>
        </div>
      </div>
      <p v-appear="{ delay: 650 }">
        {{ $i18n.locale == 'en' ? item.category.en : item.category.ru }},
        {{ $i18n.locale == 'en' ? item.country.en : item.country.ru }}
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

      <div class="cards-subitems" v-if="item.siblingItems && item.siblingItems.length > 0">
        <router-link v-for="(sibling, index) in item.siblingItems" :key="sibling.uniqueId"
          :to="{ name: 'Item', params: { itemId: sibling.uniqueId } }" class="card-subitem"
          v-appear="{ delay: 200 + 50 * index }">
          <img :src="sibling.image" alt="">
        </router-link>
      </div>
    </section>

    <section class="simple-offers-wrapper">
      <h4>{{ $t('similarOffers') }}</h4>
      <div class="similar-products">
        <ItemCard v-for="(similarItem, index) in similarProducts" :key="similarItem.uniqueId" :id="similarItem.id"
          :title="similarItem.title" :images="similarItem.images" :color="similarItem.color" :sizes="similarItem.sizes"
          :availability="similarItem.availability" :minPrice="similarItem.minPrice" :tags="similarItem.tags"
          :uniqueId="similarItem.uniqueId" :delay="200 + index * 100" v-appear="{ delay: 200 + index * 100 }" />
      </div>
    </section>

    <Transition name="modal">
      <div v-if="isModalOpen" class="item-photos-modal" @click="closeModal">
        <div class="modal-content">
          <div class="view-photo">
            <img :src="currentImage" alt="">
          </div>
          <div class="pagination-cards" @click.stop>
            <swiper :free-mode="true" :slides-per-view="'auto'" direction="horizontal" space-between="5"
              :initial-slide="startFromFirst ? 0 : currentIndex" @click.stop>
              <swiper-slide v-for="(img, index) in item.images" :key="index">
                <button :class="{ active: currentIndex === index }" @click="setCurrentImage(index)">
                  <img :src="img" alt="">
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
                    <strong>{{ $t('category') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.category.en : item.category.ru }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('color') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.color[0]?.en : item.color[0]?.ru }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('type') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.type.en : item.type.ru }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('brand') }}:</strong>
                    {{ item.brand }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('country') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.country.en : item.country.ru }}
                  </p>
                  <p class="structure-item">
                    <strong>{{ $t('description') }}:</strong>
                    {{ $i18n.locale == 'en' ? item.description.en : item.description.ru }}
                  </p>
                  <div class="structure">
                    <div class="structure-items">
                      <span v-for="material in item.structure" :key="material.name.en" class="structure-item">
                        {{ $i18n.locale == 'en' ? material.name.en : material.name.ru }}: {{ material.percent }}%
                      </span>
                    </div>
                  </div>

                  <div class="size-items">
                    <div v-for="sizeItem in item.sizes" :key="sizeItem.size"
                      :class="['size-item', getSizeStatusClass(sizeItem), { 'selected': selectedSize === sizeItem.size }]"
                      @click="selectSize(sizeItem.size)">
                      <span class="size">{{ sizeItem.size }}</span>
                      <span class="price">{{ sizeItem.price.toLocaleString() }} ₽</span>
                      <span class="status">{{ $t(getAvailabilityStatus([sizeItem])) }}</span>
                      <span class="quantity" v-if="sizeItem.quantity > 0 && !sizeItem.isOnRequest">
                        ({{ sizeItem.quantity }} {{ $t('pieces') }})
                      </span>
                      <span class="item-actions">
                        <button @click.stop="toggleFavourite(sizeItem.size)" :disabled="togglingSize === sizeItem.size">
                          <img src="../assets/img/favourite.png" :alt="$t('favouriteAlt')"
                            :style="{ 'filter': favouritesStore.isFavourite(currentUniqueId, sizeItem.size) ? '' : 'sepia(1)' }">
                        </button>
                        <button class="item-cart">
                          <img src="../assets/img/cart.png" :alt="$t('cartAlt')">
                        </button>
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
  </main>
</template>

<script setup>
import { Swiper, SwiperSlide } from 'swiper/vue';
import 'swiper/css';
import { ref, onMounted, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ItemCard from '../components/ItemCard.vue';
import { api } from '../api';
import { useFavouritesStore } from '../stores/favourites';
import { useAuthStore } from '../stores/auth';

const favouritesStore = useFavouritesStore();

const emit = defineEmits(['page-loaded']);

const props = defineProps({ itemId: String });

const route = useRoute();
const router = useRouter();

const item = ref(null);
const similarProducts = ref([]);
const isPageLoaded = ref(false);

const isModalOpen = ref(false);
const currentIndex = ref(0);
const startFromFirst = ref(false);
const currentImage = ref('');
const isCharacteristicsModalOpen = ref(false);

const selectedSize = ref(null);
const togglingSize = ref(null);

const currentUniqueId = computed(() => {
  return item.value?.uniqueId || route.params.itemId;
});

function selectSize(size) {
  selectedSize.value = size;
}

function getSizeStatusClass(sizeItem) {
  if (sizeItem.quantity > 0 && !sizeItem.isOnRequest) return 'available';
  if (sizeItem.quantity > 0 && sizeItem.isOnRequest) return 'on-request';
  return 'out-of-stock';
}

function getAvailabilityStatus(sizes) {
  const available = sizes.some(s => s.quantity > 0 && !s.isOnRequest);
  const onRequest = sizes.some(s => s.quantity > 0 && s.isOnRequest);
  if (available) return 'available';
  if (onRequest) return 'on_request';
  return 'out_of_stock';
}

async function loadItemData(itemId) {
  isPageLoaded.value = false;
  try {
    const itemResp = await api.getItem(itemId);
    item.value = itemResp.data;

    const similarResp = await api.getSimilar(itemId, 4);
    similarProducts.value = similarResp.data;

    isPageLoaded.value = true;
    emit('page-loaded', true);
  } catch (error) {
    console.error('Ошибка загрузки товара:', error);
    isPageLoaded.value = false;
    emit('page-loaded', false);
  }
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

async function toggleFavourite(size) {
  const id = currentUniqueId.value;
  if (!id || !size) return;
  if (togglingSize.value === size) return;

  const authStore = useAuthStore();
  if (!authStore.isAuthenticated) {
    redirectToAuthWithAction('favourite', id, size);
    return;
  }

  togglingSize.value = size;
  try {
    if (favouritesStore.isFavourite(id, size)) {
      await favouritesStore.removeFromFavourites(id, size);
    } else {
      await favouritesStore.addToFavourites(id, size);
    }
  } catch (error) {
    console.error('Error toggling favourite:', error);
    if (error.response?.status === 401) {
      redirectToAuthWithAction('favourite', id, size);
    }
  } finally {
    togglingSize.value = null;
  }
}

watch(() => route.params.itemId, async (newId) => {
  if (newId) await loadItemData(newId);
}, { immediate: true });

onMounted(async () => {
  const id = route.params.itemId || props.itemId;
  if (id) await loadItemData(id);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/item';
</style>