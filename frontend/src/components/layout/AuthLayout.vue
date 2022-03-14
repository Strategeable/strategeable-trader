<template>
  <div class="auth">
    <div class="inner">
      <img src="@/assets/img/logo-white.svg" alt="logo"/>
      <login v-if="showLogin"/>
      <register v-else/>
      <p v-if="showLogin" class="or">Or <span @click="showLogin = false">make an account ></span></p>
      <p v-else class="or">Or <span @click="showLogin = true">log in ></span></p>
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
    padding: 2rem 1rem;
    background-color: var(--background-lighten);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    img {
      width: 250px;
      margin-bottom: 2rem;
    }
    .or {
      margin-top: 1.5rem;
      span {
        color: var(--primary);
        cursor: pointer;
        &:hover {
          text-decoration: underline;
        }
      }
    }
  }
}
</style>
