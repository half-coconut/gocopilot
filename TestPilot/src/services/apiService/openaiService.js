import instance from "./api";

const openaiService = {
  ask: async (data) => {
    try {
      const res = await instance.post("/openai/ask/deepseek", data);
      return res;
    } catch (error) {
      console.error("Openai API error:", error);
      throw error;
    }
  },
};
export default openaiService;
