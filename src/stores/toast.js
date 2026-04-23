import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useToastStore = defineStore('toast', () => {
  const messages = ref([]);

  function show(text, type = 'info', duration = 3000) {
    const id = Date.now() + Math.random();
    messages.value.push({ id, text, type });
    setTimeout(() => {
      messages.value = messages.value.filter(m => m.id !== id);
    }, duration);
  }

  function success(text, duration = 3000) { show(text, 'success', duration); }
  function error(text, duration = 3000) { show(text, 'error', duration); }
  function warning(text, duration = 3000) { show(text, 'warning', duration); }
  function info(text, duration = 3000) { show(text, 'info', duration); }

  return { messages, show, success, error, warning, info };
});