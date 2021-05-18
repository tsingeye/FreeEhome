export default [
    {
        path: '/manager',
        component: () => import(/*webpackChunkName:'manager'*/'@/view/manager/index.vue'),
        children:[
            {
                path: '/manager/device',
                component: () => import(/*webpackChunkName:'manager'*/'@/view/manager/device.vue')
            },
            {
                path: '/manager/channel',
                component: () => import(/*webpackChunkName:'manager'*/'@/view/manager/channel.vue')
            },
            {
                path: '/manager/video',
                component: () => import(/*webpackChunkName:'manager'*/'@/view/manager/video.vue')
            },
        ]
    },
  
]