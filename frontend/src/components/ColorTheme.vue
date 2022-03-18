<template>
  <div @click="change">
    <fa-icon size="lg" :icon="theme === 'light' ? 'sun' : 'moon'"/>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { useStore } from '@/store'

import { ActionTypes } from '@/types/store/action-types'

export default defineComponent({
  setup () {
    const store = useStore()
    const theme = computed(() => store.getters.theme)

    function change () {
      store.dispatch(ActionTypes.CHANGE_COLOR_THEME, undefined)
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
