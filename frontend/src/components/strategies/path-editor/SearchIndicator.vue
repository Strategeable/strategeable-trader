<template>
  <div class="search-indicator">
    <input
      placeholder="Search for an indicator . . ."
      type="text"
      v-model="input"
    >
    <div class="indicators">
      <div
        class="indicator"
        v-for="indicator in filteredIndicators"
        :key="indicator.key"
        @click="$emit('select', indicator.key)"
      >{{ indicator.name }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import indicators from '@/assets/data/indicators'

export default defineComponent({
  emits: ['select'],
  setup () {
    const input = ref<string>('')
    const filteredIndicators = computed(() => indicators.filter(i => (i.shortName + i.name).toLowerCase().includes(input.value.toLowerCase())))

    return {
      indicators,
      input,
      filteredIndicators
    }
  }
})
</script>

<style lang="scss" scoped>
.search-indicator {
  input {
    margin-bottom: 1rem;
    width: 100%;
    border: none;
    padding: 1rem 0;
    border-bottom: 2px solid var(--border-color);
    outline: none;
  }
  .indicator {
    color: var(--text-secondary);
    padding: 0.5rem 1rem;
    border: 1px solid var(--border-color);
    border-top: none;
    cursor: pointer;
    &:first-child {
      border-top: 1px solid var(--border-color);
    }
    &:hover {
    color: var(--text);
    }
  }
}
</style>
