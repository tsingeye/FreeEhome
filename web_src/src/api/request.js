import axios from 'axios'

axios.defaults.baseURL = 'http://192.168.31.254:8080'

// axios.interceptors.request.use(
//   function(config) {
//     const token = window.localStorage.getItem('token')
//     if (token) {
//       config.headers.token = token
//     }
//     return config
//   },
//   function(error) {
//     return Promise.reject(error)
//   }
// )

export default axios
