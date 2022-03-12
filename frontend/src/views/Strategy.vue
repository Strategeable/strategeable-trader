<template>
  <div class="strategy">
    <div v-if="emptyStrategy" class="upload">
      <p>Import strategy <span>.json</span> file</p>
      <input type="file" @change="handleUploadStrategy">
    </div>
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
    <div class="backtests section">
      <h2>Backtests</h2>
      <button @click="backtest">Backtest</button>
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
import { BacktestRequestParameters } from '@/types/Backtest'
import axios from '@/helpers/axios'

import PathEditor from '@/components/strategies/path-editor/PathEditor.vue'
import ControlBar from '@/components/strategies/path-editor/ControlBar.vue'
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'

type EditorType = 'BUY' | 'SELL'

export default defineComponent({
  components: {
    PathEditor, ControlBar
  },
  setup () {
    const store = useStore()
    const route = useRoute()

    const paths = ref<Path[]>([])
    const openEditor = ref<{ SELL: string | undefined, BUY: string | undefined }>({ SELL: undefined, BUY: undefined })
    const chunks = ref<Chunk[]>([])
    const name = ref<string>('')
    const symbols = ref<string[]>([])
    const editingChunk = ref<Chunk>()
    const strategyId = ref<string | undefined>()
    const strategyCreatedAt = ref<Date>(new Date())
    const strategyLastEdited = ref<Date>(new Date())

    const editHistory: string[] = []

    const strategy = computed(() => {
      const strat: Strategy = {
        id: strategyId.value,
        version: '0.0.1',
        createdAt: new Date(),
        lastEdited: new Date(),
        name: name.value,
        symbols: symbols.value,
        chunks: chunks.value,
        paths: paths.value
      }
      return strat
    })
    const emptyStrategy = computed(() => {
      if (paths.value.some(p => p.steps.length > 0)) return false
      if (chunks.value.length > 0) return false
      if (symbols.value.length > 0) return false
      if (name.value.length > 0) return false
      return true
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

    async function loadStrategy (id: string) {
      const strat: Strategy | undefined = await store.dispatch('loadStrategy', id)
      if (strat) {
        paths.value = strat.paths
        chunks.value = strat.chunks
        symbols.value = strat.symbols
        name.value = strat.name
        strategyId.value = strat.id
        strategyCreatedAt.value = strat.createdAt
        strategyLastEdited.value = strat.lastEdited

        openFirstPaths()
      }
    }

    onMounted(() => {
      if (route.params.id === 'new') {
        newPath('BUY')
        newPath('SELL')
      } else {
        loadStrategy(String(route.params.id))
      }
    })

    function openFirstPaths () {
      if (paths.value.find(p => p.type === 'BUY')) {
        openEditor.value.BUY = paths.value.find(p => p.type === 'BUY')?.id
      }

      if (paths.value.find(p => p.type === 'SELL')) {
        openEditor.value.SELL = paths.value.find(p => p.type === 'SELL')?.id
      }
    }

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
            if (step.data === id) isUsed = true
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

    async function save () {
      const success = await store.dispatch('saveStrategy', strategy.value)
      if (success) {
        alert('Saved!')
      }
    }

    function handleUploadStrategy (e: any) {
      const reader = new FileReader()
      reader.onload = (data: any) => {
        const strat = JSON.parse(data.target.result)
        if (typeof strat === 'object' && strat.name && strat.symbols && strat.chunks && strat.paths && strat.version) {
          paths.value = strat.paths
          chunks.value = strat.chunks
          symbols.value = strat.symbols
          name.value = strat.name

          openFirstPaths()
        }
      }
      reader.readAsText(e.target.files[0])
    }

    function exportStrategy () {
      const fileName = `${name.value} - Strategy`
      const exportType = exportFromJSON.types.json

      exportFromJSON({ data: strategy.value, fileName, exportType })
    }

    async function backtest () {
      try {
        const stratId = strategyId.value
        if (!stratId) {
          alert('Save strategy first')
          return
        }
        const data: BacktestRequestParameters = {
          strategyId: stratId,
          fromDate: new Date('2017-08-05'),
          toDate: new Date('2022-03-10'),
          startBalance: 1000
        }
        const result = await axios.post('/backtest', data)
        console.log(result)
      } catch (err) {
        console.error(err)
      }
    }

    return {
      paths,
      openEditor,
      chunks,
      editingChunk,
      name,
      symbols,
      emptyStrategy,
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
      exportStrategy,
      backtest,
      handleUploadStrategy
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

.upload {
  padding: 1rem;
  background-color: var(--background-darken);
  margin-bottom: 1rem;
  p {
    margin-bottom: 1rem;
    font-weight: bold;
  }
}
</style>
