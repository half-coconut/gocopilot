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
};

export default taskService;
