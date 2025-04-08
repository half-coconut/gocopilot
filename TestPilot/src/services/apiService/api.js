import axios from "axios";

// const API_BASE_URL = process.env.BASE_URL;
// const API_BASE_URL = "http://localhost:3002";
const API_BASE_URL = "http://47.239.187.141:3002/api";

const instance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 1000 * 60 * 5, // 设置超时时间：5min
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

// 请求拦截器
// 在这里让每一个请求都加上 authorization 的头部
instance.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些什么，例如添加 token
    const token = localStorage.getItem("token");
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    // 处理请求错误
    return Promise.reject(error);
  }
);

// 响应拦截器
instance.interceptors.response.use(
  (resp) => {
    const newToken = resp.headers["x-jwt-token"] || null;
    const newRefreshToken = resp.headers["x-refresh-token"] || null;

    if (newToken) {
      localStorage.setItem("token", newToken); // 保存新的 token
    }
    if (newRefreshToken) {
      localStorage.setItem("refresh_token", newRefreshToken);
    }

    if (resp.status === 401) {
      window.location.href = "/users/login";
    }

    // console.log("resp: ", resp);

    return resp;
  },
  (err) => {
    console.error(err);
    if (err.response && err.response.status === 401) {
      window.location.href = "/login";
    }
    return Promise.reject(err);
  }
);

export default instance;
