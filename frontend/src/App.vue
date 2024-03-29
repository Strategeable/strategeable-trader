<template>
  <div id="loader" v-if="loading">
    <!-- Still needs a nice loading screen -->
    Loading.
  </div>
  <auth-layout
    v-else-if="!isLoggedIn"
  />
  <div v-else>
    <navbar/>
    <div id="content">
      <router-view/>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref } from 'vue'
import jwtDecode from 'jwt-decode'

import Navbar from '@/components/layout/Nav.vue'
import AuthLayout from '@/components/layout/AuthLayout.vue'
import { useStore } from '@/store'
import { ActionTypes } from '@/types/store/action-types'
import { MutationTypes } from '@/types/store/mutation-types'
import { Theme } from './types/general'

export default defineComponent({
  components: { Navbar, AuthLayout },
  setup () {
    const store = useStore()
    const isLoggedIn = computed(() => store.getters.loggedIn)
    const loading = ref(true)

    const init = async () => {
      store.dispatch(ActionTypes.CHANGE_COLOR_THEME, localStorage.getItem('theme') as Theme)

      const hasUser = await store.dispatch(ActionTypes.CHECK_AUTH_STATE, undefined)
      if (!hasUser) {
        loading.value = false
        return
      }

      const token = localStorage.getItem('jwt')
      if (!token) {
        loading.value = false
        return
      }

      const decoded: any = jwtDecode(token)
      if (decoded && decoded.exp && Date.now() / 1000 < decoded.exp) {
        store.commit(MutationTypes.SET_JWT, token)
        store.dispatch(ActionTypes.INIT, undefined)
      }

      loading.value = false
    }

    onMounted(() => {
      init()
    })

    return {
      isLoggedIn,
      loading
    }
  },
  sockets: {
    connect () {
      const vm: any = this as any
      if (!vm.$store.getters.loggedIn) return

      vm.$socket.emit('authorization', localStorage.getItem('jwt'))
    }
  }
})
</script>

<style lang="scss" scoped>
#content {
  width: 100%;
  max-width: var(--container-width);
  margin: 0 auto;
  padding: 0 1rem;
  padding-top: 2rem;
}
</style>
