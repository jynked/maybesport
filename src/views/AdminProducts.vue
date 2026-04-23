<template>
  <main>
    <div class="container-fluid py-4">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление товарами</h1>
        <button class="btn btn-primary" @click="openCreateModal" :disabled="isLoading || isSaving">
          <i class="bi bi-plus-lg"></i> Создать товар
        </button>
      </div>

      <div v-if="isLoading" class="text-center py-5">
        <div class="spinner-border text-primary" role="status">
          <span class="visually-hidden">Загрузка...</span>
        </div>
      </div>
      <div v-else class="table-responsive">
        <table class="table table-striped table-hover">
          <thead>
            <tr>
              <th>ID</th>
              <th>Изображение</th>
              <th>Название (RU)</th>
              <th>Название (EN)</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in items" :key="item.id" v-memo="[item.id, item.title]">
              <td>{{ item.id }}</td>
              <td>
                <img :src="item.items?.[0]?.images?.[0] || '/placeholder.jpg'"
                  style="width: 50px; height: 50px; object-fit: cover;" />
              </td>
              <td>{{ item.title.ru }}</td>
              <td>{{ item.title.en }}</td>
              <td>
                <button class="btn btn-sm btn-outline-info" @click="openEditModal(item)"
                  :disabled="isLoading || isSaving">
                  Характеристики
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="modal fade" id="itemModal" tabindex="-1" aria-labelledby="itemModalLabel" aria-hidden="true"
        ref="modal">
        <div class="modal-dialog modal-xl modal-dialog-scrollable">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="itemModalLabel">
                {{ isEditing ? 'Редактирование товара' : 'Создание товара' }}
              </h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
              <div v-if="isLoadingItem" class="text-center py-5">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">Загрузка данных...</span>
                </div>
              </div>
              <fieldset v-else :disabled="isSaving">
                <form @submit.prevent="saveItem">
                  <div class="row mb-3">
                    <div class="col-md-6">
                      <label class="form-label">Название (RU)</label>
                      <input type="text" class="form-control" v-model="form.title.ru" required />
                    </div>
                    <div class="col-md-6">
                      <label class="form-label">Название (EN)</label>
                      <input type="text" class="form-control" v-model="form.title.en" required />
                    </div>
                  </div>

                  <div class="row mb-3">
                    <div class="col-md-6">
                      <label class="form-label">Описание (RU)</label>
                      <textarea class="form-control" v-model="form.description.ru" rows="2" required></textarea>
                    </div>
                    <div class="col-md-6">
                      <label class="form-label">Описание (EN)</label>
                      <textarea class="form-control" v-model="form.description.en" rows="2" required></textarea>
                    </div>
                  </div>

                  <div class="row mb-3">
                    <div class="col-md-3">
                      <label class="form-label">Категория</label>
                      <select class="form-select" v-model="selectedCategory">
                        <option value="" disabled>Не выбрано</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.ru }} / {{ cat.en }}
                        </option>
                      </select>
                    </div>
                    <div class="col-md-3">
                      <label class="form-label">Страна</label>
                      <select class="form-select" v-model="selectedCountry">
                        <option value="" disabled>Не выбрано</option>
                        <option v-for="c in countries" :key="c.id" :value="c.id">{{ c.ru }} / {{ c.en }}</option>
                      </select>
                    </div>
                    <div class="col-md-3">
                      <label class="form-label">Тип</label>
                      <select class="form-select" v-model="selectedType" :disabled="!selectedCategory">
                        <option value="" disabled>Не выбрано</option>
                        <option v-for="t in filteredTypes" :key="t.id" :value="t.id">{{ t.ru }} / {{ t.en }}</option>
                      </select>
                    </div>
                    <div class="col-md-3">
                      <label class="form-label">Бренд</label>
                      <select class="form-select" v-model="form.brand">
                        <option value="" disabled>Не выбран</option>
                        <option v-for="b in brands" :key="b" :value="b">{{ b }}</option>
                      </select>
                    </div>
                  </div>

                  <div class="mb-3 form-row-custom">
                    <label class="form-label">Состав</label>
                    <div v-for="(material, index) in form.structure" :key="index" class="row mb-2">
                      <div class="col-md-4">
                        <select class="form-select" :value="getMaterialId(material)"
                          @change="updateMaterial(index, $event.target.value)">
                          <option value="" disabled>Выберите материал</option>
                          <option v-for="m in materials" :key="m.id" :value="m.id">{{ m.ru }} / {{ m.en }}</option>
                        </select>
                      </div>
                      <div class="col-md-2">
                        <input type="number" class="form-control" placeholder="%" v-model.number="material.percent" />
                      </div>
                      <div class="col-md-2">
                        <button type="button" class="btn btn-outline-danger" @click="removeMaterial(index)">
                          Удалить
                        </button>
                      </div>
                    </div>
                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="addMaterial">
                      Добавить материал
                    </button>
                  </div>

                  <div class="mb-4">
                    <h6>Варианты товара (подтовары)</h6>
                    <div v-for="(subItem, idx) in form.items" :key="idx" class="card mb-3">
                      <div class="card-body">
                        <div class="d-flex justify-content-between align-items-center mb-2">
                          <h6>Подтовар #{{ idx + 1 }}</h6>
                          <button type="button" class="btn btn-sm btn-danger" @click="removeSubItem(idx)">
                            Удалить подтовар
                          </button>
                        </div>

                        <div class="mb-3 form-row-custom">
                          <label>Цвета</label>
                          <div v-for="(color, colorIdx) in subItem.color" :key="colorIdx" class="row mb-2">
                            <div class="col-md-5">
                              <select class="form-select" :value="getColorCode(color)"
                                @change="updateColor(idx, colorIdx, $event.target.value)">
                                <option value="" disabled>Выберите цвет</option>
                                <option v-for="c in colors" :key="c.code" :value="c.code">
                                  {{ c.ru }} / {{ c.en }}
                                </option>
                              </select>
                            </div>
                            <div class="col-md-5 d-flex align-items-center" v-if="color.ru.length > 0">
                              <span
                                :style="{ backgroundColor: getColorHex(color), width: '30px', height: '30px', display: 'inline-block', borderRadius: '4px', marginRight: '10px' }"></span>
                              <span>{{ color.ru }} / {{ color.en }}</span>
                            </div>
                            <div class="col-md-2" style="margin-left: auto;">
                              <button type="button" class="btn btn-outline-danger" @click="removeColor(idx, colorIdx)">
                                Удалить
                              </button>
                            </div>
                          </div>
                          <button type="button" class="btn btn-sm btn-outline-secondary" @click="addColor(idx)">
                            Добавить цвет
                          </button>
                        </div>

                        <div class="mb-3">
                          <label>Теги</label>
                          <div class="d-flex flex-wrap gap-2">
                            <button v-for="tag in tags" :key="tag.id" type="button" class="btn"
                              :class="isTagSelected(subItem.tags, tag) ? 'btn-danger' : 'btn-success'" style="font-size: .75rem;"
                              @click="toggleTag(subItem, tag)">
                              {{ tag.ru }} / {{ tag.en }}
                            </button>
                          </div>
                        </div>

                        <div class="mb-3 form-row-custom">
                          <div v-for="(size, sizeIdx) in subItem.sizes" :key="sizeIdx"
                            class="row mb-2 align-items-center">
                            <div class="col-md-2" style="margin-top: auto;">
                              <label>Размер</label>
                              <input type="text" class="form-control" placeholder="0" v-model="size.size" />
                            </div>
                            <div class="col-md-2">
                              <label>Цена (CNY)</label>
                              <input type="number" class="form-control" placeholder="Юани"
                                v-model.number="size.priceCny" />
                            </div>
                            <div class="col-md-2">
                              <label>Цена (RUB)</label>
                              <input type="number" class="form-control" placeholder="Рубли"
                                :value="(size.priceCny * exchangeRate).toFixed(0)" readonly disabled />
                            </div>
                            <div class="col-md-2">
                              <label>Кол-во</label>
                              <input type="number" class="form-control" placeholder="Кол-во"
                                v-model.number="size.quantity" />
                            </div>
                            <div class="col-md-2">
                              <div class="form-check">
                                <input type="checkbox" class="form-check-input" :id="'onRequest_' + idx + '_' + sizeIdx"
                                  v-model="size.isOnRequest" />
                                <label class="form-check-label" :for="'onRequest_' + idx + '_' + sizeIdx">Под
                                  заказ</label>
                              </div>
                            </div>
                            <div class="col-md-2">
                              <button type="button" class="btn btn-outline-danger" @click="removeSize(idx, sizeIdx)">
                                Удалить
                              </button>
                            </div>
                          </div>
                          <button type="button" class="btn btn-sm btn-outline-secondary" @click="addSize(idx)">
                            Добавить размер
                          </button>
                        </div>

                        <div class="mb-3">
                          <label>Изображения</label>
                          <div class="d-flex flex-wrap gap-2 mb-2">
                            <div v-for="(img, imgIdx) in subItem.images" :key="imgIdx" class="position-relative"
                              style="width: 100px; height: 100px;">
                              <img :src="img" alt=""
                                style="width: 100%; height: 100%; object-fit: cover; border-radius: 4px;" />
                              <button type="button" class="btn btn-sm btn-danger position-absolute top-0 end-0"
                                style="border-radius: 50%; padding: 2px 6px;" @click="removeImage(idx, imgIdx)">
                                &times;
                              </button>
                            </div>
                          </div>
                          <input type="file" multiple accept="image/*" class="form-control"
                            @change="onImagesSelected($event, idx)" />
                        </div>
                      </div>
                    </div>
                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="addSubItem">
                      Добавить подтовар
                    </button>
                  </div>
                </form>
              </fieldset>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" :disabled="isSaving">
                Отмена
              </button>
              <button v-if="isEditing" type="button" class="btn btn-danger" @click="deleteItem"
                :disabled="isDeleting || isSaving">
                <span v-if="isDeleting" class="spinner-border spinner-border-sm me-1" role="status"
                  aria-hidden="true"></span>
                {{ isDeleting ? 'Удаление...' : 'Удалить товар' }}
              </button>
              <button type="submit" class="btn btn-primary" @click="saveItem" :disabled="isSaving">
                <span v-if="isSaving" class="spinner-border spinner-border-sm me-1" role="status"
                  aria-hidden="true"></span>
                {{ isSaving ? 'Сохранение...' : 'Сохранить' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, onMounted, nextTick, computed, onUnmounted } from 'vue';
import { api } from '../api';
import { useToastStore } from '../stores/toast';

const emit = defineEmits(['page-loaded']);

const countries = [
  { id: 'germany', ru: 'Германия', en: 'Germany' },
  { id: 'usa', ru: 'США', en: 'USA' },
  { id: 'sweden', ru: 'Швеция', en: 'Sweden' },
  { id: 'china', ru: 'Китай', en: 'China' },
  { id: 'russia', ru: 'Россия', en: 'Russia' },
  { id: 'uk', ru: 'Великобритания', en: 'United Kingdom' },
  { id: 'france', ru: 'Франция', en: 'France' },
  { id: 'italy', ru: 'Италия', en: 'Italy' },
  { id: 'spain', ru: 'Испания', en: 'Spain' },
  { id: 'japan', ru: 'Япония', en: 'Japan' },
  { id: 'south_korea', ru: 'Южная Корея', en: 'South Korea' },
  { id: 'vietnam', ru: 'Вьетнам', en: 'Vietnam' },
  { id: 'india', ru: 'Индия', en: 'India' },
  { id: 'brazil', ru: 'Бразилия', en: 'Brazil' },
  { id: 'canada', ru: 'Канада', en: 'Canada' },
  { id: 'australia', ru: 'Австралия', en: 'Australia' },
  { id: 'portugal', ru: 'Португалия', en: 'Portugal' },
  { id: 'netherlands', ru: 'Нидерланды', en: 'Netherlands' },
  { id: 'poland', ru: 'Польша', en: 'Poland' },
  { id: 'czech', ru: 'Чехия', en: 'Czech Republic' },
  { id: 'turkey', ru: 'Турция', en: 'Turkey' },
  { id: 'thailand', ru: 'Таиланд', en: 'Thailand' },
  { id: 'indonesia', ru: 'Индонезия', en: 'Indonesia' },
  { id: 'mexico', ru: 'Мексика', en: 'Mexico' },
  { id: 'argentina', ru: 'Аргентина', en: 'Argentina' },
  { id: 'south_africa', ru: 'ЮАР', en: 'South Africa' },
];

const categories = [
  { id: 'footwear', ru: 'Обувь', en: 'Footwear' },
  { id: 'clothing', ru: 'Одежда', en: 'Clothing' },
  { id: 'accessories', ru: 'Аксессуары', en: 'Accessories' },
  { id: 'equipment', ru: 'Экипировка', en: 'Equipment' },
];

const types = [
  // Обувь
  { id: 'sneakers', ru: 'Кроссовки', en: 'Sneakers', categoryId: 'footwear' },
  { id: 'cleats', ru: 'Бутсы', en: 'Cleats', categoryId: 'footwear' },
  { id: 'sandals', ru: 'Сандалии', en: 'Sandals', categoryId: 'footwear' },
  { id: 'boots', ru: 'Ботинки', en: 'Boots', categoryId: 'footwear' },
  { id: 'slippers', ru: 'Тапки', en: 'Slippers', categoryId: 'footwear' },
  // Одежда
  { id: 'tshirts', ru: 'Футболки', en: 'T-shirts', categoryId: 'clothing' },
  { id: 'downJackets', ru: 'Пуховики', en: 'Down jackets', categoryId: 'clothing' },
  { id: 'hoodies', ru: 'Худи', en: 'Hoodies', categoryId: 'clothing' },
  { id: 'pants', ru: 'Штаны', en: 'Pants', categoryId: 'clothing' },
  { id: 'shorts', ru: 'Шорты', en: 'Shorts', categoryId: 'clothing' },
  // Аксессуары
  { id: 'balaclavas', ru: 'Балаклавы', en: 'Balaclavas', categoryId: 'accessories' },
  { id: 'caps', ru: 'Кепки', en: 'Caps', categoryId: 'accessories' },
  { id: 'gloves', ru: 'Перчатки', en: 'Gloves', categoryId: 'accessories' },
  { id: 'bags', ru: 'Сумки', en: 'Bags', categoryId: 'accessories' },
  { id: 'scarves', ru: 'Шарфы', en: 'Scarves', categoryId: 'accessories' },
  // Экипировка
  { id: 'shin_guards', ru: 'Щитки', en: 'Shin guards', categoryId: 'equipment' },
  { id: 'helmets', ru: 'Шлемы', en: 'Helmets', categoryId: 'equipment' },
  { id: 'protectors', ru: 'Защита', en: 'Protectors', categoryId: 'equipment' },
];

const brands = [
  'Adidas', 'Nike', 'The North Face', 'Craft', 'Puma', 'Reebok',
  'Under Armour', 'New Balance', 'Asics', 'Mizuno', 'Joma', 'Umbro',
  'Kappa', 'Lotto', 'Diadora', 'Fila', 'Wilson', 'Select', 'Molten',
  'Spalding', 'Macron', 'Errea', 'Hummel', 'Le Coq Sportif', 'Castore',
  'Patrick', 'Uhlsport', 'Kipsta', 'Decathlon', 'Demix', 'Airness'
];

const colors = [
  { code: 'blue', ru: 'Синий', en: 'Blue', hex: '#0000FF' },
  { code: 'yellow', ru: 'Желтый', en: 'Yellow', hex: '#FFFF00' },
  { code: 'black', ru: 'Черный', en: 'Black', hex: '#000000' },
  { code: 'red', ru: 'Красный', en: 'Red', hex: '#FF0000' },
  { code: 'darkblue', ru: 'Темно-синий', en: 'Dark Blue', hex: '#00008B' },
  { code: 'gray', ru: 'Серый', en: 'Gray', hex: '#808080' },
  { code: 'green', ru: 'Зеленый', en: 'Green', hex: '#008000' },
];

const tags = [
  { id: 'hit', ru: 'Хит', en: 'Hit' },
  { id: 'popular', ru: 'Популярное', en: 'Popular' },
  { id: 'new', ru: 'Новое', en: 'New' },
  { id: 'limited', ru: 'Ограничено', en: 'Limited' },
  { id: 'gift', ru: 'Идея для подарка', en: 'Gift idea' },
];

const materials = [
  { id: 'textile', ru: 'Текстиль', en: 'Textile' },
  { id: 'rubber', ru: 'Резина', en: 'Rubber' },
  { id: 'synthetic', ru: 'Синтетика', en: 'Synthetic' },
  { id: 'polymer', ru: 'Полимер', en: 'Polymer' },
  { id: 'leather', ru: 'Кожа', en: 'Leather' },
];

// ---------- Основные данные ----------
const items = ref([]);
const modalInstance = ref(null);
const modal = ref(null);
const isEditing = ref(false);
const editingId = ref(null);

const isLoading = ref(false);
const isLoadingItem = ref(false);
const isSaving = ref(false);
const isDeleting = ref(false);

const form = ref({
  type: { ru: '', en: '' },
  title: { ru: '', en: '' },
  description: { ru: '', en: '' },
  brand: '',
  country: { ru: '', en: '' },
  category: { ru: '', en: '' },
  structure: [],
  items: []
});

const exchangeRate = ref(0);
const exchangeRateInterval = ref(null);

// ---------- Вычисляемые свойства для селекторов ----------
const selectedCountry = computed({
  get: () => {
    if (!form.value.country?.ru && !form.value.country?.en) return '';
    const found = countries.find(c => c.ru === form.value.country.ru && c.en === form.value.country.en);
    return found ? found.id : '';
  },
  set: (id) => {
    const found = countries.find(c => c.id === id);
    if (found) {
      form.value.country = { ru: found.ru, en: found.en };
    } else {
      form.value.country = { ru: '', en: '' };
    }
  }
});

const selectedCategory = computed({
  get: () => {
    if (!form.value.category?.ru && !form.value.category?.en) return '';
    const found = categories.find(c => c.ru === form.value.category.ru && c.en === form.value.category.en);
    return found ? found.id : '';
  },
  set: (id) => {
    const found = categories.find(c => c.id === id);
    if (found) {
      form.value.category = { ru: found.ru, en: found.en };
      if (form.value.type?.ru || form.value.type?.en) {
        const typeFound = types.find(t => t.ru === form.value.type.ru && t.en === form.value.type.en);
        if (typeFound && typeFound.categoryId !== id) {
          form.value.type = { ru: '', en: '' };
        }
      }
    } else {
      form.value.category = { ru: '', en: '' };
    }
  }
});

const filteredTypes = computed(() => {
  if (!selectedCategory.value) return types;
  return types.filter(t => t.categoryId === selectedCategory.value);
});

const selectedType = computed({
  get: () => {
    if (!form.value.type?.ru && !form.value.type?.en) return '';
    const found = types.find(t => t.ru === form.value.type.ru && t.en === form.value.type.en);
    return found ? found.id : '';
  },
  set: (id) => {
    const found = types.find(t => t.id === id);
    if (found) {
      form.value.type = { ru: found.ru, en: found.en };
    } else {
      form.value.type = { ru: '', en: '' };
    }
  }
});

// ---------- Вспомогательные методы для работы с предопределёнными данными ----------
function getColorCode(colorObj) {
  if (!colorObj?.ru && !colorObj?.en) return '';
  const found = colors.find(c => c.ru === colorObj.ru && c.en === colorObj.en);
  return found ? found.code : '';
}

function getColorHex(colorObj) {
  const found = colors.find(c => c.ru === colorObj.ru && c.en === colorObj.en);
  return found ? found.hex : '#cccccc';
}

function updateColor(itemIndex, colorIndex, code) {
  if (!code) return;
  const found = colors.find(c => c.code === code);
  if (found) {
    form.value.items[itemIndex].color[colorIndex] = { ru: found.ru, en: found.en };
  }
}

function getMaterialId(material) {
  if (!material.name?.ru && !material.name?.en) return '';
  const found = materials.find(m => m.ru === material.name.ru && m.en === material.name.en);
  return found ? found.id : '';
}

function updateMaterial(index, materialId) {
  const found = materials.find(m => m.id === materialId);
  if (found) {
    form.value.structure[index].name = { ru: found.ru, en: found.en };
  }
}

async function fetchExchangeRate() {
  try {
    const response = await api.getExchangeRate();
    exchangeRate.value = response.data.rate;
  } catch (error) {
    console.error('Ошибка загрузки курса валют:', error);
    exchangeRate.value = 11.5; // запасной курс
  }
}

function startExchangeRateUpdater() {
  fetchExchangeRate();
  exchangeRateInterval.value = setInterval(fetchExchangeRate, 60000);
}

function stopExchangeRateUpdater() {
  if (exchangeRateInterval.value) {
    clearInterval(exchangeRateInterval.value);
  }
}

// ---------- Нормализация данных после загрузки ----------
function normalizeFormData() {
  // Страна
  if (form.value.country?.ru || form.value.country?.en) {
    const found = countries.find(c => c.ru === form.value.country.ru && c.en === form.value.country.en);
    if (found) form.value.country = { ru: found.ru, en: found.en };
  }
  // Тип
  if (form.value.type?.ru || form.value.type?.en) {
    const found = types.find(t => t.ru === form.value.type.ru && t.en === form.value.type.en);
    if (found) form.value.type = { ru: found.ru, en: found.en };
  }
  // Состав
  form.value.structure.forEach(material => {
    if (material.name?.ru || material.name?.en) {
      const found = materials.find(m => m.ru === material.name.ru && m.en === material.name.en);
      if (found) material.name = { ru: found.ru, en: found.en };
    }
  });
  // Подтовары
  form.value.items.forEach(item => {
    item.color = item.color.map(c => {
      const found = colors.find(col => col.ru === c.ru && col.en === c.en);
      return found ? { ru: found.ru, en: found.en } : c;
    });
    item.tags = item.tags.map(t => {
      const found = tags.find(tag => tag.ru === t.ru && tag.en === t.en);
      return found ? { ...found } : t;
    });
  });
}

// ---------- Загрузка товаров ----------
async function loadItems() {
  isLoading.value = true;
  try {
    const response = await api.getAdminItems();
    items.value = response.data;
  } catch (error) {
    console.error('Ошибка загрузки товаров:', error);
  } finally {
    isLoading.value = false;
  }
}

// ---------- Модальное окно ----------
function openCreateModal() {
  isEditing.value = false;
  editingId.value = null;
  resetForm();
  addSubItem();
  showModal();
}

async function openEditModal(item) {
  isEditing.value = true;
  editingId.value = item.id;
  isLoadingItem.value = true;
  showModal();

  try {
    const response = await api.getAdminItem(item.id);
    const fullItem = response.data;
    form.value = JSON.parse(JSON.stringify(fullItem));
    normalizeFormData();

    if (exchangeRate.value > 0) {
      form.value.items.forEach(subItem => {
        subItem.sizes.forEach(size => {
          size.priceCny = Math.round(size.price / exchangeRate.value);
        });
      });
    }
  } catch (error) {
    console.error('Не удалось загрузить товар для редактирования', error);
  } finally {
    isLoadingItem.value = false;
  }
}

function showModal() {
  if (modalInstance.value) {
    modalInstance.value.show();
  } else {
    nextTick(() => {
      modalInstance.value = new window.bootstrap.Modal(modal.value);
      modalInstance.value.show();
    });
  }
}

function hideModal() {
  if (modalInstance.value) {
    modalInstance.value.hide();
  }
}

function resetForm() {
  form.value = {
    type: { ru: '', en: '' },
    title: { ru: '', en: '' },
    description: { ru: '', en: '' },
    brand: '',
    country: { ru: '', en: '' },
    category: { ru: '', en: '' },
    structure: [],
    items: []
  };
}

// ---------- Работа с подтоварами ----------
function addSubItem() {
  form.value.items.push({
    uniqueId: null,
    images: [],
    color: [],
    tags: [],
    sizes: []
  });
}

function removeSubItem(index) {
  form.value.items.splice(index, 1);
}

function addColor(itemIndex) {
  form.value.items[itemIndex].color.push({ ru: '', en: '' });
}

function removeColor(itemIndex, colorIndex) {
  form.value.items[itemIndex].color.splice(colorIndex, 1);
}

function addSize(itemIndex) {
  form.value.items[itemIndex].sizes.push({
    size: '',
    priceCny: 0,
    price: 0,
    quantity: 0,
    isOnRequest: false
  });
}

function removeSize(itemIndex, sizeIndex) {
  form.value.items[itemIndex].sizes.splice(sizeIndex, 1);
}

function addMaterial() {
  form.value.structure.push({
    name: { ru: '', en: '' },
    percent: 0
  });
}

function removeMaterial(index) {
  form.value.structure.splice(index, 1);
}

function onImagesSelected(event, itemIndex) {
  const files = Array.from(event.target.files);
  files.forEach(file => {
    const reader = new FileReader();
    reader.onload = (e) => {
      form.value.items[itemIndex].images.push(e.target.result);
    };
    reader.readAsDataURL(file);
  });
}

function removeImage(itemIndex, imageIndex) {
  form.value.items[itemIndex].images.splice(imageIndex, 1);
}

function isTagSelected(tagsArray, tag) {
  return tagsArray.some(t => t.id === tag.id);
}

function toggleTag(subItem, tag) {
  const index = subItem.tags.findIndex(t => t.id === tag.id);
  if (index === -1) {
    subItem.tags.push(tag);
  } else {
    subItem.tags.splice(index, 1);
  }
}

// ---------- Сохранение ----------
async function saveItem() {
  isSaving.value = true;
  try {
    if (exchangeRate.value > 0) {
      form.value.items.forEach(subItem => {
        subItem.sizes.forEach(size => {
          if (size.priceCny !== undefined) {
            size.price = Math.round(size.priceCny * exchangeRate.value);
          }
        });
      });
    }

    if (isEditing.value) {
      const productId = editingId.value;
      let maxNumber = 0;
      form.value.items.forEach(item => {
        if (item.uniqueId && item.uniqueId.startsWith(productId + '-')) {
          const num = extractNumberFromUniqueId(item.uniqueId);
          if (num > maxNumber) maxNumber = num;
        }
      });

      const updatedItems = form.value.items.map(item => {
        if (!item.uniqueId) {
          maxNumber++;
          return { ...item, uniqueId: `${productId}-${maxNumber}` };
        }
        return item;
      });

      const payload = {
        ...form.value,
        items: updatedItems
      };

      await api.updateMainItem(productId, payload);
    } else {
      const mainPayload = {
        type: form.value.type,
        title: form.value.title,
        description: form.value.description,
        brand: form.value.brand,
        country: form.value.country,
        category: form.value.category,
        structure: form.value.structure,
        createdAt: new Date().toISOString(),
        items: [],
      };

      const createResponse = await api.createMainItem(mainPayload);
      const newId = createResponse.data.id;

      const newItems = form.value.items.map((item, index) => ({
        ...item,
        uniqueId: `${newId}-${index + 1}`
      }));

      const updatePayload = {
        ...form.value,
        id: newId,
        items: newItems
      };

      await api.updateMainItem(newId, updatePayload);
    }

    await loadItems();
    hideModal();
  } catch (error) {
    console.error('Ошибка сохранения товара:', error);
    useToastStore().error('Не удалось сохранить товар');
  } finally {
    isSaving.value = false;
  }
}

async function deleteItem() {
  if (!editingId.value) return;
  if (!confirm('Вы уверены, что хотите удалить этот товар?')) return;

  isDeleting.value = true;
  try {
    await api.deleteMainItem(editingId.value);
    await loadItems();
    hideModal();
  } catch (error) {
    console.error('Ошибка удаления:', error);
    useToastStore().error('Не удалось удалить товар');
  } finally {
    isDeleting.value = false;
  }
}

function extractNumberFromUniqueId(uniqueId) {
  if (!uniqueId) return 0;
  const parts = uniqueId.split('-');
  return parts.length > 1 ? parseInt(parts[1], 10) : 0;
}

onMounted(async () => {
  await loadItems();
  startExchangeRateUpdater();
  emit('page-loaded', true);
});

onUnmounted(() => {
  stopExchangeRateUpdater();
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/admin-products';
</style>