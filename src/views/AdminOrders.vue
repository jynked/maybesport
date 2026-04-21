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
      <div v-if="orders.length === 0" class="empty">Нет заказов</div>
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
              <td>
                <span :class="'status-' + order.status">{{ getStatusText(order.status) }}</span>
              </td>
              <td>{{ order.delivery_address || '—' }}</td>
              <td>{{ formatDate(order.created_at) }}</td>
              <td>
                <button @click="openChangeStatus(order)" class="btn-change">Изменить статус</button>
              </td>
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
      <div v-if="selectedOrder" class="status-modal" @click.self="closeStatusModal">
        <div class="modal-content">
          <h3>Изменить статус заказа #{{ selectedOrder.id }}</h3>
          <div class="form-group">
            <label>Новый статус</label>
            <select v-model="newStatus">
              <option value="created">Создан</option>
              <option value="processing">Обработка</option>
              <option value="shipped">Отправлен</option>
              <option value="delivered">Доставлен</option>
              <option value="received">Получен</option>
              <option value="cancelled">Отменён</option>
            </select>
          </div>
          <div v-if="showHistory && statusHistory.length > 0" class="status-history">
            <h4>История статусов:</h4>
            <ul>
              <li v-for="(entry, idx) in statusHistory" :key="idx">
                {{ getStatusText(entry.status) }} — {{ formatDate(entry.timestamp) }}
                <span v-if="entry.description">({{ entry.description }})</span>
              </li>
            </ul>
            <button @click="deleteLastStatus" :disabled="updating" class="btn-danger">
              Удалить последний статус
            </button>
          </div>
          <div class="form-group">
            <label>Комментарий (необязательно)</label>
            <textarea v-model="statusDescription" rows="2"></textarea>
          </div>
          <div class="actions">
            <button @click="updateStatus" :disabled="updating">Сохранить</button>
            <button @click="closeStatusModal">Отмена</button>
          </div>
        </div>
      </div>
    </Transition>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api, http } from '../api'


const emit = defineEmits(['page-loaded'])
const orders = ref([])
const loading = ref(false)
const page = ref(1)
const totalPages = ref(1)
const filters = ref({ userEmail: '', status: '' })
const selectedOrder = ref(null)
const newStatus = ref('')
const statusDescription = ref('')
const updating = ref(false)
let searchTimeout = null

const statusHistory = ref([])
const showHistory = ref(false)

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
    emit('page-loaded', false)
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

const fetchOrderHistory = async (orderId) => {
  try {
    const resp = await api.getAdminOrderDetails(orderId);
    statusHistory.value = resp.data.statusHistory || [];
  } catch (err) {
    console.error(err);
  }
};

const openChangeStatus = async (order) => {
  selectedOrder.value = order
  newStatus.value = order.status
  statusDescription.value = ''
  await fetchOrderHistory(order.id)
  showHistory.value = true
}

const deleteLastStatus = async () => {
  if (!selectedOrder.value) return
  if (!confirm('Удалить последнее изменение статуса?')) return
  updating.value = true
  try {
    await api.deleteLastOrderStatus(selectedOrder.value.id)
    await loadOrders()
    await fetchOrderHistory(selectedOrder.value.id)
    selectedOrder.value.status = newStatus.value
    await loadOrders()
  } catch (err) {
    console.error(err)
    alert('Ошибка удаления статуса')
  } finally {
    updating.value = false
  }
}

const closeStatusModal = () => {
  selectedOrder.value = null
  newStatus.value = ''
  statusDescription.value = ''
}

const updateStatus = async () => {
  if (!selectedOrder.value) return
  updating.value = true
  try {
    await api.updateOrderStatus(selectedOrder.value.id, newStatus.value, statusDescription.value)
    await loadOrders()
    await fetchOrderHistory(selectedOrder.value.id)
  } catch (err) {
    console.error(err)
    alert('Ошибка обновления статуса')
  } finally {
    updating.value = false
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/admin-orders';
</style>