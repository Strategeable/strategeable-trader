<template>
  <auth-layout v-if="!isLoggedIn"/>
  <div v-else>
    <navbar/>
    <div id="content">
      <router-view/>
    </div>
  </div>
</template>

<script lang="ts">
import { useStore } from 'vuex'
import { computed, onMounted } from 'vue'
import jwtDecode from 'jwt-decode'

import Navbar from '@/components/layout/Nav.vue'
import AuthLayout from '@/components/layout/AuthLayout.vue'

export default {
  components: { Navbar, AuthLayout },
  setup () {
    const store = useStore()
    const isLoggedIn = computed(() => store.getters.loggedIn)

    onMounted(() => {
      const token = localStorage.getItem('jwt')
      if (!token) return

      const decoded: any = jwtDecode(token)
      if (decoded && decoded.exp && Date.now() / 1000 < decoded.exp) {
        store.commit('SET_JWT', token)
        store.dispatch('init')
      }
    })

    return {
      isLoggedIn
    }
  }
}
</script>

<style lang="scss" scoped>
#content {
  width: 100%;
  max-width: var(--container-width);
  margin: 0 auto;
  padding: 0 1rem;
}
</style>
