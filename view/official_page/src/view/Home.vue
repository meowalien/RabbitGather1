<template>
  <div class="h-[100vh] w-[100vw] bg-blue-300">

  </div>

</template>

<script lang="ts">
import { ref, defineComponent } from "vue";
export default defineComponent({
  name: "Home",
  props: {
    // msg: {
    //   type: String,
    //   required: true,
    // },
  },
  setup: () => {
    const count = ref(200);
    return { count };
  },
});
</script>