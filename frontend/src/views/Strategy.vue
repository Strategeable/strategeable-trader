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
      <h2>Variables</h2>
      <div class="variables">
        <div
          class="variable"
          v-for="variable in variables"
          :key="variable.id"
        >
          <div class="input">
            <p>Name</p>
            <input type="text" v-model="variable.key">
          </div>
          <div class="input">
            <p>Value</p>
            <input type="number" v-model="variable.value">
          </div>
        </div>
      </div>
      <button
        @click="newVariable"
        class="outline small"
      >New variable</button>
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
      <button
        @click="newChunk"
        class="outline small"
      >+ New chunk</button>
    </div>
    <div class="paths section">
      <div class="buy-sell-paths"
        v-for="type in ['BUY', 'SELL']"
        :key="type"
        :class="type.toLowerCase()"
      >
        <h2>{{ type }} paths</h2>
        <button
          @click="newPath(getType(type))"
          class="outline small"
        >+ New {{ type.toLowerCase() }} path</button>
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
            :variables="variables"
            :key="openEditor[getType(type)]"
          />
        </div>
      </div>
    </div>
    <div class="backtests section">
      <h2>Backtests</h2>
      <div class="run-backtest">
        <div class="backtest">
          <div class="input">
            <p>Start balance</p>
            <input
              type="number"
              v-model="backtestParameters.startBalance"
            >
          </div>
          <div class="input">
            <p>From date</p>
            <input
              type="date"
              :value="moment(backtestParameters.fromDate.toString()).format('YYYY-MM-DD')"
              @change="e => backtestParameters.fromDate = new Date(e.target.value)"
            >
          </div>
          <div class="input">
            <p>To date</p>
            <input
              type="date"
              :value="moment(backtestParameters.toDate.toString()).format('YYYY-MM-DD')"
              @change="e => backtestParameters.toDate = new Date(e.target.value)"
            >
          </div>
        </div>
        <button @click="backtest">Backtest</button>
      </div>
      <div v-if="runningBacktest">
        <p>{{ runningBacktest }}</p>
      </div>
      <backtest-result-comp
        v-for="backtest in backtestResults"
        :key="backtest.id"
        :backtest="backtest"
        @restore="() => restoreStrategy(backtest.strategy)"
      />
    </div>
    <div
      class="editor-overlay"
      v-if="!!editingChunk"
      @click.self="closeChunkEditor"
    >
      <path-editor :chunk="editingChunk" :chunks="[]" :variables="[]"/>
    </div>
    <control-bar
      :canSave="canSave"
      :canUndo="false"
      :canRedo="false"
      @save="save"
      @exportStrategy="exportStrategy"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref, watch } from '@vue/runtime-core'
import { v4 } from 'uuid'
import exportFromJSON from 'export-from-json'
import moment from 'moment'

import { Chunk, Path } from '@/types/Path'
import { Strategy, Variable } from '@/types/Strategy'
import { BacktestRequestParameters, BacktestResult } from '@/types/Backtest'

import PathEditor from '@/components/strategies/path-editor/PathEditor.vue'
import ControlBar from '@/components/strategies/path-editor/ControlBar.vue'
import BacktestResultComp from '@/components/strategies/BacktestResult.vue'
import { useStore } from 'vuex'
import { useRoute, useRouter } from 'vue-router'

type EditorType = 'BUY' | 'SELL'

export default defineComponent({
  components: {
    PathEditor, ControlBar, BacktestResultComp
  },
  setup () {
    const store = useStore()
    const route = useRoute()
    const router = useRouter()

    const paths = ref<Path[]>([])
    const openEditor = ref<{ SELL: string | undefined, BUY: string | undefined }>({ SELL: undefined, BUY: undefined })
    const chunks = ref<Chunk[]>([])
    const name = ref<string>('')
    const symbols = ref<string[]>([])
    const variables = ref<Variable[]>([])
    const editingChunk = ref<Chunk>()
    const canSave = ref<boolean>(false)

    const strategyId = ref<string | undefined>()
    const strategyCreatedAt = ref<Date>(new Date())
    const strategyLastEdited = ref<Date>(new Date())

    const backtestParameters = ref<BacktestRequestParameters>({
      strategyId: '',
      fromDate: new Date('2022-01-01'),
      toDate: new Date('2022-03-10'),
      startBalance: 1000
    })
    const backtestResults = computed<BacktestResult[]>(() => {
      const id = route.path.split('/')[route.path.split('/').length - 1]
      if (id === 'new') return []
      return (store.getters.backtests[id] || []).sort((a: BacktestResult, b: BacktestResult) => new Date(b.startedOn).getTime() - new Date(a.startedOn).getTime())
    })
    const runningBacktest = ref<string>()

    const strategy = computed(() => {
      const strat: Strategy = {
        id: strategyId.value,
        version: '0.0.1',
        createdAt: new Date(strategyCreatedAt.value.toString()),
        lastEdited: new Date(strategyLastEdited.value.toString()),
        name: name.value,
        symbols: symbols.value,
        chunks: chunks.value,
        paths: paths.value,
        variables: variables.value
      }
      return strat
    })

    watch(strategy, () => {
      canSave.value = true
    }, { deep: true })

    const emptyStrategy = computed(() => {
      if (paths.value.some(p => p.steps.length > 0)) return false
      if (chunks.value.length > 0) return false
      if (symbols.value.length > 0) return false
      if (name.value.length > 0) return false
      return true
    })

    async function loadBacktests () {
      const id = route.path.split('/')[route.path.split('/').length - 1]
      if (id === 'new') return
      store.dispatch('loadBacktests', id)
    }

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
        variables.value = strat.variables

        openFirstPaths()
      }
    }

    onMounted(() => {
      loadBacktests()
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

    async function save () {
      const stratId = await store.dispatch('saveStrategy', strategy.value)
      if (stratId && route.path.endsWith('new')) {
        router.push(`/strategies/${stratId}`)
      }
      canSave.value = false
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
        if (canSave.value) {
          await save()
        }

        const stratId = strategyId.value
        if (!stratId) {
          alert('Save strategy first')
          return
        }
        backtestParameters.value.strategyId = stratId
        runningBacktest.value = 'Backtesting...'

        await store.dispatch('runBacktest', backtestParameters.value)

        runningBacktest.value = undefined
      } catch (err) {
        console.error(err)
        runningBacktest.value = 'Something went wrong with the backtest... Try again.'
      }
    }

    function restoreStrategy (strat: Strategy) {
      // TODO
    }

    function newVariable () {
      variables.value.push({
        type: 'number',
        id: v4(),
        key: '',
        value: undefined
      })
    }

    return {
      paths,
      openEditor,
      chunks,
      editingChunk,
      name,
      symbols,
      backtestResults,
      emptyStrategy,
      backtestParameters,
      canSave,
      runningBacktest,
      variables,
      newPath,
      newChunk,
      deletePath,
      deleteChunk,
      editChunk,
      closeChunkEditor,
      editPath,
      getType,
      save,
      exportStrategy,
      backtest,
      handleUploadStrategy,
      restoreStrategy,
      newVariable,
      moment
    }
  }
})
</script>

<style lang="scss" scoped>
.strategy {
  padding-bottom: 80px;
}

.general-settings {
  display: flex;
  flex-wrap: wrap;
  > div {
    margin-right: 1rem;
    margin-bottom: 1rem;
    display: flex;
    flex-direction: column;
  }
  .input {
    max-width: 100%;
    input {
      max-width: 100% !important;
      @media(max-width: 400px) {
        width: 100%;
      }
    }
    .v-select {
      max-width: 100%;
      width: unset;
    }
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
  @media(max-width: 1080px) {
    grid-template-columns: 1fr;
  }
}

.edit-space {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: 1rem;
  margin-top: 1rem;
  @media(max-width: 600px) {
    grid-template-columns: 1fr;
  }
  .list {
    .path {
      padding: 0.5rem;
      background-color: var(--background-lighten);
      color: var(--text-inverse);
      margin-bottom: 0.5rem;
      cursor: pointer;
      display: flex;
      justify-content: space-between;
      align-items: center;
      user-select: none;
      svg {
        cursor: pointer;
        color: var(--text);
      }
      &.active, &:hover {
        background-color: var(--primary-darken);
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

.run-backtest {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  justify-content: space-between;
  background-color: var(--background-darken);
  padding: 1rem;
  border: 1px solid var(--border-color);
  margin-bottom: 2rem;
  padding-bottom: 0;
  button {
    margin-bottom: 1rem;
  }
  .backtest {
    display: flex;
    flex-wrap: wrap;
    .input {
      display: flex;
      flex-direction: column;
      margin-right: 1rem;
      margin-bottom: 1rem;
      input {
        max-width: 100% !important;
        @media(max-width: 400px) {
          width: 100%;
        }
      }
    }
  }
}

.variables {
  display: flex;
  flex-wrap: wrap;
}

.variable {
  display: flex;
  margin-bottom: 1rem;
  margin-right: 2rem;
  .input {
    margin-right: 0.5rem;
    flex-direction: column;
    width: 125px;
    input {
      width: 100%;
    }
  }
}
</style>
