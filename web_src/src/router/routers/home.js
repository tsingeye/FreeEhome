export default [
    {
        path: '/',
        redirect: '/manager/device'
    },
    {
        path: '/home',
        redirect: '/manager/device'
    },
    {
        path: '*',
        component: () => import(/*webpackChunkName:'404'*/'@/view/404.vue'),
    },
]