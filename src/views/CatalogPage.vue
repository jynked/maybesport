<template>
  <main v-if="isPageLoaded">
    <div class="catalog">
      <div class="catalog-block">
        <h1 v-appear="{ delay: 400 }">{{ $t('catalog') }}</h1>
        <div class="filters-main-buttons-block">
          <div class="search-container" v-appear="{ delay: 500 }" :class="{ 'expanded': isSearchExpanded }" ref="inputBox">
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
            {{ $t('filters') }}
            <span v-if="activeFiltersCount > 0" class="filters-counter">{{ activeFiltersCount }}</span>
          </button>
        </div>
      </div>

      <div v-if="activeFilters.length > 0" class="active-filters">
        <button class="clear-all-filters" @click="clearAllFilters">
          {{ $t('clearAllFilters') }}
        </button>
        <div class="active-filters-list">
          <div v-for="filter in activeFilters" :key="filter.key + filter.value" class="active-filter">
            <span>{{ getFilterDisplayName(filter) }}</span>
            <button @click="removeFilter(filter.key, filter.value)" class="remove-filter">×</button>
          </div>
        </div>
      </div>

      <div class="catalog-items">
        <div v-for="(row, rowIndex) in chunkedVisibleItems" :key="`row-${rowIndex}`" class="items-row">
          <ItemCard v-for="(item, colIndex) in row" :key="item.uniqueId"
            v-appear="{ delay: 100 + colIndex * 100 }" :id="item.id" :title="item.title"
            :images="item.images" :color="item.color" :sizes="item.sizes" :availability="item.availability"
            :uniqueId="item.uniqueId" :minPrice="item.minPrice" :tags="item.tags"
            :delay="600 + colIndex * 200" />
        </div>
        <button v-if="hasMoreItems" class="show-more" @click="loadMoreItems" :disabled="isLoadingMore">
          <span v-if="!isLoadingMore">{{ $t('showMore') }}</span>
          <span v-else>{{ $t('Loading') }}...</span>
        </button>
      </div>
    </div>


        <Transition name="modal">
            <div v-if="isFiltersModalOpen" class="filters-modal" @click="closeFiltersModal">
                <div class="modal-content-filters" @click.stop>
                    <div class="modal-header">
                        <h2>{{ $t('filters').toUpperCase() }}</h2>
                        <button class="close-button" @click="closeFiltersModal">×</button>
                    </div>
                    <div class="filters-content">
                        <div class="filters-container">
                            <div v-if="activeFilters.length > 0" class="modal-active-filters">
                                <button class="clear-all-filters" @click="clearAllFilters">
                                    {{ $t('clearAllFilters') }}
                                </button>
                                <div class="active-filters-list">
                                    <div v-for="filter in activeFilters" :key="filter.key + filter.value"
                                        class="active-filter">
                                        <span>{{ getFilterDisplayName(filter) }}</span>
                                        <button @click="removeFilter(filter.key, filter.value)" class="remove-filter">
                                            ×
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('priceRange') }}</h3>
                                <div class="price-inputs">
                                    <input type="number" v-model="filters.price.min" :placeholder="$t('minPrice')"
                                        @input="applyPriceFilter">
                                    <span>-</span>
                                    <input type="number" v-model="filters.price.max" :placeholder="$t('maxPrice')"
                                        @input="applyPriceFilter">
                                </div>
                                <div class="price-slider">
                                    <input type="range" :min="minAvailablePrice" :max="maxAvailablePrice"
                                        v-model="filters.price.min" @input="applyPriceFilter" class="slider-min">
                                    <input type="range" :min="minAvailablePrice" :max="maxAvailablePrice"
                                        v-model="filters.price.max" @input="applyPriceFilter" class="slider-max">
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('brand') }}</h3>
                                <div class="filter-options">
                                    <label v-for="brand in availableFilters.brands" :key="brand" class="filter-option">
                                        <input type="checkbox" :value="brand" v-model="filters.brands">
                                        <span class="checkmark"></span>
                                        {{ brand }}
                                        <span class="option-count">({{ getBrandCount(brand) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('country') }}</h3>
                                <div class="filter-options">
                                    <label v-for="country in availableFilters.countries" :key="country"
                                        class="filter-option">
                                        <input type="checkbox" :value="country" v-model="filters.countries">
                                        <span class="checkmark"></span>
                                        {{ country }}
                                        <span class="option-count">({{ getCountryCount(country) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('material') }}</h3>
                                <div class="filter-options">
                                    <label v-for="material in availableFilters.materials" :key="material.en"
                                        class="filter-option">
                                        <input type="checkbox" :value="material" v-model="filters.materials">
                                        <span class="checkmark"></span>
                                        {{ $i18n.locale == 'en' ? material.en : material.ru }}
                                        <span class="option-count">({{ getMaterialCount(material) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('category') }}</h3>
                                <div class="filter-options">
                                    <label v-for="category in availableFilters.categories" :key="category.en"
                                        class="filter-option">
                                        <input type="checkbox" :value="category" v-model="filters.categories">
                                        <span class="checkmark"></span>
                                        {{ $i18n.locale == 'en' ? category.en : category.ru }}
                                        <span class="option-count">({{ getCategoryCount(category) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('type') }}</h3>
                                <div class="filter-options">
                                    <label v-for="type in availableFilters.types" :key="type.en" class="filter-option">
                                        <input type="checkbox" :value="type" v-model="filters.types">
                                        <span class="checkmark"></span>
                                        {{ $i18n.locale == 'en' ? type.en : type.ru }}
                                        <span class="option-count">({{ getTypeCount(type) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('color') }}</h3>
                                <div class="filter-options">
                                    <label v-for="color in availableFilters.colors" :key="color.en"
                                        class="filter-option">
                                        <input type="checkbox" :value="color" v-model="filters.colors">
                                        <span class="checkmark"></span>
                                        {{ $i18n.locale == 'en' ? color.en : color.ru }}
                                        <span class="option-count">({{ getColorCount(color) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('tags') }}</h3>
                                <div class="filter-options">
                                    <label v-for="tag in availableFilters.tags" :key="tag.en" class="filter-option">
                                        <input type="checkbox" :value="tag" v-model="filters.tags">
                                        <span class="checkmark"></span>
                                        {{ $i18n.locale == 'en' ? tag.en : tag.ru }}
                                        <span class="option-count">({{ getTagCount(tag) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('sizes') }}</h3>
                                <div class="filter-options">
                                    <label v-for="size in availableFilters.sizes" :key="size" class="filter-option">
                                        <input type="checkbox" :value="size" v-model="filters.sizes">
                                        <span class="checkmark"></span>
                                        {{ size }}
                                        <span class="option-count">({{ getSizeCount(size) }})</span>
                                    </label>
                                </div>
                            </div>

                            <div class="filter-group">
                                <h3>{{ $t('availability') }}</h3>
                                <div class="filter-options">
                                    <label v-for="status in availableFilters.availability" :key="status"
                                        class="filter-option">
                                        <input type="checkbox" :value="status" v-model="filters.availability">
                                        <span class="checkmark"></span>
                                        {{ $t(status) }}
                                        <span class="option-count">({{ getAvailabilityCount(status) }})</span>
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button class="apply-filters" @click="applyFilters">{{ $t('applyFilters') }}</button>
                        <button class="reset-filters" @click="resetFilters">{{ $t('resetFilters') }}</button>
                    </div>
                </div>
            </div>
        </Transition>
    </main>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import ItemCard from '../components/ItemCard.vue'
import { api } from '../api'

const emit = defineEmits(['page-loaded'])

const isPageLoaded = ref(false)
const items = ref([])
const totalItems = ref(0)
const currentPage = ref(1)
const itemsPerLoad = 20
const isLoadingMore = ref(false)

const currentSort = ref('')
const searchQuery = ref('')
const filters = ref({
  price: { min: null, max: null },
  brands: [],
  countries: [],
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

const activeFiltersCount = computed(() => activeFilters.value.length)

const hasMoreItems = computed(() => {
  return items.value.length < totalItems.value
})

const chunkedVisibleItems = computed(() => {
  const itemsPerRow = 4
  const result = []
  for (let i = 0; i < items.value.length; i += itemsPerRow) {
    result.push(items.value.slice(i, i + itemsPerRow))
  }
  return result
})

const getFilterLabel = (type) => {
    const labels = {
        brands: 'Бренд',
        countries: 'Страна',
        materials: 'Материал',
        categories: 'Категория',
        types: 'Тип',
        colors: 'Цвет',
        tags: 'Тег',
        sizes: 'Размер',
        availability: 'Наличие'
    };
    return labels[type] || type;
};

function getFilterDisplayName(filter) { return filter.displayName }

async function fetchItems() {
  isLoadingMore.value = currentPage.value > 1
  try {
    const params = {
      page: currentPage.value,
      limit: itemsPerLoad,
      sort: currentSort.value || undefined,
      search: searchQuery.value || undefined,
      lang: 'ru',
      priceMin: filters.value.price.min,
      priceMax: filters.value.price.max,
      'brands[]': filters.value.brands,
      'countries[]': filters.value.countries,
      'materials[]': filters.value.materials.map(m => m.ru || m),
      'categories[]': filters.value.categories.map(c => c.ru || c),
      'types[]': filters.value.types.map(t => t.ru || t),
      'colors[]': filters.value.colors.map(c => c.ru || c),
      'tags[]': filters.value.tags.map(t => t.ru || t),
      'sizes[]': filters.value.sizes,
      'availability[]': filters.value.availability
    }

    const response = await api.getItems(params)
    const data = response.data

    if (currentPage.value === 1) {
      items.value = data.items
    } else {
      items.value = [...items.value, ...data.items]
    }
    totalItems.value = data.total

    // Также можно получить доступные фильтры из ответа, если сервер их возвращает (опционально)
    // if (data.availableFilters) availableFilters.value = data.availableFilters
  } catch (error) {
    console.error('Ошибка загрузки каталога:', error)
  } finally {
    isLoadingMore.value = false
    isPageLoaded.value = true
    emit('page-loaded', true)
  }
}

function resetAndFetch() {
  currentPage.value = 1
  fetchItems()
}

function sortItems(sortType) {
  currentSort.value = sortType
  isMainSortDropdownOpen.value = false
  resetAndFetch()
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

let searchTimeout = null
watch(searchQuery, (newVal) => {
  if (newVal !== '') isSearchExpanded.value = true
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    resetAndFetch()
  }, 300)
})

function applyFilters() {
  resetAndFetch()
  closeFiltersModal()
}

function resetFilters() {
  filters.value = {
    price: { min: null, max: null },
    brands: [],
    countries: [],
    materials: [],
    categories: [],
    types: [],
    colors: [],
    tags: [],
    sizes: [],
    availability: []
  }
  resetAndFetch()
  closeFiltersModal()
}

function clearAllFilters() {
  resetFilters()
}

function removeFilter(key, value) {
  if (key === 'price') {
    filters.value.price = { min: null, max: null }
  } else {
    filters.value[key] = filters.value[key].filter(v => v !== value)
  }
  resetAndFetch()
}

function openFiltersModal() {
  isFiltersModalOpen.value = true
  document.body.style.overflow = 'hidden'
}
function closeFiltersModal() {
  isFiltersModalOpen.value = false
  document.body.style.overflow = 'auto'
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

onMounted(() => {
  fetchItems()
  document.addEventListener('click', (event) => {
    if (!event.target.closest('.dropdown-wrapper')) {
      isMainSortDropdownOpen.value = false
    }
  })
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/catalog';
</style>
