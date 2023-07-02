import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        name: 'TerminalPage',
        component: () => import('../views/TerminalPage.vue')
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})


export default router