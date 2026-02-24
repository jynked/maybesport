<template>
  <main>
    <div class="container-fluid py-4">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление товарами</h1>
        <button class="btn btn-primary" @click="openCreateModal">
          <i class="bi bi-plus-lg"></i> Создать товар
        </button>
      </div>

      <!-- Таблица товаров -->
      <div class="table-responsive">
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
            <tr v-for="item in items" :key="item.uniqueId">
              <td>{{ item.uniqueId }}</td>
              <td>
                <img :src="item.images[0]" alt="" style="width: 50px; height: 50px; object-fit: cover;" />
              </td>
              <td>{{ item.title.ru }}</td>
              <td>{{ item.title.en }}</td>
              <td>
                <button class="btn btn-sm btn-outline-info" @click="openEditModal(item)">
                  Характеристики
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Модальное окно редактирования/создания -->
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
              <form @submit.prevent="saveItem">
                <!-- Основная информация -->
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
                  <div class="col-md-4">
                    <label class="form-label">Категория (RU)</label>
                    <input type="text" class="form-control" v-model="form.category.ru" required />
                  </div>
                  <div class="col-md-4">
                    <label class="form-label">Категория (EN)</label>
                    <input type="text" class="form-control" v-model="form.category.en" required />
                  </div>
                  <div class="col-md-4">
                    <label class="form-label">Бренд</label>
                    <input type="text" class="form-control" v-model="form.brand" required />
                  </div>
                </div>

                <div class="row mb-3">
                  <div class="col-md-4">
                    <label class="form-label">Страна (RU)</label>
                    <input type="text" class="form-control" v-model="form.country.ru" required />
                  </div>
                  <div class="col-md-4">
                    <label class="form-label">Страна (EN)</label>
                    <input type="text" class="form-control" v-model="form.country.en" required />
                  </div>
                  <div class="col-md-4">
                    <label class="form-label">Тип (RU)</label>
                    <input type="text" class="form-control" v-model="form.type.ru" required />
                  </div>
                  <div class="col-md-4 mt-2">
                    <label class="form-label">Тип (EN)</label>
                    <input type="text" class="form-control" v-model="form.type.en" required />
                  </div>
                </div>

                <!-- Цвета (динамический список) -->
                <div class="mb-3">
                  <label class="form-label">Цвета</label>
                  <div v-for="(color, index) in form.color" :key="index" class="row mb-2">
                    <div class="col-md-5">
                      <input type="text" class="form-control" placeholder="Цвет (RU)" v-model="color.ru" />
                    </div>
                    <div class="col-md-5">
                      <input type="text" class="form-control" placeholder="Цвет (EN)" v-model="color.en" />
                    </div>
                    <div class="col-md-2">
                      <button type="button" class="btn btn-outline-danger" @click="removeColor(index)">
                        Удалить
                      </button>
                    </div>
                  </div>
                  <button type="button" class="btn btn-sm btn-outline-secondary" @click="addColor">
                    Добавить цвет
                  </button>
                </div>

                <!-- Теги (динамический список) -->
                <div class="mb-3">
                  <label class="form-label">Теги</label>
                  <div v-for="(tag, index) in form.tags" :key="index" class="row mb-2">
                    <div class="col-md-5">
                      <input type="text" class="form-control" placeholder="Тег (RU)" v-model="tag.ru" />
                    </div>
                    <div class="col-md-5">
                      <input type="text" class="form-control" placeholder="Тег (EN)" v-model="tag.en" />
                    </div>
                    <div class="col-md-2">
                      <button type="button" class="btn btn-outline-danger" @click="removeTag(index)">
                        Удалить
                      </button>
                    </div>
                  </div>
                  <button type="button" class="btn btn-sm btn-outline-secondary" @click="addTag">
                    Добавить тег
                  </button>
                </div>

                <!-- Размеры (сложная структура) -->
                <div class="mb-3">
                  <label class="form-label">Размеры</label>
                  <div v-for="(size, index) in form.sizes" :key="index" class="row mb-2 align-items-center">
                    <div class="col-md-2" style="margin-top: auto;">
                      <input type="text" class="form-control" placeholder="Размер" v-model="size.size"/>
                    </div>
                    <div class="col-md-2">
                      <label for="new-item-price">Цена</label>
                      <input type="number" id="new-item-price" class="form-control" placeholder="Цена" v-model.number="size.price" />
                    </div>
                    <div class="col-md-2">
                      <label for="new-item-quantity">Кол-во</label>
                      <input type="number" id="new-item-quantity" class="form-control" placeholder="Кол-во" v-model.number="size.quantity" />
                    </div>
                    <div class="col-md-2">
                      <div class="form-check">
                        <input type="checkbox" class="form-check-input" :id="'onRequest_' + index"
                          v-model="size.isOnRequest" />
                        <label class="form-check-label" :for="'onRequest_' + index">Под заказ</label>
                      </div>
                    </div>
                    <div class="col-md-2">
                      <button type="button" class="btn btn-outline-danger" @click="removeSize(index)">
                        Удалить
                      </button>
                    </div>
                  </div>
                  <button type="button" class="btn btn-sm btn-outline-secondary" @click="addSize">
                    Добавить размер
                  </button>
                </div>

                <!-- Состав (structure) -->
                <div class="mb-3">
                  <label class="form-label">Состав</label>
                  <div v-for="(material, index) in form.structure" :key="index" class="row mb-2">
                    <div class="col-md-4">
                      <input type="text" class="form-control" placeholder="Название (RU)" v-model="material.name.ru" />
                    </div>
                    <div class="col-md-4">
                      <input type="text" class="form-control" placeholder="Название (EN)" v-model="material.name.en" />
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

                <!-- Изображения -->
                <div class="mb-3">
                  <label class="form-label">Изображения</label>
                  <div class="d-flex flex-wrap gap-2 mb-2">
                    <div v-for="(img, idx) in form.images" :key="idx" class="position-relative"
                      style="width: 100px; height: 100px;">
                      <img :src="img" alt=""
                        style="width: 100%; height: 100%; object-fit: cover; border-radius: 4px;" />
                      <button type="button" class="btn btn-sm btn-danger position-absolute top-0 end-0"
                        style="border-radius: 50%; padding: 2px 6px;" @click="removeImage(idx)">
                        &times;
                      </button>
                    </div>
                  </div>
                  <input type="file" multiple accept="image/*" class="form-control" @change="onImagesSelected" />
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">
                Отмена
              </button>
              <button v-if="isEditing" type="button" class="btn btn-danger" @click="deleteItem">
                Удалить товар
              </button>
              <button type="submit" class="btn btn-primary" @click="saveItem">
                Сохранить
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

// Данные
const items = ref([]);
const modalInstance = ref(null);
const modal = ref(null);
const isEditing = ref(false);
const currentUniqueId = ref(null);

// Форма
const form = ref({
  title: { ru: '', en: '' },
  description: { ru: '', en: '' },
  category: { ru: '', en: '' },
  brand: '',
  country: { ru: '', en: '' },
  type: { ru: '', en: '' },
  color: [],
  tags: [],
  sizes: [],
  structure: [],
  images: []
});

// Загрузка товаров
async function loadItems() {
  try {
    const response = await api.getItems();
    items.value = response.data.items; 
  } catch (error) {
    console.error('Ошибка загрузки товаров:', error);
  }
}

// Модальное окно
function openCreateModal() {
  isEditing.value = false;
  currentUniqueId.value = null;
  resetForm();
  showModal();
}

function openEditModal(item) {
  isEditing.value = true;
  currentUniqueId.value = item.uniqueId;
  currentParentId.value = item.id; // сохраняем родительский ID
  form.value = JSON.parse(JSON.stringify(item));
  showModal();
}

function showModal() {
  if (modalInstance.value) {
    modalInstance.value.show();
  } else {
    // Инициализация при первом открытии
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

// Сброс формы (для создания)
function resetForm() {
  form.value = {
    title: { ru: '', en: '' },
    description: { ru: '', en: '' },
    category: { ru: '', en: '' },
    brand: '',
    country: { ru: '', en: '' },
    type: { ru: '', en: '' },
    color: [],
    tags: [],
    sizes: [],
    structure: [],
    images: []
  };
}

// Добавление/удаление цветов
function addColor() {
  form.value.color.push({ ru: '', en: '' });
}
function removeColor(index) {
  form.value.color.splice(index, 1);
}

// Добавление/удаление тегов
function addTag() {
  form.value.tags.push({ ru: '', en: '' });
}
function removeTag(index) {
  form.value.tags.splice(index, 1);
}

// Добавление/удаление размеров
function addSize() {
  form.value.sizes.push({
    size: '',
    price: 0,
    quantity: 0,
    isOnRequest: false
  });
}
function removeSize(index) {
  form.value.sizes.splice(index, 1);
}

// Добавление/удаление материалов состава
function addMaterial() {
  form.value.structure.push({
    name: { ru: '', en: '' },
    percent: 0
  });
}
function removeMaterial(index) {
  form.value.structure.splice(index, 1);
}

// Работа с изображениями
function onImagesSelected(event) {
  const files = Array.from(event.target.files);
  // Здесь можно либо сразу загрузить файлы на сервер и получить URL,
  // либо временно хранить File объекты для последующей отправки через FormData.
  // Для упрощения добавим локальные URL для предпросмотра.
  files.forEach(file => {
    const reader = new FileReader();
    reader.onload = (e) => {
      form.value.images.push(e.target.result);
    };
    reader.readAsDataURL(file);
  });
  // Также сохраняем сами файлы для отправки, если нужно.
  // Можно хранить их в отдельном reactive-массиве.
}

function removeImage(index) {
  form.value.images.splice(index, 1);
}

async function saveItem() {
  try {
    const newUniqueId = 'item-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);

    const newItem = {
      title: form.value.title,
      description: form.value.description,
      category: form.value.category,
      brand: form.value.brand,
      country: form.value.country,
      type: form.value.type,
      structure: form.value.structure,
      createdAt: new Date().toISOString(),
      items: [
        {
          uniqueId: newUniqueId,
          images: form.value.images,
          color: form.value.color,
          tags: form.value.tags,
          sizes: form.value.sizes.map(s => ({
            size: s.size,
            price: s.price,
            quantity: s.quantity,
            isOnRequest: s.isOnRequest
          }))
        }
      ]
    };

    if (!isEditing.value) {
      await api.createItem(newItem);
    } else {
      // Пока редактирование не реализовано
      alert('Редактирование временно отключено');
      return;
    }

    await loadItems();
    hideModal();
  } catch (error) {
    console.error('Ошибка сохранения товара:', error);
    if (error.response) {
      alert(`Ошибка ${error.response.status}: ${JSON.stringify(error.response.data)}`);
    } else {
      alert('Ошибка при сохранении');
    }
  }
}

async function deleteItem() {
  if (!currentParentId.value) return;
  if (!confirm('Вы уверены, что хотите удалить этот товар?')) return;

  try {
    await api.deleteItem(currentParentId.value);
    await loadItems();
    hideModal();
  } catch (error) {
    console.error('Ошибка удаления:', error);
    alert('Не удалось удалить товар');
  }
}

onMounted(() => {
  loadItems();
});
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/admin-products';
</style>