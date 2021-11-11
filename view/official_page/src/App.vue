<template>
  <router-view />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
// import Home from '@view/Home.vue'
import Try from '@/components/ForTest.vue'

export default defineComponent({
  name: 'App',
  components: {
    // Home
  }
})
</script>

<style lang="postcss">
  #app {
    @apply text-center
  }
</style>