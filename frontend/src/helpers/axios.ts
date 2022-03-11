import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: process.env.VUE_APP_API_URL
})

axiosInstance.interceptors.request.use(cfg => {
  const token = localStorage.getItem('jwt')
  if (token && cfg && cfg.headers) {
    cfg.headers.Authorization = `Bearer ${token}`
  }

  return cfg
})

export default axiosInstance
