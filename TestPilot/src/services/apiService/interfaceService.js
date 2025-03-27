import instance from "./api";

const interfaceService = {
  edit: async (data) => {
    try {
      const res = await instance.post("/api/edit", data);
      return res;
    } catch (error) {
      console.error("Edit API error:", error);
      throw error;
    }
  },

  list: async () => {
    try {
      const res = await instance.get("/api/list");
      return res;
    } catch (error) {
      console.error("List API error:", error);
      throw error;
    }
  },
  detail: async (id) => {
    try {
      const res = await instance.get("/api/detail/" + id);
      return res;
    } catch (error) {
      console.error("List API error:", error);
      throw error;
    }
  },
};

export default interfaceService;
