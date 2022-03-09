<template>
  <div class="strategy">
    <div class="general-settings section">
      <div class="input">
        <p>Name</p>
        <input
          type="text"
          v-model="name"
        >
      </div>
      <div class="input">
        <p>Symbols</p>
        <v-select
          multiple
          :options="['BTC/USDT', 'ETH/USDT', 'LTC/USDT', 'UNI/USDT', 'DOT/USDT', 'AVAX/USDT']"
          v-model="symbols"
        />
      </div>
    </div>
    <div class="section">
      <h2>Chunks</h2>
      <div class="chunks">
        <div
          v-for="chunk in chunks"
          :key="chunk.id"
          class="chunk"
        >
          <p>{{ chunk.name }}</p>
          <fa-icon
            icon="times"
            @click="deleteChunk(chunk.id)"
          />
          <fa-icon
            icon="pencil-alt"
            @click="editChunk(chunk.id)"
          />
        </div>
      </div>
      <button @click="newChunk">+ New chunk</button>
    </div>
    <div class="paths section">
      <div class="buy-sell-paths"
        v-for="type in ['BUY', 'SELL']"
        :key="type"
        :class="type.toLowerCase()"
      >
        <h2>{{ type }} paths</h2>
        <button @click="newPath(getType(type))">+ New {{ type.toLowerCase() }} path</button>
        <div class="edit-space">
          <div class="list">
            <div class="path"
              v-for="(path, index) in paths.filter(p => p.type === type)"
              :key="path.id + index"
              :class="{ active: path.id === openEditor[getType(type)] }"
              @click.self="editPath(getType(type), path.id)"
            >
              <p @click="editPath(getType(type), path.id)">
                {{ path.name && path.name !== '<signal path>' ? path.name : `#${index + 1}` }}
              </p>
              <fa-icon
                icon="times"
                @click="deletePath(path.id)"
              />
            </div>
          </div>
          <path-editor
            v-if="openEditor[getType(type)] && paths.some(p => p.id === openEditor[getType(type)])"
            :chunks="chunks"
            :path="paths.find(p => p.id === openEditor[getType(type)])"
            :key="openEditor[getType(type)]"
          />
        </div>
      </div>
    </div>
    <div
      class="editor-overlay"
      v-if="!!editingChunk"
      @click.self="closeChunkEditor"
    >
      <path-editor :chunk="editingChunk" :chunks="[]"/>
    </div>
    <control-bar
      :canSave="true"
      :canUndo="true"
      :canRedo="false"
      @save="save"
      @undo="undo"
      @redo="redo"
      @exportStrategy="exportStrategy"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref, watch } from '@vue/runtime-core'
import { v4 } from 'uuid'
import exportFromJSON from 'export-from-json'

import { Chunk, Path } from '@/types/Path'
import { Strategy } from '@/types/Strategy'

import PathEditor from '@/components/strategies/path-editor/PathEditor.vue'
import ControlBar from '@/components/strategies/path-editor/ControlBar.vue'

type EditorType = 'BUY' | 'SELL'

export default defineComponent({
  components: {
    PathEditor, ControlBar
  },
  setup () {
    const paths = ref<Path[]>([])
    const openEditor = ref<{ SELL: string | undefined, BUY: string | undefined }>({ SELL: undefined, BUY: undefined })
    const chunks = ref<Chunk[]>([])
    const name = ref<string>('')
    const symbols = ref<string[]>([])
    const editingChunk = ref<Chunk>()
    const editHistory: string[] = []
    const strategy = computed(() => {
      const strat: Strategy = {
        name: name.value,
        symbols: [],
        chunks: chunks.value,
        paths: paths.value
      }
      return strat
    })

    let timeout: number
    watch(strategy, strat => {
      if (timeout) {
        clearTimeout(timeout)
      }
      timeout = setTimeout(() => {
        if (editHistory.includes(JSON.stringify(strat))) return
        editHistory.push(JSON.stringify(strat))
      }, 1500)
    }, {
      deep: true
    })

    onMounted(() => {
      newPath('BUY')
      newPath('SELL')
    })

    function newPath (type: EditorType): void {
      const path: Path = { id: v4(), name: undefined, whitelist: [], steps: [], type }
      paths.value.push(path)
      openEditor.value[type] = path.id
    }

    function deletePath (id: string): void {
      const path = paths.value.find(p => p.id === id)
      if (!path) return
      if (paths.value.filter(p => p.type === path.type).length === 1) return

      for (const [key, val] of Object.entries(openEditor.value)) {
        if (val === id) {
          openEditor.value[key as EditorType] = undefined
        }
      }

      paths.value = paths.value.filter(p => p.id !== id)
    }

    function editPath (type: EditorType, id: string) {
      openEditor.value[type] = id
    }

    function newChunk (): void {
      chunks.value.push({ id: v4(), name: 'New chunk', steps: [] })
    }

    function editChunk (id: string) {
      editingChunk.value = chunks.value.find(c => c.id === id)
    }

    function deleteChunk (id: string): void {
      let isUsed = false
      for (const path of paths.value) {
        for (const step of path.steps) {
          if (step.type === 'CHUNK_ID') {
            if (step.chunkId === id) isUsed = true
          }
        }
      }

      if (isUsed) {
        alert('Cannot delete this chunk, it is being used in a path.')
      } else {
        chunks.value = chunks.value.filter(c => c.id !== id)
      }
    }

    function closeChunkEditor () {
      editingChunk.value = undefined
    }

    function getType (type: string): EditorType {
      if (type === 'BUY') return 'BUY'
      return 'SELL'
    }

    function undo () {
      // todo
    }

    function redo () {
      // todo
    }

    function save () {
      // todo
    }

    function exportStrategy () {
      const fileName = `${name.value} - Strategy`
      const exportType = exportFromJSON.types.json

      exportFromJSON({ data: strategy.value, fileName, exportType })
    }

    return {
      paths,
      openEditor,
      chunks,
      editingChunk,
      name,
      symbols,
      newPath,
      newChunk,
      deletePath,
      deleteChunk,
      editChunk,
      closeChunkEditor,
      editPath,
      getType,
      undo,
      redo,
      save,
      exportStrategy
    }
  }
})
</script>

<style lang="scss" scoped>
.general-settings {
  display: flex;
  > div {
    margin-right: 1rem;
    display: flex;
    flex-direction: column;
  }
}

.input {
  display: flex;
  p {
    margin-bottom: 0.5rem;
    font-size: 14px;
    font-weight: 600;
  }
  input {
    padding: 0.5rem;
    outline: none;
    border: 1px solid var(--border-color);
    color: var(--text-color);
    width: 300px;
  }
  .v-select {
    width: 300px;
  }
}

.chunks {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  .chunk {
    margin-right: 0.5rem;
    margin-bottom: 0.5rem;
    padding: 0.8rem;
    background-color: var(--primary);
    color: var(--text-inverse);
    display: flex;
    svg {
      cursor: pointer;
      margin-left: 1rem;
    }
  }
}

.editor-overlay {
  position: fixed;
  z-index: 10;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
}

.paths {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.edit-space {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: 1rem;
  margin-top: 1rem;
  .list {
    .path {
      padding: 0.5rem;
      background-color: var(--primary);
      color: var(--text-inverse);
      margin-bottom: 0.5rem;
      cursor: pointer;
      display: flex;
      justify-content: space-between;
      align-items: center;
      user-select: none;
      svg {
        cursor: pointer;
      }
      &.active, &:hover {
        background-color: var(--primary-lighten);
      }
    }
  }
}
</style>
