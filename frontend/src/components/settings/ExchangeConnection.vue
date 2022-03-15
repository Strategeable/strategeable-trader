<template>
  <div class="exchange-connection">
    <div class="name">
      <img :src="icon" alt="logo">
      <p>{{ exchangeConnection.name }}</p>
    </div>
    <p><strong>3</strong> bots (1 running)</p>
    <p>Created on {{ moment(exchangeConnection.createdOn).format('DD MMM YYYY') }}</p>
    <div class="controls">
      <button
        class="rounded small outline"
        v-if="canDelete"
        @click="$emit('delete')"
      >
        <fa-icon icon="trash"/>
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { ExchangeConnection, exchangeDetails } from '@/types/Exchange'
import { defineComponent, PropType, ref } from 'vue'
import moment from 'moment'

export default defineComponent({
  emits: ['delete'],
  props: {
    exchangeConnection: {
      type: Object as PropType<ExchangeConnection>,
      required: true
    }
  },
  setup (props) {
    const icon = exchangeDetails[props.exchangeConnection.exchange].icon
    const canDelete = ref<boolean>(true)

    return {
      icon,
      moment,
      canDelete
    }
  }
})
</script>

<style lang="scss" scoped>
.exchange-connection {
  display: grid;
  align-items: center;
  grid-template-columns: 300px 200px 200px 1fr;
  margin-bottom: 0.5rem;
  background-color: var(--background-lighten);
  border: 1px solid var(--border-color);
  > * {
    margin-right: 1rem;
  }
  p {
    font-size: 14px;
  }
  .name {
    display: flex;
    align-items: center;
    img {
      height: 50px;
      margin-right: 1rem;
    }
    p {
      font-weight: bold;
      font-size: 16px;
    }
  }
  .controls {
    display: flex;
    justify-content: flex-end;
  }
}
</style>
