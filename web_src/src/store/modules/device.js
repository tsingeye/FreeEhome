import { deviceList,channelList } from '@/api/device'



export default {
    state: {
        loading: false,
        deviceConfig: {
            deviceList: []
        },
        channelConfig:{
            channelList:[]
        }
    },
    mutations: {
        setDevice (state, has) {
            state.deviceConfig = has
        },
        setLoading (state, has) {
            state.loading = has
        },
        setChannel(state, has){
            state.channelConfig = has
        }
    },
    actions: {
        getDeviceList ({ commit, state }, options) {
            commit('setLoading', true)
            deviceList(options)
                .then(res => {
                    commit('setLoading', false)
                    commit('setDevice', res)
                })
        },
        getChannelList ({ commit, state }, options) {
            console.log(options)
            commit('setLoading', true)
            channelList(options)
                .then(res => {
                    commit('setLoading', false)
                    commit('setChannel', res)
                })
        }
    }
}
