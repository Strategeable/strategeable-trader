<template>
  <div class="any-signal-tile">
    <div class="top">
      <p>ANY</p>
      <input
        type="number"
        v-model="tiles.amount"
      >
    </div>
    <div class="signals">
      <signal-tile-comp
        v-for="tile in tiles.signals"
        :key="tile.id"
        :tile="tile"
        @delete="() => deleteSignal(tile.id)"
      />
      <div class="buttons">
        <button @click="$emit('delete')" class="delete">Delete</button>
        <button @click="addSignal">Add signal</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { AnySignal } from '@/types/Path'
import { defineComponent, PropType } from 'vue'

import SignalTileComp from '@/components/strategies/path-editor/SignalTile.vue'
import { v4 } from 'uuid'

export default defineComponent({
  components: { SignalTileComp },
  props: {
    tiles: {
      required: true,
      type: Object as PropType<AnySignal>
    }
  },
  emits: ['delete'],
  setup (props) {
    function addSignal () {
      props.tiles.signals.push({ id: v4(), name: '', persistence: 1 })
    }

    function deleteSignal (id: string) {
      props.tiles.signals = props.tiles.signals.filter(s => s.id !== id)
    }

    return {
      addSignal,
      deleteSignal
    }
  }
})
</script>

<style lang="scss" scoped>
.any-signal-tile {
  padding: 0.5rem;
  background-color: var(--background-darken);
  .top {
    display: flex;
    align-items: center;
    justify-content: center;
    input {
      background-color: transparent;
      border: none;
      margin-left: 0.2rem;
      outline: none;
      padding: 0.3rem;
      width: 40px;
      font-size: 17px;
    }
  }
  .signals {
    > * {
      margin-top: 0.5rem;
    }
  }
  button {
    width: 100%;
    background-color: var(--background);
    color: var(--text);
    border: 1px dashed var(--primary);
  }
  .buttons {
    display: grid;
    grid-template-columns: 100px 1fr;
    gap: 0.5rem;
    .delete {
      border-color: var(--text-secondary);
      color: var(--text-secondary);
      background-color: var(--background-darken);
      &:hover {
      border-color: var(--text-color);
      color: var(--text-color);
      }
    }
  }
}
</style>
