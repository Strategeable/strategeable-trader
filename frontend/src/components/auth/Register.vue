<template>
  <div class="register">
    <p>Set a password for your Strategable trader</p>
    <form>
      <div class="input">
        <p>Password</p>
        <input
          type="password"
          v-model="password"
        >
      </div>
      <div class="input">
        <p>Repeat password</p>
        <input
          type="password"
          v-model="repeatPassword"
        >
      </div>
      <button
        :disabled="!valid"
        @click="register"
        type="submit"
      >
        {{ loading ? 'Loading...' : 'Register' }}
      </button>
      <p v-if="errorMsg" class="error">{{ errorMsg }}</p>
    </form>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import { useStore } from '@/store'

import { ActionTypes } from '@/types/store/action-types'

export default defineComponent({
  setup () {
    const store = useStore()
    const password = ref<string>()
    const repeatPassword = ref<string>()
    const valid = computed(() => {
      if (!password.value || password.value === '') return false
      if (repeatPassword.value !== password.value) return false
      return true
    })
    const loading = ref<boolean>(false)
    const errorMsg = ref<string>()

    async function register () {
      loading.value = true
      errorMsg.value = await store.dispatch(ActionTypes.REGISTER_ACCOUNT, {
        password: password.value as string
      }
      )
      loading.value = false
    }

    return {
      password,
      repeatPassword,
      valid,
      errorMsg,
      loading,
      register
    }
  }
})
</script>

<style lang="scss" scoped>
.register {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  p {
    margin-bottom: 1.5rem;
  }
  form {
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
  .input {
    width: 100%;
    margin-bottom: 1rem;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    p {
      margin-bottom: 0.5rem;
      text-align: center;
    }
    input {
      width: 100%;
      padding: 0.5rem;
      border: 1px solid var(--border-color);
      max-width: 300px;
      &:active, &:focus {
        outline: var(--primary) solid 1px;
      }
    }
  }
  .error {
    color: red;
    padding-top: 2rem;
  }
}
</style>
