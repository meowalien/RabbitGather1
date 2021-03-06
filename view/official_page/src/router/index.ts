import {createRouter, createWebHistory, RouteRecordRaw} from "vue-router";
// import Login from "../view/Login.vue"

const routes: Array<RouteRecordRaw> = [
    {
        path: "/login",
        name: "Login",
        component:() => import('../view/Login.vue'),// Login,
        meta: {
            roles: ['admin'],
        }
    },
    {
        path: "/",
        name: "Home",
        component:() => import('../view/Home.vue'),// Login,
        meta: {
            roles: ['admin'],
        }
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;