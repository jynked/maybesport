const appearDirective = {
  mounted(el, binding) {
    const defaultOptions = {
      threshold: 0.1,
      rootMargin: '0px 0px -20px 0px',
      transition: 'all 0.6s ease',
      delay: 0,
      forceAnimate: false
    };

    const options = { ...defaultOptions, ...binding.value };

    if (options.forceAnimate) {
      applyAnimation(el, options);
      return;
    }

    el.style.opacity = '0';
    el.style.transform = 'translateY(20px)';
    el.style.transition = options.transition;

    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          setTimeout(() => {
            el.style.opacity = '1';
            el.style.transform = 'translateY(0)';
            
            if (!binding.modifiers.repeat) {
              observer.unobserve(el);
            }
          }, options.delay);
        } else if (binding.modifiers.repeat) {
          setTimeout(() => {
            el.style.opacity = '0';
            el.style.transform = 'translateY(20px)';
          }, options.delay);
        }
      });
    }, {
      threshold: options.threshold,
      rootMargin: options.rootMargin
    });

    observer.observe(el);

    el._appearObserver = observer;
  },
  unmounted(el) {
    if (el._appearObserver) {
      el._appearObserver.disconnect();
    }
  }
};

export default appearDirective;
