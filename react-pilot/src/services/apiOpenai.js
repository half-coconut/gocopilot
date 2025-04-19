import opensiService from "./apiService/openaiService";

export async function startAskOpenai(newAsk) {
  return opensiService.ask(newAsk);
}
