// 对axios 二次封装

import axios from 'axios';
import store from '../store/index'
import router from '@/router/index'

import { Loading } from 'element-ui';
import {getLocal} from './local'

let loadingInstance;

// 当页面切换时 删除不必要的请求
class Http{
    constructor(){
        // this.timeout = 3000; // 超时时间
        console.log(process.env.NODE_ENV)
        //this.baseURL = process.env.NODE_ENV=='development'?'/':'http://www.tsingeye.com:8080';
        this.baseURL = '/';

    }
    mergeOptions(options){ // 合并参数
        return {
            timeout:this.timeout,
            baseURL:this.baseURL,
            ...options
        }
    }
    setInterceptor(instance,url){
        instance.interceptors.request.use((config)=>{
        //    console.log(config)
            // config.headers.authorization = 'Bearer ' + getLocal('token');
            return config;
        });
        instance.interceptors.response.use(res=>{
            if(res.status == 200){
                if(res.data.errCode == 401){
                    console.log(router)
                    router.push('/login')
                    return Promise.reject(res.data);
                }
                return Promise.resolve(res.data);
            }
            else if(res.status == 401){
                router.push('/login')
            }
            else{
                // 401 403 .... switch-case 去判断每个状态码代表的含义
                // ...
                return Promise.reject(res);
            }
           
            return res.data
        },err=>{ 
            return Promise.reject(err);
        })
    }
    request(options){ // 用户的参数 + 默认参数 = 总共的参数
        const opts = this.mergeOptions(options);
        const axiosInstance = axios.create(); // axios()
        // 添加拦截器
        this.setInterceptor(axiosInstance,opts.url);
        // 当调用axios.request 时 内部会创建一个 axios实例 并且给这个实例传入配置属性
        return axiosInstance(opts)
    }
    // 这两个方法只是对request方法 一个简写而已
    get(url,params = {}){ // params
        ///api/v1/system/login
        params.authCode = getLocal('authCode')
        return this.request({
            url,
            method:'get',
            params
            // ...config
        })
    }
    post(url,data){  // data
        // 对data进行格式化
        data.authCode = getLocal('authCode')
        return this.request({
            method:'post',
            url,
            data
        })
    }
}
export default new Http

