<template>
  <div class="navbar">
    <div class="wrapper">
      <div class="bar">
        <div class="inner">
          <div class="left">
            <nav-item
              name="Bots"
              to="/bots"
            />
            <nav-item
              name="Strategies"
              to="/strategies"
            />
          </div>
          <div class="right">
            <nav-item
              name="Settings"
              to="/settings"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import NavItem from '@/components/layout/NavItem.vue'

export default {
  components: { NavItem }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 140px;
}
.wrapper {
  height: 80px;
  position: relative;
  width: 100%;
  background-color: var(--primary-darken);
  display: flex;
  align-items: flex-end;
  .bar {
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
}
</style>
