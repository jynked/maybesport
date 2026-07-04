<template>
  <main>
    <div class="admin-orders">
      <h1>Управление заказами</h1>

      <div class="filters">
        <input type="text" v-model="filters.userEmail" placeholder="Email пользователя" @input="debounceSearch" />
        <select v-model="filters.status">
          <option value="">Все статусы</option>
          <option value="created">Создан</option>
          <option value="processing">Обработка</option>
          <option value="shipped">Отправлен</option>
          <option value="delivered">Доставлен</option>
          <option value="received">Получен</option>
          <option value="cancelled">Отменён</option>
        </select>
        <button @click="loadOrders">Применить</button>
      </div>

      <div v-if="loading" class="loading">Загрузка...</div>
      <div v-else-if="orders.length === 0" class="empty">Нет заказов</div>
      <div v-else class="table-container-custom">
        <table class="orders-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Пользователь</th>
              <th>Сумма</th>
              <th>Статус</th>
              <th>Адрес</th>
              <th>Дата</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="order in orders" :key="order.id">
              <td>{{ order.id }}</td>
              <td>{{ order.user_email }} (ID: {{ order.user_id }})</td>
              <td>{{ order.total_amount.toLocaleString() }} ₽</td>
              <td><span :class="'status-' + order.status">{{ getStatusText(order.status) }}</span></td>
              <td>{{ order.delivery_address || '—' }}</td>
              <td>{{ formatDate(order.created_at) }}</td>
              <td><button @click="openManageItems(order)" class="btn-change">Управление</button></td>
             </tr>
          </tbody>
        </table>
        <div class="pagination" v-if="totalPages > 1">
          <button :disabled="page === 1" @click="changePage(page-1)">«</button>
          <span>Страница {{ page }} из {{ totalPages }}</span>
          <button :disabled="page === totalPages" @click="changePage(page+1)">»</button>
        </div>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="selectedOrder" class="status-modal" @click.self="closeModal">
        <div class="modal-content" style="max-width: 900px; width: 90%;">
          <div class="modal-header">
            <h2>Товары в заказе #{{ selectedOrder.id }}</h2>
            <button class="close-button" @click="closeModal">×</button>
          </div>
          <div v-if="loadingItems" class="loading">Загрузка товаров...</div>
          <div v-else class="order-items-list">
            <div v-for="item in orderItems" :key="item.id" class="order-item-row">
              <img :src="item.image" class="item-img" />
              <div class="item-info">
                <div><strong>{{ item.title.ru }} / {{ item.title.en }}</strong></div>
                <div>Размер: {{ item.size }}, кол-во: {{ item.quantity }}, цена: {{ item.price.toLocaleString() }} ₽</div>
                <div>Текущий статус: <span :class="'status-' + item.status">{{ getStatusText(item.status) }}</span></div>
                <div v-if="item.statusHistory && item.statusHistory.length" class="status-history-short">
                  <button @click="toggleItemHistory(item.id)" class="history-toggle">История изменений</button>
                  <div v-if="expandedHistory === item.id" class="history-list">
                    <div v-for="(h, idx) in item.statusHistory" :key="idx" class="history-item">
                      {{ getStatusText(h.status) }} — {{ formatDate(h.timestamp) }}<span v-if="h.description"> ({{ h.description }})</span>
                    </div>
                  </div>
                </div>
              </div>
              <div class="item-status-control">
                <select v-model="item.newStatus">
                  <option value="created">Создан</option>
                  <option value="processing">Обработка</option>
                  <option value="shipped">Отправлен</option>
                  <option value="delivered">Доставлен</option>
                  <option value="received">Получен</option>
                  <option value="cancelled">Отменён</option>
                </select>
                <input type="text" v-model="item.comment" placeholder="Комментарий" />
                <button @click="updateItemStatus(item)" :disabled="item.updating">Сохранить</button>
              </div>
            </div>
          </div>
          <div class="bulk-actions">
            <label>Применить ко всем:</label>
            <select v-model="bulkStatus">
              <option value="">-- не менять --</option>
              <option value="created">Создан</option>
              <option value="processing">Обработка</option>
              <option value="shipped">Отправлен</option>
              <option value="delivered">Доставлен</option>
              <option value="received">Получен</option>
              <option value="cancelled">Отменён</option>
            </select>
            <input type="text" v-model="bulkComment" placeholder="Комментарий для всех" />
            <button @click="applyBulkStatus" :disabled="!bulkStatus">Применить</button>
          </div>
        </div>
      </div>
    </Transition>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../api'

const emit = defineEmits(['page-loaded'])

const orders = ref([])
const loading = ref(false)
const page = ref(1)
const totalPages = ref(1)
const filters = ref({ userEmail: '', status: '' })
let searchTimeout = null

const selectedOrder = ref(null)
const orderItems = ref([])
const loadingItems = ref(false)
const bulkStatus = ref('')
const bulkComment = ref('')
const expandedHistory = ref(null)

const getStatusText = (status) => {
  const map = {
    created: 'Создан',
    processing: 'Обработка',
    shipped: 'Отправлен',
    delivered: 'Доставлен',
    received: 'Получен',
    cancelled: 'Отменён'
  }
  return map[status] || status
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleString()
}

const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      limit: 20,
      status: filters.value.status || undefined,
      user_email: filters.value.userEmail || undefined
    }
    const resp = await api.getAdminOrders(params)
    orders.value = resp.data.orders
    totalPages.value = resp.data.totalPages
    emit('page-loaded', true)
  } catch (err) {
    console.error(err)
    emit('page-loaded', false)
  } finally {
    loading.value = false
  }
}

const changePage = (newPage) => {
  page.value = newPage
  loadOrders()
}

const debounceSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    loadOrders()
  }, 500)
}

const openManageItems = async (order) => {
  selectedOrder.value = order
  loadingItems.value = true
  try {
    const resp = await api.getAdminOrderItems(order.id)
    orderItems.value = resp.data.map(item => ({
      ...item,
      newStatus: item.status,
      comment: '',
      updating: false
    }))
  } catch (err) {
    console.error(err)
  } finally {
    loadingItems.value = false
  }
}

const updateItemStatus = async (item) => {
  if (item.updating) return
  item.updating = true
  try {
    await api.updateOrderItemStatus(item.id, item.newStatus, item.comment)
    item.status = item.newStatus
    item.comment = ''
    await loadOrders()
  } catch (err) {
    console.error(err)
  } finally {
    item.updating = false
  }
}

const applyBulkStatus = async () => {
  if (!bulkStatus.value) return
  for (const item of orderItems.value) {
    if (item.status !== bulkStatus.value) {
      await updateItemStatus({ ...item, newStatus: bulkStatus.value, comment: bulkComment.value })
    }
  }
  bulkStatus.value = ''
  bulkComment.value = ''
}

const toggleItemHistory = (id) => {
  expandedHistory.value = expandedHistory.value === id ? null : id
}

const closeModal = () => {
  selectedOrder.value = null
  orderItems.value = []
  bulkStatus.value = ''
  bulkComment.value = ''
  expandedHistory.value = null
}

onMounted(() => {
  loadOrders()
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/admin-orders';
</style>