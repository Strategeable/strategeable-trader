<template>
  <navbar/>
  <div id="content">
    <router-view/>
  </div>
</template>

<script lang="ts">
import Navbar from '@/components/layout/Nav.vue'

export default {
  components: { Navbar }
}
</script>

<style lang="scss" scoped>
#content {
  width: 100%;
  max-width: var(--container-width);
  margin: 0 auto;
  padding: 0 1rem;
}
</style>
