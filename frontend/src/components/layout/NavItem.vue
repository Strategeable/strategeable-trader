<template>
  <div class="nav-item">
    <router-link :to="to">
      {{ name }}
    </router-link>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    name: String,
    to: String
  }
}
</script>

<style lang="scss" scoped>
.nav-item {
  a {
    padding: 1.2rem 1rem;
    display: block;
    color: var(--text-inverse);
    text-decoration: none;
  }
  &.active, &:hover {
    background-color: var(--primary-lighten);
  }
}
</style>
