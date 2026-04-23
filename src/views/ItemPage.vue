<template>
  <main v-if="isPageLoaded">
    <section class="heading-container">
      <div class="heading">
        <h1 v-appear="{ delay: 600 }">
          {{ $i18n.locale == 'en' ? item.title.en : item.title.ru }}
        </h1>
        <div class="action-item-buttons">
          <button v-appear="{ delay: 800 }" @click="openMainItemSizesModal">
            <img src="../assets/img/characteristics.png" :alt="$t('characteristicsAlt')">
          </button>
          <button v-appear="{ delay: 900 }" @click="openCharacteristicsModal">
            <img src="../assets/img/info.png" :alt="$t('info')">
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
        <ItemCard v-for="(similarItem, index) in similarProducts" :key="similarItem.uniqueId" v-memo="[similarItem.uniqueId, similarItem.availability, similarItem.minPrice]"
          :id="similarItem.id"
          :title="similarItem.title" :images="similarItem.images" :color="similarItem.color" :sizes="similarItem.sizes"
          :availability="similarItem.availability" :minPrice="similarItem.minPrice" :tags="similarItem.tags"
          :uniqueId="similarItem.uniqueId" :delay="200 + index * 100" v-appear="{ delay: 200 + index * 100 }"
          @openSizeModal="openSimilarItemSizesModal" />
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
import { ref, onMounted, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ItemCard from '../components/ItemCard.vue';
import ItemSizesModal from '../components/ItemSizesModal.vue';
import { api } from '../api';
import { useFavouritesStore } from '../stores/favourites';
import { useAuthStore } from '../stores/auth';
import { useCartStore } from '../stores/cart';

const cartStore = useCartStore();
const favouritesStore = useFavouritesStore();

const emit = defineEmits(['page-loaded']);

const props = defineProps({ itemId: String });

const route = useRoute();
const router = useRouter();

const item = ref(null);
const similarProducts = ref([]);
const isPageLoaded = ref(false);

const isSizesModalOpen = ref(false);
const currentModalUniqueId = ref('');
const currentModalSizes = ref([]);

const isModalOpen = ref(false);
const currentIndex = ref(0);
const startFromFirst = ref(false);
const currentImage = ref('');
const isCharacteristicsModalOpen = ref(false);

const currentUniqueId = computed(() => {
  return item.value?.uniqueId || route.params.itemId;
});

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