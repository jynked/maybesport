<template>
  <Transition name="modal">
    <div v-if="isOpen" class="quantity-modal" @click="close">
      <div class="modal-content" @click.stop>
        <div class="quantity-control">
          <button @click="decrement" :disabled="quantity <= 1">-</button>
          <input type="text" v-model="inputValue" @input="validateInput" @blur="fixLeadingZero" />
          <button @click="increment">+</button>
        </div>
        <div class="actions">
          <button @click="confirm">{{ $t('add') }}</button>
          <button @click="close">{{ $t('cancel') }}</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
  isOpen: Boolean,
  initialQuantity: { type: Number, default: 1 },
});
const emit = defineEmits(['close', 'confirm']);

const quantity = ref(props.initialQuantity);
const inputValue = ref(String(props.initialQuantity));

watch(() => props.isOpen, (val) => {
  if (val) {
    quantity.value = props.initialQuantity;
    inputValue.value = String(props.initialQuantity);
  }
});

function validateInput(e) {
  let val = e.target.value;
  val = val.replace(/[^0-9]/g, '');
  if (val === '' || val === '0') val = '1';
  if (val.length > 1 && val[0] === '0') val = val.replace(/^0+/, '');
  const num = parseInt(val, 10);
  if (!isNaN(num) && num > 0) {
    quantity.value = num;
    inputValue.value = String(num);
  } else {
    inputValue.value = String(quantity.value);
  }
}

function fixLeadingZero() {
  let val = inputValue.value;
  if (val.length > 1 && val[0] === '0') {
    val = val.replace(/^0+/, '');
    if (val === '') val = '1';
    quantity.value = parseInt(val, 10);
    inputValue.value = String(quantity.value);
  }
}

function increment() {
  quantity.value++;
  inputValue.value = String(quantity.value);
}

function decrement() {
  if (quantity.value > 1) {
    quantity.value--;
    inputValue.value = String(quantity.value);
  }
}

function confirm() {
  emit('confirm', quantity.value);
  close();
}

function close() {
  emit('close');
}
</script>

<style scoped lang="scss">
@use '../assets/styles/components/quantity-modal';
</style>