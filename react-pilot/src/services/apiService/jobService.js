import instance from "./api";

const jobService = {
  edit: async (data) => {
    try {
      const res = await instance.post("/job/add", data);
      return res;
    } catch (error) {
      console.error("Edit Task error:", error);
      throw error;
    }
  },
  open: async (id) => {
    try {
      const res = await instance.get("/open/" + id);
      return res;
    } catch (error) {
      console.error("List Task error:", error);
      throw error;
    }
  },
  close: async (id) => {
    try {
      const res = await instance.get("/close/" + id);
      return res;
    } catch (error) {
      console.error("List Task error:", error);
      throw error;
    }
  },
};

export default jobService;
