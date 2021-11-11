<template>
  <div class="h-[100vh] w-[100vw] bg-white flex flex-col justify-center items-center">

    <div class="bg-white h-96 w-72 rounded-4xl shadow-2rg p-7 flex flex-col justify-between">
      <div>
        <img src="/favicon.ico" alt="" class="block">
      </div>
      <div class="w-full flex flex-row items-baseline">
        <p class="min-w-max mr-2"> 帳號：</p><input type="text" class="py-3 px-4 bg-white rounded-lg placeholder-gray-400 text-gray-900 appearance-none inline-block w-full shadow-md focus:outline-none focus:ring-2 focus:ring-pink-300">
      </div>
      <div class="w-full flex flex-row items-baseline">
        <p class="min-w-max mr-2"> 帳號：</p><input type="text" class="py-3 px-4 bg-white rounded-lg placeholder-gray-400 text-gray-900 appearance-none inline-block w-full shadow-md focus:outline-none focus:ring-2 focus:ring-pink-300">
      </div>
    </div>
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