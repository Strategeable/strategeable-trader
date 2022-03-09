<template>
  <div class="control-bar">
    <div class="group">
    <button @click="$emit('exportStrategy')">Export strategy</button>
      <button :disabled="!canUndo" @click="$emit('undo')">
        <fa-icon icon="undo"/>
      </button>
      <button :disabled="!canRedo" @click="$emit('redo')">
        <fa-icon icon="redo"/>
      </button>
    </div>
    <button :disabled="!canSave" @click="$emit('save')">
      Save
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  emits: ['save', 'undo', 'redo', 'exportStrategy'],
  props: {
    canRedo: {
      required: true,
      type: Boolean
    },
    canUndo: {
      required: true,
      type: Boolean
    },
    canSave: {
      required: true,
      type: Boolean
    }
  }
})
</script>

<style lang="scss" scoped>
.control-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  padding: 0.8rem;
  background-color: var(--background-darken);
  z-index: 1;
  max-width: var(--container-width);
  margin: 0 auto;
  display: flex;
  justify-content: flex-end;
  .group {
    margin-right: 1rem;
  }
  button {
    margin-left: 0.5rem;
  }
}
</style>
