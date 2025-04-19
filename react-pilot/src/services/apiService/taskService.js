import instance from "./api";

const taskService = {
  edit: async (data) => {
    try {
      const res = await instance.post("/task/edit", data);
      return res;
    } catch (error) {
      console.error("Edit Task error:", error);
      throw error;
    }
  },
  list: async () => {
    try {
      const res = await instance.get("/task/list");
      return res;
    } catch (error) {
      console.error("List Task error:", error);
      throw error;
    }
  },
  detail: async (id) => {
    try {
      const res = await instance.get("/task/detail/" + id);
      return res;
    } catch (error) {
      console.error("Task Detail error:", error);
      throw error;
    }
  },

  debugInterfaces: async (id) => {
    try {
      const res = await instance.get("/task/debug/interfaces/" + id);
      return res;
    } catch (error) {
      console.error("Task Debug Interface error:", error);
      throw error;
    }
  },

  debugTask: async (id) => {
    try {
      const res = await instance.get("/task/debug/" + id);
      return res;
    } catch (error) {
      console.error("Task Debug Task error:", error);
      throw error;
    }
  },

  executeTask: async (id) => {
    try {
      const res = await instance.get("/task/execute/" + id);
      return res;
    } catch (error) {
      console.error("Task Debug Task error:", error);
      throw error;
    }
  },
};

export default taskService;
