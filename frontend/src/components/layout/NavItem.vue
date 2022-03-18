<template>
  <div class="nav-item" :class="{ active }">
    <router-link :to="to">
      {{ name }}
    </router-link>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    name: String,
    to: String,
    active: {
      type: Boolean,
      default: false
    }
  }
}
</script>

<style lang="scss" scoped>
.nav-item {
  a {
    padding: 1.2rem 1rem;
    display: block;
    color: var(--button-text);
    font-weight: bold;
    text-decoration: none;
  }
  &.active, &:hover {
    background-color: var(--primary-darken);
  }
}
</style>
