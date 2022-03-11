<template>
  <div class="login">
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
    >
      Login
    </button>
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

</style>
