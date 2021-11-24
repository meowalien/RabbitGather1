import {createApp} from 'vue'
import App from './App.vue'
import './index.css'
import {key, store} from './store'
import router from './router'
// import axios from 'axios'
//
// import vuerouter from 'vue-router'
//
// axios.interceptors.request.use(async function (config) {
//     return config
// }, function (error) {
//     return Promise.reject(error)
// })
//
// axios.interceptors.response.use(async function (response) {
//     return response
// }, function (error) {
//     return Promise.reject(error)
// })
//
//
// // function hasPermission(roles: string[], permissionRoles: string) {
// //     if (roles.indexOf('admin') >= 0) return true // admin permission passed directly
// //     if (!permissionRoles) return true
// //     return roles.some(role => permissionRoles.indexOf(role) >= 0)
// // }
//
// function hasPermission(router: vuerouter.RouteLocationNormalized, accessMenu) {
//     if (whiteList.indexOf(router.path) !== -1) {
//         return true;
//     }
// }
//
// function getToken(): string {
//     return store.state.app_token
// }
//
// const whiteList = ['/login'] // no redirect whitelist
// router.beforeEach(async (to, from, next) => {
//     if (getToken()) {
//         let userInfo = store.state.user.userInfo;
//         if (!userInfo.name) {
//             try {
//                 await store.dispatch("GetUserInfo")
//                 await store.dispatch('updateAccessMenu')
//                 if (to.path === '/login') {
//                     next({name: 'home_index'})
//                 } else {
// //Util.toDefaultPage([...routers], to.name, router, next);
//                     next({...to, replace: true})//菜單權限更新完成,重新進一次當前路由
//                 }
//             } catch (e) {
//                 if (whiteList.indexOf(to.path) !== -1) { // 在免登錄白名單，直接進入
//                     next()
//                 } else {
//                     next('/login')
//                 }
//             }
//         } else {
//             if (to.path === '/login') {
//                 next({name: 'home_index'})
//             } else {
//                 if (hasPermission(to, store.getters.accessMenu)) {
//                     Util.toDefaultPage(store.getters.accessMenu, to, routes, next);
//                 } else {
//                     next({path: '/403', replace: true})
//                 }
//             }
//         }
//     } else {
//         if (whiteList.indexOf(to.path) !== -1) { // 在免登錄白名單，直接進入
//             next()
//         } else {
//             next('/login')
//         }
//     }
// })
//
//
// router.afterEach(() => {
//     // NProgress.done() // finish progress bar
// })


const app = createApp(App)

app.use(store, key)
app.use(router)

app.mount('#app')
