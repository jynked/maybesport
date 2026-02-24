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
            <tr v-for="item in items" :key="item.id">
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
                    <div class="col-md-6">
                      <label class="form-label">Категория (RU)</label>
                      <input type="text" class="form-control" v-model="form.category.ru" required />
                    </div>
                    <div class="col-md-6">
                      <label class="form-label">Категория (EN)</label>
                      <input type="text" class="form-control" v-model="form.category.en" required />
                    </div>
                  </div>
                  <div class="row mb-3">
                    <div class="col-md-3">
                      <label class="form-label">Страна (RU)</label>
                      <input type="text" class="form-control" v-model="form.country.ru" required />
                    </div>
                    <div class="col-md-3">
                      <label class="form-label">Страна (EN)</label>
                      <input type="text" class="form-control" v-model="form.country.en" required />
                    </div>
                    <div class="col-md-3">
                      <label class="form-label">Тип (RU)</label>
                      <input type="text" class="form-control" v-model="form.type.ru" required />
                    </div>
                    <div class="col-md-3">
                      <label class="form-label">Тип (EN)</label>
                      <input type="text" class="form-control" v-model="form.type.en" required />
                    </div>
                  </div>
                  <div class="col-md-12 mb-4">
                    <label class="form-label">Бренд</label>
                    <input type="text" class="form-control" v-model="form.brand" required />
                  </div>

                  <div class="mb-3 form-row-custom">
                    <label class="form-label">Состав</label>
                    <div v-for="(material, index) in form.structure" :key="index" class="row mb-2">
                      <div class="col-md-4">
                        <input type="text" class="form-control" placeholder="Название (RU)"
                          v-model="material.name.ru" />
                      </div>
                      <div class="col-md-4">
                        <input type="text" class="form-control" placeholder="Название (EN)"
                          v-model="material.name.en" />
                      </div>
                      <div class="col-md-2">
                        <input type="number" class="form-control" placeholder="%" v-model.number="material.percent" />
                      </div>
                      <div class="col-md-2">
                        <button type="button" class="btn btn-outline-danger" @click="removeMaterial(index)"
                          :disabled="isSaving">
                          Удалить
                        </button>
                      </div>
                    </div>
                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="addMaterial"
                      :disabled="isSaving">
                      Добавить материал
                    </button>
                  </div>

                  <div class="mb-4">
                    <h6>Варианты товара (подтовары)</h6>
                    <div v-for="(subItem, idx) in form.items" :key="idx" class="card mb-3">
                      <div class="card-body">
                        <div class="d-flex justify-content-between align-items-center mb-2">
                          <h6>Подтовар #{{ idx + 1 }}</h6>
                          <button type="button" class="btn btn-sm btn-danger" @click="removeSubItem(idx)"
                            :disabled="isSaving">
                            Удалить подтовар
                          </button>
                        </div>

                        <div class="mb-3 form-row-custom">
                          <label>Цвета</label>
                          <div v-for="(color, colorIdx) in subItem.color" :key="colorIdx" class="row mb-2">
                            <div class="col-md-5">
                              <input type="text" class="form-control" placeholder="Цвет (RU)" v-model="color.ru" />
                            </div>
                            <div class="col-md-5">
                              <input type="text" class="form-control" placeholder="Цвет (EN)" v-model="color.en" />
                            </div>
                            <div class="col-md-2">
                              <button type="button" class="btn btn-outline-danger" @click="removeColor(idx, colorIdx)"
                                :disabled="isSaving">
                                Удалить
                              </button>
                            </div>
                          </div>
                          <button type="button" class="btn btn-sm btn-outline-secondary" @click="addColor(idx)"
                            :disabled="isSaving">
                            Добавить цвет
                          </button>
                        </div>

                        <div class="mb-3 form-row-custom">
                          <label>Теги</label>
                          <div v-for="(tag, tagIdx) in subItem.tags" :key="tagIdx" class="row mb-2">
                            <div class="col-md-5">
                              <input type="text" class="form-control" placeholder="Тег (RU)" v-model="tag.ru" />
                            </div>
                            <div class="col-md-5">
                              <input type="text" class="form-control" placeholder="Тег (EN)" v-model="tag.en" />
                            </div>
                            <div class="col-md-2">
                              <button type="button" class="btn btn-outline-danger" @click="removeTag(idx, tagIdx)"
                                :disabled="isSaving">
                                Удалить
                              </button>
                            </div>
                          </div>
                          <button type="button" class="btn btn-sm btn-outline-secondary" @click="addTag(idx)"
                            :disabled="isSaving">
                            Добавить тег
                          </button>
                        </div>

                        <div class="mb-3 form-row-custom">
                          <label>Размеры</label>
                          <div v-for="(size, sizeIdx) in subItem.sizes" :key="sizeIdx"
                            class="row mb-2 align-items-center">
                            <div class="col-md-2">
                              <input type="text" class="form-control" placeholder="Размер" v-model="size.size" />
                            </div>
                            <div class="col-md-2">
                              <label>Цена</label>
                              <input type="number" class="form-control" placeholder="Цена"
                                v-model.number="size.price" />
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
                              <button type="button" class="btn btn-outline-danger" @click="removeSize(idx, sizeIdx)"
                                :disabled="isSaving">
                                Удалить
                              </button>
                            </div>
                          </div>
                          <button type="button" class="btn btn-sm btn-outline-secondary" @click="addSize(idx)"
                            :disabled="isSaving">
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
                                style="border-radius: 50%; padding: 2px 6px;" @click="removeImage(idx, imgIdx)"
                                :disabled="isSaving">
                                &times;
                              </button>
                            </div>
                          </div>
                          <input type="file" multiple accept="image/*" class="form-control"
                            @change="onImagesSelected($event, idx)" :disabled="isSaving" />
                        </div>
                      </div>
                    </div>
                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="addSubItem"
                      :disabled="isSaving">
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
import { ref, onMounted, nextTick } from 'vue';
import { api } from '../api';

const emit = defineEmits(['page-loaded']);

// Данные
const items = ref([]);
const modalInstance = ref(null);
const modal = ref(null);
const isEditing = ref(false);
const editingId = ref(null);

// Состояния загрузки
const isLoading = ref(false);
const isLoadingItem = ref(false);
const isSaving = ref(false);
const isDeleting = ref(false);

// Форма
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

// Загрузка товаров
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

// Модальное окно
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

// Сброс формы
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

// Работа с подтоварами
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


function addTag(itemIndex) {
  form.value.items[itemIndex].tags.push({ ru: '', en: '' });
}
function removeTag(itemIndex, tagIndex) {
  form.value.items[itemIndex].tags.splice(tagIndex, 1);
}


function addSize(itemIndex) {
  form.value.items[itemIndex].sizes.push({
    size: '',
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


async function saveItem() {
  isSaving.value = true;
  try {
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
    alert('Не удалось сохранить товар');
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
    alert('Не удалось удалить товар');
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
  emit('page-loaded', true);
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/admin-products';
</style>