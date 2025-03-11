import instance from "./api";

const userService = {
  getUsers: () => {
    return instance.get("/users");
  },

  getUser: (id) => {
    return instance.get(`/users/${id}`);
  },

  login: async (data) => {
    try {
      const res = await instance.post("/users/login", data);
      return res;
    } catch (error) {
      console.error("Login API error:", error);
      throw error;
    }
  },

  logout: async () => {
    try {
      const res = await instance.get("/users/logout");
      localStorage.removeItem("token");
      localStorage.removeItem("refresh_token");
      return res;
    } catch (error) {
      console.error("Logout API error:", error);
      throw error;
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
