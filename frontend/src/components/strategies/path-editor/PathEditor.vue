<template>
  <div class="path-editor">
    <div class="header">
      <input
        v-if="finalPath"
        placeholder="Add custom name to path"
        type="text"
        v-model="finalPath.name"
      >
    </div>
    <div class="content">
      <div v-if="finalPath && finalPath.steps.length > 0" class="steps">
        <div class="step"
          v-for="(step, i) in finalPath.steps"
          :key="step.id"
        >
          <signal-tile-comp
            v-if="step.type === 'SIGNAL_TILE'"
            :tile="step.data"
            @delete="deleteStep(step.id)"
          />
          <any-signal-tile-comp
            v-if="step.type === 'ANY_SIGNAL_TILE'"
            :tiles="step.data"
            @delete="deleteStep(step.id)"
          />
          <chunk-tile
            v-if="step.type === 'CHUNK_ID'"
            :name="getChunkName(typeof step.data === 'string' ? step.data || '' : '')"
            @delete="deleteStep(step.id)"
          />
          <div class="next" v-if="finalPath && i != finalPath.steps.length - 1">
            <fa-icon icon="angle-down"/>
          </div>
        </div>
      </div>
      <div v-else>
        Add some tiles to this path
      </div>
      <div class="add">
        <button
          @click="addStep(getStepType('SIGNAL_TILE'))"
        >
          Add signal
        </button>
        <button
          @click="addStep(getStepType('ANY_SIGNAL_TILE'))"
        >
          Add ANY group
        </button>
        <button
          @click="addStep(getStepType('CHUNK_ID'))"
          v-if="chunks.length > 0"
        >
          Add chunk
        </button>
      </div>
    </div>
    <select-chunk
      v-if="selectNewChunk"
      @close="selectNewChunk = false"
      @select="id => addChunk(id)"
      :chunks="chunks"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from 'vue'
import { v4 } from 'uuid'

import { AnySignal, Chunk, Operand, Path, SignalTile, StepType } from '@/types/Path'
import SignalTileComp from '@/components/strategies/path-editor/SignalTile.vue'
import AnySignalTileComp from '@/components/strategies/path-editor/AnySignalTile.vue'
import ChunkTile from '@/components/strategies/path-editor/ChunkTile.vue'
import SelectChunk from '@/components/strategies/path-editor/SelectChunk.vue'

export default defineComponent({
  components: {
    SignalTileComp, AnySignalTileComp, ChunkTile, SelectChunk
  },
  props: {
    path: {
      type: Object as PropType<Path>
    },
    chunk: {
      type: Object as PropType<Chunk>
    },
    chunks: {
      required: true,
      type: Array as PropType<Chunk[]>
    }
  },
  setup (props) {
    const drag = ref<boolean>(false)
    const selectNewChunk = ref<boolean>(false)

    function addStep (type: StepType) {
      const p = props.path || props.chunk
      if (!p) return

      const signalTile: SignalTile = { id: v4(), name: '', operand: Operand.GREATER_THAN, persistence: 1 }
      const anySignal: AnySignal = { signals: [], amount: 1 }

      switch (type) {
        case StepType.SIGNAL_TILE:
          p.steps.push({ id: v4(), type: StepType.SIGNAL_TILE, data: signalTile })
          break
        case StepType.ANY_SIGNAL_TILE:
          p.steps.push({ id: v4(), type: StepType.ANY_SIGNAL_TILE, data: anySignal })
          break
        case StepType.CHUNK_ID:
          selectNewChunk.value = true
      }
    }

    function addChunk (id: string) {
      const p = props.path || props.chunk
      if (!p) return

      p.steps.push({ id: v4(), type: StepType.CHUNK_ID, data: id })
      selectNewChunk.value = false
    }

    function deleteStep (id: string) {
      const p = props.path || props.chunk
      if (!p) return
      p.steps = p.steps.filter(s => s.id !== id)
    }

    function getName (): string | undefined {
      if (props.path) return props.path.name
      if (props.chunk) return props.chunk.name
      return 'undefined'
    }

    function getWhitelist (): string[] | null {
      if (props.path) return props.path.whitelist
      return null
    }

    function getStepType (type: string): StepType {
      if (type === StepType.SIGNAL_TILE) return StepType.SIGNAL_TILE
      if (type === StepType.ANY_SIGNAL_TILE) return StepType.ANY_SIGNAL_TILE
      return StepType.CHUNK_ID
    }

    function getChunkName (id: string): string {
      const found = props.chunks.find(c => c.id === id)
      if (found) {
        return found.name
      } else {
        return 'unknown'
      }
    }

    return {
      name: getName(),
      whitelist: getWhitelist(),
      finalPath: props.path || props.chunk,
      drag,
      selectNewChunk,
      addStep,
      deleteStep,
      getStepType,
      getChunkName,
      addChunk
    }
  }
})
</script>

<style lang="scss" scoped>
.path-editor {
  width: 100%;
  max-width: 800px;
  &:hover {
    .header {
      input {
        border: 1px dashed var(--primary-darken);
      }
    }
  }
  .header {
    padding: 1rem;
    background-color: var(--primary);
    color: var(--text-inverse);
    input {
      background-color: var(--primary);
      border: 1px dashed var(--primary-darken);
      color: var(--text-inverse);
      padding: 0.5rem;
      outline: none;
      &::placeholder {
        color: var(--primary-light);
      }
    }
  }
  .content {
    padding: 0.5rem;
    background-color: var(--background);
    border: 1px solid var(--border-color);
    max-height: 80vh;
    overflow-y: auto;
  }
  .next {
    display: flex;
    justify-content: center;
    margin: 0.5rem 0;
  }
  .add {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border-color);
    display: flex;
    flex-wrap: wrap;
    button {
      margin-right: 0.5rem;
      margin-bottom: 0.5rem;
    }
  }
}
</style>
