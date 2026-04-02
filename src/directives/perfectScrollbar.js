// к компоненту всегда фиксированная максимальная высота и relative
// overflow убираем и не ставим

import PerfectScrollbar from 'perfect-scrollbar';

export default {
  mounted(el, binding) {
    const defaultOptions = {
      suppressScrollX: true,       // отключаем горизонтальный скролл
      wheelPropagation: false,     // не передаём прокрутку родителям
      wheelSpeed: .2,               // скорость прокрутки колесом (1 = нормально)
      swipeEasing: true,           // плавное замедление при свайпе (для мобильных)
      minScrollbarLength: 20,      // минимальная длина ползунка
      maxScrollbarLength: null,    // максимальная длина ползунка (null = авто)
      useBothWheelAxes: false,     // использовать оба направления мыши? (false - только вертикаль)
      scrollXMarginOffset: 0,      // отступ для горизонтального скролла (не нужно)
      scrollYMarginOffset: 0,      // отступ для вертикального скролла
      stopPropagationOnClick: true, // остановить всплытие клика на ползунке
    };

    const options = binding.value ? { ...defaultOptions, ...binding.value } : defaultOptions;
    const ps = new PerfectScrollbar(el, options);
    el._ps = ps;

    const observer = new MutationObserver(() => {
      ps.update();
    });
    observer.observe(el, { childList: true, subtree: true });
    el._observer = observer;
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
  }
};