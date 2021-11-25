import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'login',
    meta: {
      title: '登录'
    },
    component: () => import('@/views/login.vue')
  },
  {
    path: '/home',
    name: 'home',
    component: () => import('@/views/home.vue'),
    redirect: '/monitor',
    children: [
      {
        path: '/monitor',
        name: 'monitor',
        meta: {
          title: '设备点播'
        },
        component: () => import('@/views/monitor.vue')
      },
      {
        path: '/devices',
        name: 'devices',
        meta: {
          title: '设备列表'
        },
        component: () => import('@/views/devices.vue')
      },
      {
        path: '/channels',
        name: 'channels',
        meta: {
          title: '通道列表'
        },
        component: () => import('@/views/channels.vue')
      }
    ]
  }
]

const router = new VueRouter({
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = 'Ehome | ' + to.meta.title
  }
  next()
})

export default router
