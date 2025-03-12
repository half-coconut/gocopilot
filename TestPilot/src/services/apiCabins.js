import interfaceService from "./apiService/interfaceService";
import supabase from "./supabase";

export default async function getCabins() {
  const resp = await interfaceService.list();
  return resp;
}

export async function createEditCabin(newInterface) {
  return interfaceService.edit(newInterface);
}

// 删除 cabin
export async function deleteCabin(id) {
  const { data, error } = await supabase.from("cabins").delete().eq("id", id);

  if (error) {
    console.error(error);
    throw new Error("Cabin could not be deleted");
  }

  return data;
}
