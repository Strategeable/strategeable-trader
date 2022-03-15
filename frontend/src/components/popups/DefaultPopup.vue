<template>
  <div @click.self="$emit('close')" class="popup">
    <div class="inner">
      <slot/>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.popup {
  position: fixed;
  display: flex;
  justify-content: center;
  align-items: center;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: rgba(0, 0, 0, 0.3);
  z-index: 10;
  .inner {
    padding: 1.5rem;
    background-color: var(--background);
    border: 1px solid var(--border-color);
    width: 500px;
  }
}
</style>
