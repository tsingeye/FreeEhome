import Vue from 'vue'
import VueRouter from 'vue-router'
import hooks from './hook'

Vue.use(VueRouter)

const files = require.context('./routers', false, /\.js$/);
let routes = [];
files.keys().forEach(key => {
  routes.push(...files(key).default)
})

let router = new VueRouter({
  base: '',
  routes
})

// 需要给路由增加多个钩子 每个钩子实现一个具体功能 beforeEach next
Object.values(hooks).forEach(hook=>{
  // 绑定hook中的this是路由的实例  
  router.beforeEach(hook.bind(router));
})


export default router
