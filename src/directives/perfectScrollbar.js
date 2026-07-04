// к компоненту всегда фиксированная максимальная высота и relative
// overflow убираем и не ставим

import PerfectScrollbar from 'perfect-scrollbar';

export default {
  mounted(el, binding) {
    const defaultOptions = {
      suppressScrollX: true,       // отключаем горизонтальный скролл
      wheelPropagation: false,     // не передаём прокрутку родителям
      wheelSpeed: 0.2,            
      swipeEasing: true,           // плавное замедление при свайпе
      minScrollbarLength: 20,
      maxScrollbarLength: null,
      useBothWheelAxes: false,
      scrollXMarginOffset: 0,
      scrollYMarginOffset: 0,
      stopPropagationOnClick: true,
      // === ВАЖНО ДЛЯ МОБИЛЬНЫХ УСТРОЙСТВ ===
      touchEvents: true,           // включаем обработку touch-событий
      swipePropagation: false,     // НЕ передаём свайп родительским элементам
    };

    const options = binding.value ? { ...defaultOptions, ...binding.value } : defaultOptions;

    const checkAndInit = () => {
      const needsScroll = el.scrollHeight > el.clientHeight;

      if (needsScroll && !el._ps) {
        el._ps = new PerfectScrollbar(el, options);
      } else if (!needsScroll && el._ps) {
        el._ps.destroy();
        el._ps = null;
      } else if (needsScroll && el._ps) {
        el._ps.update();
      }
    };

    checkAndInit();

    const observer = new MutationObserver(() => {
      if (el._ps) {
        el._ps.update();
      } else {
        checkAndInit();
      }
    });
    observer.observe(el, { childList: true, subtree: true });
    el._observer = observer;

    const resizeObserver = new ResizeObserver(() => {
      if (el._ps) el._ps.update();
    });
    resizeObserver.observe(el);
    el._resizeObserver = resizeObserver;
  },

  updated(el) {
    if (el._ps) {
      el._ps.update();
    }
  },

  beforeUnmount(el) {
    if (el._ps) {
      el._ps.destroy();
      delete el._ps;
    }
    if (el._observer) {
      el._observer.disconnect();
      delete el._observer;
    }
    if (el._resizeObserver) {
      el._resizeObserver.disconnect();
      delete el._resizeObserver;
    }
  }
};