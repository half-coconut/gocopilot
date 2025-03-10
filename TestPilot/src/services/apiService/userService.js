import instance from "./api";

const userService = {
  hello: () => {
    return instance.get("/hello");
  },
  getUsers: () => {
    return instance.get("/users");
  },

  getUser: (id) => {
    return instance.get(`/users/${id}`);
  },

  login: async (data) => {
    try {
      const res = await instance.post("/users/login", data);
      return res; // 返回整个 response 对象
    } catch (error) {
      console.error("Login API error:", error);
      throw error; // 重新抛出错误，让调用者处理
    }
  },

  logout: async () => {
    try {
      const res = await instance.get("/users/logout");
      localStorage.removeItem("token");
      localStorage.removeItem("refresh_token");
      return res; // 返回整个 response 对象
    } catch (error) {
      console.error("Login API error:", error);
      throw error; // 重新抛出错误，让调用者处理
    }
  },

  signup: (data) => {
    return instance.post("/users/signup", data);
  },

  editUser: (data) => {
    return instance.put(`/users/edit`, data);
  },

  profile: () => {
    return instance.get(`/users/profile`);
  },
};

export default userService;
