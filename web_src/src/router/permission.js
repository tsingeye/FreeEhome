export default [
    {
        path: 'setting',
        meta: {
            auth: 'setting'
        },
        name: 'setting',
        component: () => import( /*webpackChunkName:'manager'*/ '@/views/manager/setting.vue')
    },
    {
        path: 'user',
        meta: {
            auth: 'user'
        },
        name: 'user',
        component: () => import( /*webpackChunkName:'manager'*/ '@/views/manager/user.vue')
    },
]