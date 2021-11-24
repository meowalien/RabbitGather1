<template>
  <div class="w-full bg-white flex flex-col justify-items-start items-center">
    <div

        class=" gap-5 flex flex-col sm:flex-row flex-wrap items-center sm:justify-center bg-black w-full sm:w-11/12 max-w-[1080px] p-5 lg:p-4">
      <ProjectCard v-for="(item, key, index) in project_list" v-bind:text="item" class="ProjectCard">
      </ProjectCard>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from "vue";

import ProjectCard from "../components/ProjectCard.vue"
import {ABC} from "@/components/ProjectCard.vue";

export default defineComponent({
  components: {
    ProjectCard
  },
  name: "Home",
  props: {},
  setup: () => {

    const project_list: ABC[] = [
      {
        title:"hut",
        content:"88888",
      },
    ];
    // const password = ref("");
    return { project_list };
  },
  methods: {}
});
</script>


<style scoped>
.ProjectCard {
  @apply bg-blue-300 h-72 w-full sm:w-72
}
</style>