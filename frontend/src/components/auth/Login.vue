<template>
  <div class="login">
    <form @submit="login">
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

<script>
import { ActionTypes } from '@/types/store/action-types'

export default {
  data () {
    return {
      password: ''
    }
  },
  computed: {
    valid () {
      if (!this.password) return false
      return true
    }
  },
  methods: {
    async login () {
      this.$store.dispatch(ActionTypes.LOGIN, { password: this.password })
    }
  }
}
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
