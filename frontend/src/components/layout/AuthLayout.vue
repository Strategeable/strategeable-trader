<template>
  <div class="auth">
    <div class="inner">
      <img :src="theme === 'dark' ? require('@/assets/img/logo-white.svg') : require('@/assets/img/logo-purple.svg')" alt="logo"/>
      <login v-if="hasUser"/>
      <register v-else/>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'

import Login from '@/components/auth/Login.vue'
import Register from '@/components/auth/Register.vue'
import { useStore } from '@/store'

export default defineComponent({
  components: { Login, Register },
  setup () {
    const store = useStore()
    const theme = computed(() => store.getters.theme)
    const hasUser = computed(() => store.getters.hasUser)
    const showLogin = ref<boolean>(true)

    return {
      showLogin,
      theme,
      hasUser
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
