<template>
    <div class="order-card" :class="statusClass" v-appear="{ delay: delay }">
        <div class="order-header">
            <div class="order-date">{{ formattedDate }}</div>
            <div class="order-status" :class="statusClass">
                {{ $t(statusText) }}
            </div>
        </div>
        <div class="order-body">
            <div class="order-info">
                <div class="info-item">
                    <span class="info-label">{{ $t('itemsCount') }}:</span>
                    <span class="info-value">{{ order.items.length }} {{ $t('pieces') }}</span>
                </div>
                <div class="info-item">
                    <span class="info-label">{{ $t('totalAmount') }}:</span>
                    <span class="info-value">{{ formattedAmount }} ₽</span>
                </div>
            </div>
        </div>
        <div class="order-footer">
            <button class="repeat-order-btn" @click="repeatOrder">
                {{ $t('repeatOrder') }}
            </button>
            <button class="details-btn" @click="$emit('details', order.id)">
                {{ $t('details') }}
            </button>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
    order: Object,
    delay: { type: Number, default: 0 }
})
const emit = defineEmits(['details', 'repeat'])

const formattedDate = computed(() => {
    const date = new Date(props.order.createdAt)
    return date.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
})

const formattedAmount = computed(() => {
    return props.order.totalAmount.toLocaleString('ru-RU')
})

const statusClass = computed(() => {
    const status = props.order.status
    switch (status) {
        case 'created':
            return 'status-created'
        case 'processing':
        case 'shipped':
            return 'status-shipped'
        case 'delivered':
            return 'status-delivered'
        case 'received':
            return 'status-received'
        case 'cancelled':
            return 'status-cancelled'
        default:
            return ''
    }
})

const statusText = computed(() => {
    const map = {
        created: 'orderStatusCreated',
        processing: 'orderStatusProcessing',
        shipped: 'orderStatusShipped',
        delivered: 'orderStatusDelivered',
        received: 'orderStatusReceived',
        cancelled: 'orderStatusCancelled',
    }
    return map[props.order.status] || props.order.status
})

const repeatOrder = () => {
    emit('repeat', props.order.id)
}
</script>

<style lang="scss" scoped>
@use '../assets/styles/components/order-card';
</style>