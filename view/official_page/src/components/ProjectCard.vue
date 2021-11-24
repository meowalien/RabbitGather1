<template>
<div>fdsfdklfjsdlk</div>
</template>

<script>
// in a vue component
import { useStore } from "@/store"

export default {
  setup () {
    const store = useStore()

    console.log(store.state.count) // typed as number
  }
}
</script>

<style scoped>

</style>