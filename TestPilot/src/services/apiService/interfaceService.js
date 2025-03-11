import instance from "./api";

const interfaceService = {
  add: async (data) => {
    try {
      const res = await instance.post("/api/add", data);
      return res;
    } catch (error) {
      console.error("Add API error:", error);
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
};

export default interfaceService;
