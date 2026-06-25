import axios from 'axios'

// 创建 axios 实例
const request = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // 指向刚刚写好的 Go 后端
  timeout: 15000
})

// 请求拦截器 (注入 JWT Token)
request.interceptors.request.use(
  config => {
    // 每次发送请求之前，检查本地是否存在 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = 'Bearer ' + token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器 (统一错误处理)
request.interceptors.response.use(
  response => {
    const res = response.data
    // 我们的后端规定了标准响应体结构，code !== 200 即为业务报错
    if (res.code !== 200) {
      console.error("API Error: ", res.msg)
      
      // 如果后端返回 401 鉴权失败，直接踢回登录页
      if (res.code === 401) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      }
      return Promise.reject(new Error(res.msg || 'API Error'))
    }
    // 返回真实的 data 载荷
    return res.data
  },
  error => {
    console.error("Network Error: ", error)
    return Promise.reject(error)
  }
)

export default request
