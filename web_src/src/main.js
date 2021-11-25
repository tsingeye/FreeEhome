import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import '@/style/index'
import fullscreen from 'vue-fullscreen'

Vue.config.productionTip = false
Vue.use(ElementUI, { size: 'mini' })
Vue.use(fullscreen)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
