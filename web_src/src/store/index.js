import Vue from 'vue'
import Vuex from 'vuex'
import rootModule from './rootModule'

Vue.use(Vuex)

const files = require.context('./modules', false, /\.js$/)
let module = rootModule.modules = (rootModule.modules || {})
files.keys().forEach(key => {
    let moduleName = key.replace(/\.\//, '').replace(/\.js/, '')
    let store = files(key).default
    module[moduleName] = store
    module[moduleName].namespaced = true
})
export default new Vuex.Store(rootModule)
