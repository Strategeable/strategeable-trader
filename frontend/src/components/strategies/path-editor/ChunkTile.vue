<template>
  <div class="tile">
    <p>{{ name }}</p>
    <div class="controls">
      <fa-icon icon="times"
        @click="$emit('delete')"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    name: {
      required: true,
      type: String
    }
  }
})
</script>

<style lang="scss" scoped>
.tile {
  p {
    display: flex;
    align-items: center;
    &::before {
      content: 'CHUNK';
      font-size: 12px;
      margin-right: 0.75rem;
    }
  }
  padding: 0.7rem;
  background-color: var(--primary);
  border: 1px solid var(--primary-darken);
  color: var(--text-inverse);
  display: flex;
  justify-content: space-between;
  align-items: center;
  svg {
    cursor: pointer;
    color: var(--text-inverse);
  }
  .controls {
    display: flex;
    align-items: center;
    svg {
      margin-left: 1rem;
    }
  }
}
</style>
