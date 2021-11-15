import {InjectionKey} from 'vue'
// @ts-ignore
import {createStore, Store, useStore as baseUseStore} from 'vuex'
import createPersistedState from "vuex-persistedstate";

// define your typings for the store state
export interface State {
    api_base_url: string
    app_token: string
    refresh_token: string
}

// define injection key
export const key: InjectionKey<Store<State>> = Symbol()

export enum StoreKey {
    api_base_url = "api_base_url",
    app_token = "app_token",
    refresh_token = "refresh_token"
}

const mystate: State = {
    // api_base_url: "https://rabbit_gather_api.meowalien.com",
    api_base_url: "http://localhost:2001/",
    refresh_token: "",
    app_token: ""
}


export const store = createStore<State>({
    plugins: [createPersistedState(
        {
            paths: [
                StoreKey.refresh_token,
                StoreKey.app_token,
            ]
        }
    )],
    state: mystate,
    mutations: {
        [StoreKey.refresh_token](state: any, new_refresh_token: string) {
            console.log("new_refresh_token: ", new_refresh_token)
            state.refresh_token = new_refresh_token
        },
        [StoreKey.app_token](state: any, newToken: string) {
            console.log("newToken: ", newToken)
            state.app_token = newToken
        }
    }
})


// define your own `useStore` composition function
export function useStore() {
    return baseUseStore(key)
}

