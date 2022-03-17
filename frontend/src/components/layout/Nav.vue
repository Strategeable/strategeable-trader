<template>
  <div class="navbar">
    <div class="wrapper">
      <img src="@/assets/img/logo.svg" alt="logo"/>
      <div class="bar">
        <div class="inner">
          <div class="left">
            <nav-item
              :active="activeRoute === 'Bots' || activeRoute === 'Bot'"
              name="Bots"
              to="/bots"
            />
            <nav-item
              :active="activeRoute === 'Strategies' || activeRoute === 'Strategy'"
              name="Strategies"
              to="/strategies"
            />
          </div>
          <div class="right">
            <nav-item
              :active="activeRoute === 'Settings'"
              name="Settings"
              to="/settings"
            />
            <color-theme/>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import NavItem from '@/components/layout/NavItem.vue'
import ColorTheme from '@/components/ColorTheme.vue'

// TODO: this should preferably be automatically generated
// using the router config
export default {
  components: { NavItem, ColorTheme },
  computed: {
    activeRoute () {
      return (this as any).$route.name
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 140px;
}
.wrapper {
  height: 100px;
  position: relative;
  width: 100%;
  background-color: var(--primary-darken);
  display: flex;
  align-items: flex-end;
  img {
    width: 200px;
    position: absolute;
    top: 23px;
    left: 50%;
    transform: translateX(-50%);
  }
  .bar {
    z-index: 10;
    width: 100%;
    max-width: var(--container-width);
    margin: 0 1rem;
    margin: 0 auto;
    margin-bottom: -25px;
    .inner {
      display: flex;
      justify-content: space-between;
      margin: 0 1rem;
      background-color: var(--primary);
      box-shadow: 0 -4px 7px rgba(0,0,0,0.1);
      padding: 0 1rem;
    }
    .left, .right {
      display: flex;
    }
  }
  @media(max-width: 600px) {
    .inner {
      flex-direction: column;
      justify-content: unset !important;
    }
    height: unset;
  }
}
</style>
