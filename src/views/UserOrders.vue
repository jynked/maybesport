<template>
    <main>
        <div class="orders-heading">
            <h1 v-appear="{ delay: 200 }">{{ $t('userOrders') }}</h1>
        </div>

        <div v-if="ordersStore.orders?.length === 0 && !ordersStore.loading" class="undefined-items-container"
            v-appear="{ delay: 400 }">
            <p>{{ $t('undefinedItems') }}</p>
            <img src="../assets/img/fail.png" :alt="$t('failAlt')" />
        </div>

        <div v-else class="orders-grid">
            <div v-for="(row, rowIndex) in chunkedOrders" :key="rowIndex" class="orders-row">
                <OrderCard v-for="(order, colIndex) in row" :key="order.id" v-memo="[order.id, order.status, order.totalAmount]" 
                    :order="order" :delay="100 + colIndex * 100"
                    @details="openOrderDetails(order.id)" v-appear="{ delay: 100 + colIndex * 100 }" />
            </div>
        </div>

        <Transition name="modal">
            <OrderDetailsModal v-if="selectedOrderId" :orderId="selectedOrderId" @close="closeOrderDetails" />
        </Transition>
    </main>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useOrdersStore } from '../stores/orders'
import OrderCard from '../components/OrderCard.vue'
import OrderDetailsModal from '../components/OrderDetailsModal.vue'

const ordersStore = useOrdersStore()
const emit = defineEmits(['page-loaded'])

const selectedOrderId = ref(null)

const chunkedOrders = computed(() => {
    const ordersList = ordersStore.orders || []
    const itemsPerRow = 3
    const result = []
    for (let i = 0; i < ordersList?.length; i += itemsPerRow) {
        result.push(ordersList.slice(i, i + itemsPerRow))
    }
    return result
})

const openOrderDetails = (orderId) => {
    selectedOrderId.value = orderId
    document.body.style.overflow = 'hidden'
}

const closeOrderDetails = () => {
    selectedOrderId.value = null
    document.body.style.overflow = 'auto'
}

onMounted(async () => {
    await ordersStore.fetchOrders()
    emit('page-loaded', true)
})

onBeforeUnmount(() => {
    document.body.style.overflow = 'auto'
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/pages/user-orders';
</style>