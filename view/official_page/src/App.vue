<template>
  <router-view class="h-full" />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import HelloWorld from './components/HelloWorld.vue'
import Try from './components/ForTest.vue'

export default defineComponent({
  name: 'App',
  components: {
    Try,
    HelloWorld
  }
})
</script>

<style lang="postcss">
  #app {
    @apply text-center
  }
</style>