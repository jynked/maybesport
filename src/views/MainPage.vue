<template>
    <main v-if="isPageLoaded">
        <section class="new-item-border" v-appear="{ delay: 600 }">
            <img :src="newItem.image" :alt="$t('newItemImageAlt')">
            <h2>
                <span v-for="(char, index) in animatedTitle" :key="index" v-show="isTitleVisible[index]"
                    :class="{ 'space': char == ' ' }">
                    {{ char == ' ' ? '&nbsp;' : char }}
                </span>
                <span class="type-bar" :style="{ 'opacity': titleTypeBarVisible ? '1' : '0.5' }"></span>
            </h2>
            <router-link :to="{ name: 'Item', params: { itemId: newItem.newItemUniqueId } }" class="decorated-link">
                {{ $t('newItemMore') }}
            </router-link>
        </section>

        <section class="main-page-container" v-appear="{ delay: 800 }">
            <h2>{{ $t('lookForCatalog') }}</h2>
            <div class="catalog">
                <ItemCard v-for="(item, index) in catalogItems" :key="item.uniqueId" v-memo="[item.uniqueId, item.availability, item.minPrice]"
                    v-appear="{ delay: 700 + index * 200 }" :id="item.id" :title="item.title" :images="item.images"
                    :color="item.color" :sizes="item.sizes" :availability="item.availability" :minPrice="item.minPrice"
                    :tags="item.tags" :uniqueId="item.uniqueId" :delay="600 + index * 200"
                    @openSizeModal="openSizesModal" />
            </div>
        </section>

        <section class="main-page-container" v-appear="{ delay: 200 }">
            <h2>{{ $t('ourTGC') }}</h2>
            <div class="tgc-wrapper">
                <img src="../assets/img/telegram-bg.png" :alt="$t('TGCBGAlt')" v-appear="{ delay: 200 }">
                <div class="tgc-container" v-appear="{ delay: 500 }">
                    <p>{{ $t('TGC') }}</p>
                    <a href="https://t.me/MAYBE_SPORT" target="_blank" class="decorated-link">{{ $t('link') }}</a>
                </div>
            </div>
        </section>
        <ItemSizesModal :uniqueId="currentModalUniqueId" :sizes="currentModalSizes" :isOpen="isSizesModalOpen"
            @close="closeSizesModal" @addToCart="handleAddToCart" />
    </main>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import ItemCard from '../components/ItemCard.vue';
import ItemSizesModal from '../components/ItemSizesModal.vue';
import { api } from '../api';

const { t, locale } = useI18n();

const newItem = ref({});
const isPageLoaded = ref(null);
const catalogItems = ref([]);

const isVisible = ref([]);
const typeBarVisible = ref(true);
const typeBarInterval = ref(null);
const animatedText = ref([]);
const originalText = ref('');
const typingTimeouts = ref([]);

const animatedTitle = ref([]);
const isTitleVisible = ref([]);
const titleTypeBarVisible = ref(true);
const titleTypeBarInterval = ref(null);
const titleTypingTimeouts = ref([]);

const isSizesModalOpen = ref(false);
const currentModalUniqueId = ref('');
const currentModalSizes = ref([]);

const emit = defineEmits(['page-loaded']);

watch(() => isPageLoaded.value, (newVal) => {
    if (newVal !== null) emit('page-loaded', newVal)
})

async function loadData() {
    try {
        const newResp = await api.getMainPageNew()
        newItem.value = newResp.data

        const itemsResp = await api.getItems({ page: 1, limit: 4 })
        catalogItems.value = itemsResp.data.items
        isPageLoaded.value = true
        setTimeout(() => {
            initTypingAnimations()
        }, 600)
    } catch (error) {
        console.error('Ошибка загрузки главной:', error)
        isPageLoaded.value = false
    }
}

function clearAllAnimations() {
    if (typeBarInterval.value) {
        clearInterval(typeBarInterval.value);
        typeBarInterval.value = null;
    }
    if (titleTypeBarInterval.value) {
        clearInterval(titleTypeBarInterval.value);
        titleTypeBarInterval.value = null;
    }

    typingTimeouts.value.forEach(timeout => clearTimeout(timeout));
    typingTimeouts.value = [];

    titleTypingTimeouts.value.forEach(timeout => clearTimeout(timeout));
    titleTypingTimeouts.value = [];
}

function initTypingAnimations() {
    clearAllAnimations();
    initLinkAnimation();
    initTitleAnimation();
}

function initLinkAnimation() {
    originalText.value = 'newItemMore';
    animatedText.value = t(originalText.value).split('');
    isVisible.value = Array(animatedText.value.length).fill(false);

    startTypingEffect();
    startTypeBarBlinking();
}

function initTitleAnimation() {
    const titleText = locale.value == 'en' ? newItem.value.title?.en : newItem.value.title?.ru;
    if (!titleText) return;

    animatedTitle.value = titleText.split('');
    isTitleVisible.value = Array(animatedTitle.value.length).fill(false);

    startTitleTypingEffect();
    startTitleTypeBarBlinking();
}

function startTypingEffect() {
    animatedText.value.forEach((_, index) => {
        const timeout = setTimeout(() => {
            isVisible.value[index] = true;

            if (index == animatedText.value.length - 1) {
                const eraseTimeout = setTimeout(() => eraseText(), 1000);
                typingTimeouts.value.push(eraseTimeout);
            }
        }, 75 * index);
        typingTimeouts.value.push(timeout);
    });
}

function startTitleTypingEffect() {
    animatedTitle.value.forEach((_, index) => {
        const timeout = setTimeout(() => {
            isTitleVisible.value[index] = true;

            if (index == animatedTitle.value.length - 1) {
                const eraseTimeout = setTimeout(() => eraseTitleText(), 2000);
                titleTypingTimeouts.value.push(eraseTimeout);
            }
        }, 50 * index);
        titleTypingTimeouts.value.push(timeout);
    });
}

function eraseText() {
    for (let i = animatedText.value.length - 1; i >= 0; i--) {
        const timeout = setTimeout(() => {
            isVisible.value[i] = false;

            if (i == 0) {
                const restartTimeout = setTimeout(() => startTypingEffect(), 500);
                typingTimeouts.value.push(restartTimeout);
            }
        }, 50 * (animatedText.value.length - 1 - i));
        typingTimeouts.value.push(timeout);
    }
}

function eraseTitleText() {
    for (let i = animatedTitle.value.length - 1; i >= 0; i--) {
        const timeout = setTimeout(() => {
            isTitleVisible.value[i] = false;

            if (i == 0) {
                const restartTimeout = setTimeout(() => startTitleTypingEffect(), 1000);
                titleTypingTimeouts.value.push(restartTimeout);
            }
        }, 30 * (animatedTitle.value.length - 1 - i));
        titleTypingTimeouts.value.push(timeout);
    }
}

function startTypeBarBlinking() {
    typeBarInterval.value = setInterval(() => {
        typeBarVisible.value = !typeBarVisible.value;
    }, 400);
}

function startTitleTypeBarBlinking() {
    titleTypeBarInterval.value = setInterval(() => {
        titleTypeBarVisible.value = !titleTypeBarVisible.value;
    }, 400);
}

function openSizesModal({ uniqueId, sizes }) {
    currentModalUniqueId.value = uniqueId;
    currentModalSizes.value = sizes;
    isSizesModalOpen.value = true;
    document.body.style.overflow = 'hidden';
}

function closeSizesModal() {
    isSizesModalOpen.value = false;
    document.body.style.overflow = 'auto';
}

function handleAddToCart({ uniqueId, size }) {
}

watch(locale, async (newLocale) => {
    await nextTick();
    initTypingAnimations();
});

onMounted(() => {
    loadData()
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/main';
</style>