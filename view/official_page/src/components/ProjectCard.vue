<template>
  <div>
    title:{{  text.title }}
    content:{{text.content}}
  </div>
</template>

<script lang="ts">
// in a vue component
// import { useStore } from "@/store"
import {defineComponent, PropType} from "vue";
import {scalarOptions} from "yaml";
import Str = scalarOptions.Str;



export declare interface ABC  {
  title: string
  content:string
}


export default defineComponent({
  name: "ProjectCard",
  props: {
    text: {
      type:Object as PropType<ABC>,
      required: true
    }
  },
  setup() {
    // const store = useStore()

    // console.log(store.state.count) // typed as number
  }
});
</script>

<style scoped>

</style>