<template>
  <div @click="change">
    <fa-icon size="lg" :icon="theme === 'light' ? 'sun' : 'moon'"/>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { useStore } from 'vuex'

export default defineComponent({
  setup () {
    const store = useStore()
    const theme = computed(() => store.getters.theme)

    function change () {
      store.dispatch('changeColorTheme')
    }

    return {
      theme,
      change
    }
  }
})
</script>

<style lang="scss" scoped>
div {
  padding: 0 1rem;
  height: 100%;
  display: flex;
  align-items: center;
  cursor: pointer;
  svg {
  color: var(--button-text);
  }
}
</style>
