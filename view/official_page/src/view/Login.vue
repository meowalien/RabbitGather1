<template>
  <div class="h-[100vh] w-[100vw] bg-white flex flex-col justify-center items-center">
    <div class="bg-pink-200 h-min-96 w-min-72 rounded-4xl shadow-2rg p-7 flex flex-col justify-between">
      <div class="w-full flex flex-row justify-center">
        <img src="/icon/RabitIcon.svg" alt="" class="w-2/3">
      </div>
      <div class="w-full flex flex-row items-baseline mt-3">
        <p class="min-w-max mr-2 text-xl">&nbsp;Account：</p>
        <input  @keyup.enter="sendLogin(account , password)" v-model="account" type="text" class="h-10 py-3 px-4 bg-white rounded-lg placeholder-gray-400 text-gray-900 appearance-none inline-block w-full shadow-md focus:outline-none focus:ring-2 focus:ring-pink-300">
      </div>
      <div class="w-full flex flex-row items-baseline mt-3">
        <p class="min-w-max mr-2 text-xl">Password：</p>
        <input @keyup.enter="sendLogin(account , password)" v-model="password" type="password" class="h-10 py-3 px-4 bg-white rounded-lg placeholder-gray-400 text-gray-900 appearance-none inline-block w-full shadow-md focus:outline-none focus:ring-2 focus:ring-pink-300">
      </div>
      <div class="w-full flex flex-row items-center justify-end mt-3 ">
        <button @click="sendLogin(account , password)" class="mr-4 text-2xl bg-[#ff9600] rounded-4xl p-1 pr-2 pl-2 font-bold">OK</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, defineComponent } from "vue";
import axios from "axios";
import {StoreKey } from "@/store"
import {NotErrCode, Panic, StdResponse} from "@/module/ErrorHandler";

export default defineComponent({
  name: "Login",
  props: {
  },
  setup: () => {
    const account = ref("");
    const password = ref("");
    return { account , password };
  },
  methods :{
    async sendLogin (account:string , password:string ){
      console.debug("account: ",account)
      console.debug("password: ",password)


      let d = new URL( "/member/login",this.$store.state.api_base_url);

      await axios.post(d.toString(),{
        "account":account,
        "password":password
      }).then(res => {
        if (!NotErrCode(res.data as StdResponse)) {
          return
        }
        this.$store.commit(StoreKey.app_token,res.data.data.token)
        this.$store.commit(StoreKey.refresh_token,res.data.data.refresh_token)
      }).catch(reason => {
        Panic(reason)
      })
    }
  }
});
</script>