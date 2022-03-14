<template>
  <div class="login">
    <form @submit="login">
      <div class="input">
        <p>Username</p>
        <input
          type="text"
          v-model="username"
        >
      </div>
      <div class="input">
        <p>Password</p>
        <input
          type="password"
          v-model="password"
        >
      </div>
      <button
        :disabled="!valid"
        @click="login"
        type="submit"
      >
        Login
      </button>
    </form>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import { useStore } from 'vuex'

export default defineComponent({
  setup () {
    const store = useStore()
    const username = ref<string>()
    const password = ref<string>()
    const valid = computed(() => {
      if (!username.value || !password.value) return false
      if (username.value.length <= 2) return false
      return true
    })

    function login () {
      store.dispatch('login', { username: username.value, password: password.value })
    }

    return {
      username,
      password,
      valid,
      login
    }
  }
})
</script>

<style lang="scss" scoped>
.login {
  width: 100%;
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
}
</style>
