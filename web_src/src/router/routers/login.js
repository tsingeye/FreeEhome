export default [
    {
        path: '/login',
        component: () => import(/*webpackChunkName:'home'*/'@/view/login/login.vue'),
        meta: {
            needLogin: true
        }
    },
]