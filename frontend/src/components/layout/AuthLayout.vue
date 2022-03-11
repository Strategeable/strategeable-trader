<template>
  <div class="auth">
    <div class="inner">
      <login v-if="showLogin"/>
      <register v-else/>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

import Login from '@/components/auth/Login.vue'
import Register from '@/components/auth/Register.vue'

export default defineComponent({
  components: { Login, Register },
  setup () {
    const showLogin = ref<boolean>(true)

    return {
      showLogin
    }
  }
})
</script>

<style lang="scss" scoped>
.auth {
  max-width: var(--container-width);
  margin: 0 auto;
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  .inner {
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 500px;
    margin: 0 1rem;
    padding: 1rem;
    background-color: var(--background-lighten);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
}
</style>
