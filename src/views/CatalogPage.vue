<template>
  <main v-if="isPageLoaded">
    <div class="catalog">
      <div class="catalog-block">
        <h1 v-appear="{ delay: 400 }">{{ $t('catalog') }}</h1>
        <div class="filters-main-buttons-block">
          <div class="search-container" v-appear="{ delay: 500 }" :class="{ 'expanded': isSearchExpanded }"
            ref="inputBox">
            <button @click="expandSearch">
              <img src="../assets/img/loop.png" :alt="$t('LoopAlt')" />
            </button>
            <input type="text" name="search" id="search" :placeholder="$t('search')" @focus="expandSearch"
              ref="searchInput" @blur="handleSearchBlur" v-model="searchQuery">
          </div>

          <div class="dropdown-wrapper" v-appear="{ delay: 650 }">
            <button class="dropdown-trigger" @click="toggleDropdown('sort', 'main')"
              :class="{ 'dropdown-open': isMainSortDropdownOpen }">
              {{ $t('sortBy') }}
              <img src="../assets/img/down.png" :alt="$t('Down')" :class="{ 'rotated': isMainSortDropdownOpen }">
            </button>

            <transition name="dropdown">
              <div v-if="isMainSortDropdownOpen" class="dropdown-content" @click="closeDropdownOnClickOutside('main')">
                <div class="dropdown-group">
                  <h3>{{ $t('alphabet') }}</h3>
                  <button @click="sortItems('alphabet-asc')" :class="{ 'active': currentSort == 'alphabet-asc' }">
                    {{ $t('FirstLetter') }} → {{ $t('LastLetter') }}
                  </button>
                  <button @click="sortItems('alphabet-desc')" :class="{ 'active': currentSort == 'alphabet-desc' }">
                    {{ $t('LastLetter') }} → {{ $t('FirstLetter') }}
                  </button>
                </div>

                <div class="dropdown-group">
                  <h3>{{ $t('price') }}</h3>
                  <button @click="sortItems('price-asc')" :class="{ 'active': currentSort == 'price-asc' }">
                    {{ $t('priceAsc') }}
                  </button>
                  <button @click="sortItems('price-desc')" :class="{ 'active': currentSort == 'price-desc' }">
                    {{ $t('priceDesc') }}
                  </button>
                </div>
              </div>
            </transition>
          </div>

          <button class="filters-button" v-appear="{ delay: 700 }" @click="openFiltersModal">
            <img src="../assets/img/filter.png" :alt="$t('filterAlt')">
            {{ $t('filters') }} {{ activeFiltersCount > 0 ? `(${activeFiltersCount})` : '' }}
          </button>
        </div>
      </div>

      <div class="undefined-items-container" v-if="totalItems == 0 && !isCatalogLoading" v-appear="{ delay: 200 }">
        <p>{{ $t('undefinedItems') }}</p>
        <img src="../assets/img/fail.png" :alt="$t('failAlt')">
      </div>

      <div class="catalog-items" ref="catalogItemsRef">
        <div v-if="isCatalogLoading && items.length === 0" class="catalog-loader">
          <div class="spinner"></div>
        </div>
        <div v-for="(row, rowIndex) in chunkedVisibleItems" :key="`row-${rowIndex}`" class="items-row">
          <ItemCard v-for="(item, colIndex) in row" :key="item.uniqueId"
            v-memo="[item.uniqueId, item.availability, item.minPrice, item.images.length]"
            v-appear="{ delay: 100 + colIndex * 100 }" :id="item.id" :title="item.title" :images="item.images"
            :color="item.color" :sizes="item.sizes" :availability="item.availability" :uniqueId="item.uniqueId"
            :minPrice="item.minPrice" :tags="item.tags" :delay="600 + colIndex * 200" @openSizeModal="openSizesModal" />
        </div>
        <div v-if="hasMoreItems && !isLoadingMore" ref="loadTrigger" style="height: 1px; width: 100%;"></div>
        <div class="spinner" v-if="isLoadingMore"></div>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="isFiltersModalOpen" class="filters-modal" @click="closeFiltersModal">
        <div class="modal-content-filters" @click.stop>
          <div class="modal-header">
            <h2>{{ $t('filters').toUpperCase() }}</h2>
            <div class="modal-action-buttons">
              <button class="reset-filters" @click="resetFilters" :style="{
                opacity: activeFilters.length > 0 ? '1' : '.6',
                pointerEvents: activeFilters.length > 0 ? 'all' : 'none'
              }">
                <img src="../assets/img/broom.png" :alt="$t('resetFilters')">
                {{ $t('resetFilters') }}
              </button>
            </div>
            <button class="close-button" @click="closeFiltersModal">×</button>
          </div>
          <div class="filters-content" @click="clearHighlight">
            <div v-if="hasAppliedFiltersEver" class="checkbox-legend">
              <div class="legend-item">
                <span class="legend-color legend-server"></span>
                <span>{{ $t('legendServer') }}</span>
              </div>
              <div class="legend-item">
                <span class="legend-color legend-local"></span>
                <span>{{ $t('legendLocal') }}</span>
              </div>
              <div class="legend-item">
                <span class="legend-color legend-both"></span>
                <span>{{ $t('legendBoth') }}</span>
              </div>
            </div>
            <div class="filters-container">
              <div class="filters-heading-filters" style="grid-column: 1 / -1;">
                <div class="buttons-for-show-filters">
                  <button class="highlight-filters-btn" @click.stop="clearHighlight" @click="highlightServerFilters"
                    :disabled="isHighlighting || activeFiltersCount == 0">
                    <span>{{ $t('highlightServer') }}</span>
                  </button>
                  <button class="highlight-filters-btn" @click.stop="clearHighlight" @click="highlightLocalFilters"
                    :disabled="isHighlighting || activeFilters.length == 0">
                    <span>{{ $t('highlightActive') }}</span>
                  </button>
                </div>
                <div class="filter-group" data-filter-type="price">
                  <h3>
                    {{ $t('priceRange') }}
                    <span v-if="serverPriceRangeText" class="price-hint">
                      <span class="legend-color legend-server"></span> {{ serverPriceRangeText }}
                    </span>
                  </h3>
                  <div class="price-inputs">
                    <input type="text" :value="priceMinDisplay" @input="updatePriceMin" :placeholder="$t('minPrice')">
                    <span>-</span>
                    <input type="text" :value="priceMaxDisplay" @input="updatePriceMax" :placeholder="$t('maxPrice')">
                  </div>
                  <div class="price-slider">
                    <input type="range" :min="minAvailablePrice" :max="maxAvailablePrice" :value="filters.price.min"
                      @input="updateSliderMin" class="slider-min">
                    <input type="range" :min="minAvailablePrice" :max="maxAvailablePrice" :value="filters.price.max"
                      @input="updateSliderMax" class="slider-max">
                  </div>
                </div>
              </div>
              <div class="filter-group" data-filter-type="tags" v-if="availableFilters.tags && availableFilters.tags.length">
                <h3>{{ $t('tags') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="tag in availableFilters.tags" :key="tag.en" class="filter-option"
                    :class="getCheckboxClass('tags', tag)">
                    <input type="checkbox" :value="tag" v-model="filters.tags">
                    {{ $i18n.locale == 'en' ? tag.en : tag.ru }}
                  </label>
                </div>
              </div>

              <div class="filter-group" data-filter-type="brands" v-if="availableFilters.brands && availableFilters.brands.length">
                <h3>{{ $t('brand') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="brand in availableFilters.brands" :key="brand" class="filter-option"
                    :class="getCheckboxClass('brands', brand)">
                    <input type="checkbox" :value="brand" v-model="filters.brands">
                    {{ brand }}
                  </label>
                </div>
              </div>

              <div class="filter-group" data-filter-type="materials" v-if="availableFilters.materials && availableFilters.materials.length">
                <h3>{{ $t('material') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="material in availableFilters.materials" :key="material.en" class="filter-option"
                    :class="getCheckboxClass('materials', material)">
                    <input type="checkbox" :value="material" v-model="filters.materials">
                    {{ $i18n.locale == 'en' ? material.en : material.ru }}
                  </label>
                </div>
              </div>

              <div class="filter-group" data-filter-type="categories" v-if="availableFilters.categories && availableFilters.categories.length">
                <h3>{{ $t('category') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="category in availableFilters.categories" :key="category.en" class="filter-option"
                    :class="getCheckboxClass('categories', category)">
                    <input type="checkbox" :value="category" v-model="filters.categories">
                    {{ $i18n.locale == 'en' ? category.en : category.ru }}
                  </label>
                </div>
              </div>

              <div class="filter-group" data-filter-type="types" v-if="availableFilters.types && availableFilters.types.length">
                <h3>{{ $t('type') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="type in availableFilters.types" :key="type.en" class="filter-option"
                    :class="getCheckboxClass('types', type)">
                    <input type="checkbox" :value="type" v-model="filters.types">
                    {{ $i18n.locale == 'en' ? type.en : type.ru }}
                  </label>
                </div>
              </div>

              <div class="filter-group" data-filter-type="colors" v-if="availableFilters.colors && availableFilters.colors.length">
                <h3>{{ $t('color') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="color in availableFilters.colors" :key="color.en" class="filter-option"
                    :class="getCheckboxClass('colors', color)">
                    <input type="checkbox" :value="color" v-model="filters.colors">
                    {{ $i18n.locale == 'en' ? color.en : color.ru }}
                  </label>
                </div>
              </div>

              <div class="filter-group" data-filter-type="sizes" v-if="availableFilters.sizes && availableFilters.sizes.length">
                <h3>{{ $t('sizes') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="size in availableFilters.sizes" :key="size" class="filter-option"
                    :class="getCheckboxClass('sizes', size)">
                    <input type="checkbox" :value="size" v-model="filters.sizes">
                    {{ size }}
                  </label>
                </div>
              </div>

              <div class="filter-group" data-filter-type="availability" v-if="availableFilters.availability && availableFilters.availability.length">
                <h3>{{ $t('availability') }}</h3>
                <div class="filter-options" v-perfect-scrollbar>
                  <label v-for="status in availableFilters.availability" :key="status" class="filter-option"
                    :class="getCheckboxClass('availability', status)">
                    <input type="checkbox" :value="status" v-model="filters.availability">
                    {{ $t(status) }}
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
        <button class="apply-filters" @click="applyFilters" @click.stop="closeDontForget">
          <img src="../assets/img/apply.png" :alt="$t('applyFilters')">
          {{ $t('applyFilters') }}
        </button>
        <div class="dont-forget-apply-filters" @click.stop="closeDontForget">
          <p>{{ $t('dontForgetApplyFilters') }}</p>
          <img src="../assets/img/tap.png" :alt="$t('tap')">
        </div>
      </div>
    </Transition>
    <button class="active-filters-popover" @click="openFiltersModalAndHighlight"
      :style="{ transform: showFiltersPopover ? 'translateY(0%)' : 'translateY(-200px)', opacity: showFiltersPopover ? '1' : '0' }">
      <span v-appear.repeat="{ delay: 350 }">
        {{ $t('changeFilters') }}
        {{ activeFiltersCount > 0 ? `(${activeFiltersCount})` : '' }}
      </span>
    </button>
    <ItemSizesModal :uniqueId="currentModalUniqueId" :sizes="currentModalSizes" :isOpen="isSizesModalOpen"
      @close="closeSizesModal" @addToCart="handleAddToCart" />
    <button v-if="showScrollTopButton" class="scroll-top-button" @click="handleScrollTop">
      ↑
    </button>
  </main>
</template>

<script setup>
import { ref, onMounted, watch, computed, onBeforeUnmount, nextTick } from 'vue';
import { useIntersectionObserver, useScroll, useThrottleFn } from '@vueuse/core';
import ItemCard from '../components/ItemCard.vue';
import ItemSizesModal from '../components/ItemSizesModal.vue';
import { api } from '../api';

const emit = defineEmits(['page-loaded'])

let searchTimeout = null

const scrolledY = ref(0)
const SCROLL_THRESHOLD = 150

const isPageLoaded = ref(false)
const items = ref([])
const totalItems = ref(0)
const currentPage = ref(1)
const isLoadingMore = ref(false)

const catalogItemsRef = ref(null)
const showScrollTopButton = ref(false)
const SCROLL_TOP_THRESHOLD = 500

const loadTrigger = ref(null)

const currentSort = ref('')
const searchQuery = ref('')
const filters = ref({
  price: { min: null, max: null },
  brands: [],
  materials: [],
  categories: [],
  types: [],
  colors: [],
  tags: [],
  sizes: [],
  availability: []
})

const isSearchExpanded = ref(false)
const searchInput = ref(null)
const inputBox = ref(null)
const isMainSortDropdownOpen = ref(false)
const isFiltersModalOpen = ref(false)
const availableFilters = ref({})

const serverAppliedFilters = ref([])

const priceMinDisplay = ref('')
const priceMaxDisplay = ref('')

const isKnowForApply = ref(false);

const hasAppliedFiltersEver = ref(false)
const STORAGE_KEY_APPLIED_EVER = 'catalog_has_applied_filters'
const STORAGE_KEY = 'catalog_filters_applied_knowledge'

const isSizesModalOpen = ref(false);
const currentModalUniqueId = ref('');
const currentModalSizes = ref([]);

const isCatalogLoading = ref(false);

const windowWidth = ref(window.innerWidth);

const { y } = useScroll(window);

const loadKnowForApply = () => {
  const saved = localStorage.getItem(STORAGE_KEY)
  isKnowForApply.value = saved === 'true'
}

const saveKnowForApply = () => {
  localStorage.setItem(STORAGE_KEY, isKnowForApply.value ? 'true' : 'false')
}

const loadHasAppliedEver = () => {
  const saved = localStorage.getItem(STORAGE_KEY_APPLIED_EVER)
  hasAppliedFiltersEver.value = saved === 'true'
}

const saveHasAppliedEver = () => {
  localStorage.setItem(STORAGE_KEY_APPLIED_EVER, hasAppliedFiltersEver.value ? 'true' : 'false')
}

const itemsPerRow = computed(() => {
  if (windowWidth.value < 1440) return 2;
  return 4;
});

const itemsPerLoad = computed(() => {
  const perRow = itemsPerRow.value;
  if (perRow === 4) return 20;
  if (perRow === 2) return 12;
  return 8;
});

const minAvailablePrice = computed(() => {
  return availableFilters.value.minPrice ?? 0
})
const maxAvailablePrice = computed(() => {
  return availableFilters.value.maxPrice ?? 100000
})

const activeFilters = computed(() => {
  const active = []
  if (filters.value.price.min !== null || filters.value.price.max !== null) {
    active.push({
      key: 'price',
      value: 'price',
      displayName: `${filters.value.price.min || '0'} - ${filters.value.price.max || '∞'} ₽`
    })
  }
  Object.keys(filters.value).forEach(key => {
    if (key !== 'price' && filters.value[key].length > 0) {
      filters.value[key].forEach(value => {
        active.push({
          key: key,
          value: value,
          displayName: `${getFilterLabel(key)}: ${value}`
        })
      })
    }
  })
  return active
})

const activeFiltersCount = computed(() => {
  if (!serverAppliedFilters.value.length) return 0
  let total = 0
  for (const filter of serverAppliedFilters.value) {
    if (Array.isArray(filter.value)) {
      total += filter.value.length
    } else {
      total += 1
    }
  }
  return total
})

const hasMoreItems = computed(() => {
  return items.value.length < totalItems.value
})

const chunkedVisibleItems = computed(() => {
  const result = [];
  for (let i = 0; i < items.value.length; i += itemsPerRow.value) {
    result.push(items.value.slice(i, i + itemsPerRow.value));
  }
  return result;
});

const secondRowElement = computed(() => {
  if (!catalogItemsRef.value) return null;
  const rows = catalogItemsRef.value.querySelectorAll('.items-row');
  return rows[0] || null;
})

const getFilterLabel = (type) => {
  const labels = {
    brands: 'Бренд',
    materials: 'Материал',
    categories: 'Категория',
    types: 'Тип',
    colors: 'Цвет',
    tags: 'Тег',
    sizes: 'Размер',
    availability: 'Наличие'
  }
  return labels[type] || type
}

const updatePriceMin = (event) => {
  let raw = event.target.value.replace(/\s/g, '')
  if (raw === '') {
    filters.value.price.min = null
    priceMinDisplay.value = ''
  } else {
    let num = Number(raw)
    if (!isNaN(num)) {
      filters.value.price.min = num
      priceMinDisplay.value = num.toLocaleString('ru-RU')
    } else {
      priceMinDisplay.value = filters.value.price.min !== null
        ? filters.value.price.min.toLocaleString('ru-RU')
        : ''
    }
  }
}

const updateSliderMin = (event) => {
  let val = parseFloat(event.target.value)
  if (isNaN(val)) val = null
  filters.value.price.min = val
}

const updatePriceMax = (event) => {
  let raw = event.target.value.replace(/\s/g, '')
  if (raw === '') {
    filters.value.price.max = null
    priceMaxDisplay.value = ''
  } else {
    let num = Number(raw)
    if (!isNaN(num)) {
      filters.value.price.max = num
      priceMaxDisplay.value = num.toLocaleString('ru-RU')
    } else {
      priceMaxDisplay.value = filters.value.price.max !== null
        ? filters.value.price.max.toLocaleString('ru-RU')
        : ''
    }
  }
}

const updateSliderMax = (event) => {
  let val = parseFloat(event.target.value)
  if (isNaN(val)) val = null
  filters.value.price.max = val
}

const serverPriceRangeText = computed(() => {
  if (!serverAppliedFilters.value.length) return ''

  let min = null
  let max = null

  for (const filter of serverAppliedFilters.value) {
    let key = filter.key
    if (key === 'priceRange') key = 'price'

    if (key === 'price') {
      const val = filter.value
      const parts = val.split(/[–\-]/).map(p => parseInt(p.trim(), 10))
      min = parts[0] || null
      max = parts[1] || null
    }
    else if (key === 'priceMin') {
      const num = Number(filter.value)
      if (!isNaN(num)) min = num
    }
    else if (key === 'priceMax') {
      const num = Number(filter.value)
      if (!isNaN(num)) max = num
    }
  }

  if (min !== null && max !== null) {
    return `${min.toLocaleString('ru-RU')} – ${max.toLocaleString('ru-RU')} ₽`
  } else if (min !== null) {
    return `от ${min.toLocaleString('ru-RU')} ₽`
  } else if (max !== null) {
    return `до ${max.toLocaleString('ru-RU')} ₽`
  }
  return ''
})

const isHighlighting = ref(false)
let highlightTimeouts = []
let highlightGroups = []

const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms))

const isServerFilterGroupActive = (type) => {
  if (!serverAppliedFilters.value.length) return false

  for (const filter of serverAppliedFilters.value) {
    let key = filter.key
    if (key === 'priceRange') key = 'price'
    if (key === 'material') key = 'materials'
    if (key === 'category') key = 'categories'
    if (key === 'type') key = 'types'
    if (key === 'color') key = 'colors'

    if (key === type) {
      if (Array.isArray(filter.value)) {
        return filter.value.length > 0
      } else {
        return filter.value !== null && filter.value !== undefined && filter.value !== ''
      }
    }
  }
  return false
}

const serverValuesMap = computed(() => {
  const map = {}
  if (!serverAppliedFilters.value.length) return map

  for (const filter of serverAppliedFilters.value) {
    let key = filter.key
    if (key === 'priceRange') key = 'price'
    if (key === 'material') key = 'materials'
    if (key === 'category') key = 'categories'
    if (key === 'type') key = 'types'
    if (key === 'color') key = 'colors'

    let values = []
    if (Array.isArray(filter.value)) values = filter.value
    else if (filter.value !== null && filter.value !== '') values = [filter.value]

    map[key] = new Set(values.map(v => typeof v === 'object' ? (v.ru || v.en) : String(v)))
  }
  return map
})

const getCheckboxClass = (filterType, value) => {
  const valueStr = typeof value === 'object' ? (value.ru || value.en) : String(value)
  const isServer = serverValuesMap.value[filterType]?.has(valueStr) || false
  let isLocal = false
  const localArray = filters.value[filterType]

  if (localArray && Array.isArray(localArray)) {
    isLocal = localArray.some(item => {
      const itemStr = typeof item === 'object' ? (item.ru || item.en) : String(item)
      return itemStr === valueStr
    })
  }

  if (isServer && isLocal) return 'checkbox-both'
  if (isServer) return 'checkbox-server'
  if (isLocal) return 'checkbox-local'
  return ''
}

const isLocalFilterGroupActive = (type) => {
  switch (type) {
    case 'price':
      return filters.value.price.min !== null || filters.value.price.max !== null
    case 'tags':
      return filters.value.tags.length > 0
    case 'brands':
      return filters.value.brands.length > 0
    case 'materials':
      return filters.value.materials.length > 0
    case 'categories':
      return filters.value.categories.length > 0
    case 'types':
      return filters.value.types.length > 0
    case 'colors':
      return filters.value.colors.length > 0
    case 'sizes':
      return filters.value.sizes.length > 0
    case 'availability':
      return filters.value.availability.length > 0
    default:
      return false
  }
}

const highlightServerFilters = async () => {
  await highlightFiltersByType('server')
}

const highlightLocalFilters = async () => {
  await highlightFiltersByType('local')
}

const highlightFiltersByType = async (type) => {
  if (isHighlighting.value) {
    clearHighlight()
    await delay(50)
  }

  isHighlighting.value = true

  try {
    const groups = document.querySelectorAll('.filters-container .filter-group')
    const activeGroups = []

    for (const group of groups) {
      const filterType = group.dataset.filterType
      const isActive = type === 'server'
        ? isServerFilterGroupActive(filterType)
        : isLocalFilterGroupActive(filterType)
      if (isActive) activeGroups.push(group)
    }

    if (activeGroups.length === 0) {
      isHighlighting.value = false
      return
    }

    const highlightClass = type === 'server'
      ? 'highlight-filter-group-server'
      : 'highlight-filter-group-local'

    const highlightOne = (group) => {
      const isInViewport = (el) => {
        const rect = el.getBoundingClientRect()
        const container = document.querySelector('.filters-content')
        if (!container) return true
        const containerRect = container.getBoundingClientRect()
        return rect.top >= containerRect.top && rect.bottom <= containerRect.bottom
      }
      if (!isInViewport(group)) {
        group.scrollIntoView({ behavior: 'smooth', block: 'center' })
      }

      const timeoutId = setTimeout(() => {
        group.classList.add(highlightClass)
        highlightGroups.push(group)
      }, 50)
      highlightTimeouts.push(timeoutId)
    }

    for (let i = 0; i < activeGroups.length; i++) {
      const timeoutId = setTimeout(highlightOne, i * 150, activeGroups[i])
      highlightTimeouts.push(timeoutId)
    }

    const totalDuration = (activeGroups.length - 1) * 150 + 150
    await delay(totalDuration)
  } finally {
    isHighlighting.value = false
  }
}

const clearHighlight = () => {
  highlightTimeouts.forEach(timeout => clearTimeout(timeout))
  highlightTimeouts = []
  highlightGroups.forEach(group => {
    group.classList.remove('highlight-filter-group-server', 'highlight-filter-group-local')
  })
  highlightGroups = []
  isHighlighting.value = false
}

async function fetchItems() {
  const isFirstPage = currentPage.value === 1;

  if (isFirstPage) {
    isCatalogLoading.value = true;
  } else {
    isLoadingMore.value = true;
  }

  try {
    const params = {
      page: currentPage.value,
      limit: itemsPerLoad.value,
      sort: currentSort.value || undefined,
      search: searchQuery.value || undefined,
      lang: 'ru',
      priceMin: filters.value.price.min,
      priceMax: filters.value.price.max,
      'brands[]': filters.value.brands,
      'materials[]': filters.value.materials.map(m => m.ru || m),
      'categories[]': filters.value.categories.map(c => c.ru || c),
      'types[]': filters.value.types.map(t => t.ru || t),
      'colors[]': filters.value.colors.map(c => c.ru || c),
      'tags[]': filters.value.tags.map(t => t.ru || t),
      'sizes[]': filters.value.sizes,
      'availability[]': filters.value.availability
    };

    const response = await api.getItems(params);
    const data = response.data;

    if (currentPage.value === 1) {
      items.value = data.items;
      serverAppliedFilters.value = data.appliedFilters || [];
    } else {
      items.value = [...items.value, ...data.items];
    }
    totalItems.value = data.total;
  } catch (error) {
    console.error('Ошибка загрузки каталога:', error);
  } finally {
    if (isFirstPage) {
      isCatalogLoading.value = false;
    } else {
      isLoadingMore.value = false;
    }
    isPageLoaded.value = true;
    emit('page-loaded', true);
  }
}

function sortItems(sortType) {
  if (currentSort.value === sortType) {
    currentSort.value = ''
  } else {
    currentSort.value = sortType
  }
  isMainSortDropdownOpen.value = false
  resetCatalogAndReload()
}

function expandSearch() {
  isSearchExpanded.value = true
  inputBox.value.style.transform = 'scale(1.02)'
  setTimeout(() => {
    searchInput.value?.focus()
  }, 500)
}

function handleSearchBlur() {
  if (!searchQuery.value) {
    isSearchExpanded.value = false
  }
  inputBox.value.style.transform = 'scale(1)'
}

function applyFilters() {
  if (!hasAppliedFiltersEver.value) {
    hasAppliedFiltersEver.value = true
    saveHasAppliedEver()
  }
  isKnowForApply.value = true
  saveKnowForApply()
  resetCatalogAndReload()
  closeFiltersModal()
  window.scrollTo(0, 0)
}

function resetFilters() {
  clearHighlight()
  filters.value = {
    price: { min: null, max: null },
    brands: [],
    materials: [],
    categories: [],
    types: [],
    colors: [],
    tags: [],
    sizes: [],
    availability: []
  }
}

function openFiltersModal() {
  isFiltersModalOpen.value = true
  document.body.style.overflow = 'hidden'
}

const openFiltersModalAndHighlight = () => {
  openFiltersModal()
  nextTick(() => {
    setTimeout(() => {
      highlightServerFilters()
    }, 200)
  })
}

function closeFiltersModal() {
  clearHighlight()
  if (!isKnowForApply.value && activeFilters.value.length > 0) {
    const applyModal = document.querySelector('.dont-forget-apply-filters')
    applyModal.style.display = 'flex'
    setTimeout(() => {
      applyModal.style.opacity = '1'
    }, 50)
  } else {
    isFiltersModalOpen.value = false
    document.body.style.overflow = 'auto'
  }
}

function closeDropdownOnClickOutside(loc) {
  if (loc === 'main') isMainSortDropdownOpen.value = false
}

function toggleDropdown(type, loc) {
  if (type === 'sort' && loc === 'main') isMainSortDropdownOpen.value = !isMainSortDropdownOpen.value
}

function loadMoreItems() {
  if (!hasMoreItems.value || isLoadingMore.value) return
  currentPage.value++
  fetchItems()
}

function closeDontForget() {
  isKnowForApply.value = true
  saveKnowForApply()
  const applyModal = document.querySelector('.dont-forget-apply-filters')
  applyModal.style.opacity = '0'
  setTimeout(() => {
    applyModal.style.display = 'none'
  }, 300)
}

const updateScroll = () => {
  scrolledY.value = window.scrollY
}

const showFiltersPopover = computed(() => {
  return scrolledY.value > SCROLL_THRESHOLD && activeFiltersCount.value > 0
})

function openSizesModal({ uniqueId, sizes }) {
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

let resizeTimer = null;
function handleResize() {
  if (resizeTimer) clearTimeout(resizeTimer);
  resizeTimer = setTimeout(() => {
    const newWidth = window.innerWidth;
    if ((newWidth < 1440 && windowWidth.value >= 1440) || (newWidth > 1440 && windowWidth.value <= 1440)) {
      windowWidth.value = newWidth;
      resetCatalogAndReload();
    }
  }, 200);
}

function resetCatalogAndReload() {
  currentPage.value = 1;
  items.value = [];
  totalItems.value = 0;
  isLoadingMore.value = false;
  fetchItems();
}

function handleAddToCart({ uniqueId, size }) { }

async function loadFilters() {
  try {
    const response = await api.getFilters()
    availableFilters.value = response.data
  } catch (error) {
    console.error('Ошибка загрузки фильтров:', error)
  }
}

async function handleScrollTop() {
  const secondRow = secondRowElement.value;
  if (secondRow) {
    secondRow.scrollIntoView({ behavior: 'auto', block: 'end' });
    await nextTick();
  }
  window.scrollTo({ top: 0, behavior: 'auto' });
  currentPage.value = 1;
  await fetchItems();
  window.scrollTo({ top: 0, behavior: 'auto' });
}

useIntersectionObserver(
  loadTrigger,
  ([{ isIntersecting }]) => {
    if (isIntersecting && hasMoreItems.value && !isLoadingMore.value) {
      loadMoreItems();
    }
  },
  { threshold: 0.1, rootMargin: '0px 0px 0px 0px' }
);

watch(y, (value) => {
  showScrollTopButton.value = value > SCROLL_TOP_THRESHOLD;
});

watch(searchQuery, (newVal) => {
  if (newVal !== '') isSearchExpanded.value = true
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    resetCatalogAndReload()
  }, 300)
})

watch(() => filters.value.price.min, (newVal) => {
  if (newVal !== null && !isNaN(newVal)) {
    priceMinDisplay.value = newVal.toLocaleString('ru-RU')
  } else {
    priceMinDisplay.value = ''
  }
}, { immediate: true })

watch(() => filters.value.price.max, (newVal) => {
  if (newVal !== null && !isNaN(newVal)) {
    priceMaxDisplay.value = newVal.toLocaleString('ru-RU')
  } else {
    priceMaxDisplay.value = ''
  }
}, { immediate: true })

onMounted(() => {
  loadKnowForApply()
  loadHasAppliedEver()
  fetchItems()
  loadFilters()
  document.addEventListener('click', (event) => {
    if (!event.target.closest('.dropdown-wrapper')) {
      isMainSortDropdownOpen.value = false
    }
  })
  window.addEventListener('scroll', updateScroll)
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateScroll)
  window.removeEventListener('resize', handleResize)
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/catalog';
</style>