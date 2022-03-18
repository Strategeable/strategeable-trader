<template>
  <div @click.self="$emit('close')" class="select-chunk">
    <div class="inner">
      <v-select
        :options="chunks"
        label="name"
        :reduce="(x: any) => x.id"
        v-model="chunkId"
      />
      <button
        @click="$emit('select', chunkId)"
        :disabled="!chunkId"
      >Select</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from 'vue'

import { Chunk } from '@/types/Path'

export default defineComponent({
  props: {
    chunks: {
      required: true,
      type: Array as PropType<Chunk[]>
    }
  },
  setup () {
    const chunkId = ref<string>()

    return {
      chunkId
    }
  }
})
</script>

<style lang="scss" scoped>
.select-chunk {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.2);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10;
  .inner {
    padding: 2rem;
    background-color: var(--background);
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
    width: 100%;
    max-width: 300px;
    margin: 0 1rem;
    button {
      width: 100%;
      margin-top: 1rem;
      &:disabled {
        background-color: var(--background-darken);
        cursor: unset;
        color: var(--text-color);
        border-color: var(--border-color);
      }
    }
  }
}
</style>
