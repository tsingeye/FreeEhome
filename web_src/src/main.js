import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store'
import './plugins/element'
// import './test.scss'

Vue.config.productionTip = false
/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
