<template>
  <div class="exchange-connection">
    <div class="name">
      <img :src="icon" alt="logo">
      <p>{{ exchangeConnection.name }}</p>
    </div>
    <p><strong>3</strong> (1 running)</p>
    <p>{{ moment(exchangeConnection.createdOn).format('DD MMM YYYY') }}</p>
  </div>
</template>

<script lang="ts">
import { ExchangeConnection, exchangeDetails } from '@/types/Exchange'
import { defineComponent, PropType } from 'vue'
import moment from 'moment'

export default defineComponent({
  props: {
    exchangeConnection: {
      type: Object as PropType<ExchangeConnection>,
      required: true
    }
  },
  setup (props) {
    const icon = exchangeDetails[props.exchangeConnection.exchange].icon

    return {
      icon,
      moment
    }
  }
})
</script>

<style lang="scss" scoped>
.exchange-connection {
  display: flex;
  align-items: center;
  margin-bottom: 0.5rem;
  background-color: var(--background-lighten);
  border: 1px solid var(--border-color);
  > * {
    margin-right: 1rem;
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
    }
  }
}
</style>
