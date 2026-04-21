<template>
  <div class="order-details-modal" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>{{ $t('orderDetails') }} #{{ order.id }}</h2>
        <button class="close-button" @click="$emit('close')">×</button>
      </div>
      <div v-if="loading" class="loading-spinner">
        <Loader />
      </div>
      <div v-else-if="order" class="modal-body">
        <div class="timeline">
          <h3>{{ $t('deliveryTimeline') }}</h3>
          <div class="timeline-steps">
            <div
              v-for="(entry, idx) in order.statusHistory"
              :key="idx"
              class="timeline-step"
              :class="{ active: idx === order.statusHistory.length - 1 }"
            >
              <div class="step-marker"></div>
              <div class="step-content">
                <div class="step-status">
                  {{ $t(`orderStatus${capitalize(entry.status)}`) }}
                </div>
                <div class="step-date">{{ formatDate(entry.timestamp) }}</div>
                <div class="step-desc" v-if="entry.description">{{ entry.description }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="order-items">
          <h3>{{ $t('orderedItems') }}</h3>
          <div class="items-list">
            <div
              v-for="item in order.items"
              :key="item.uniqueId + item.size"
              class="order-item"
            >
              <router-link :to="{ name: 'Item', params: { itemId: item.uniqueId } }" class="item-link" @click="$emit('close')">
                <img :src="item.image" :alt="item.title.ru" class="item-image" />
                <div class="item-info">
                  <div class="item-title">
                    {{ $i18n.locale === 'en' ? item.title.en : item.title.ru }}
                  </div>
                  <div class="item-size">{{ $t('size') }}: {{ item.size }}</div>
                  <div class="item-price">{{ item.price.toLocaleString() }} ₽</div>
                  <div class="item-quantity">{{ $t('quantity') }}: {{ item.quantity }}</div>
                </div>
              </router-link>
            </div>
          </div>
        </div>
      </div>
      <div class="total-amount">
        <span>{{ $t('totalAmount') }}:</span>
        {{ order.totalAmount.toLocaleString() }} ₽
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useOrdersStore } from '../stores/orders'

const props = defineProps({
  orderId: Number,
})
const emit = defineEmits(['close'])

const ordersStore = useOrdersStore()
const order = ref(null)
const loading = ref(false)

const loadOrder = async () => {
  loading.value = true
  try {
    const data = await ordersStore.fetchOrderDetails(props.orderId)
    order.value = data
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const formatDate = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleString()
}

const capitalize = (str) => str.charAt(0).toUpperCase() + str.slice(1)

watch(
  () => props.orderId,
  () => {
    if (props.orderId) loadOrder()
  },
  { immediate: true }
)
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/order-details-modal';
</style>