import { supabase, supabaseUrl } from "./supabase";

// 查询 cabin
export default async function getCabins() {
  const { data, error } = await supabase.from("cabins").select("*");

  if (error) {
    console.error(error);
    throw new Error("Cabins could not be loaded");
  }
  return data;
}

// 创建新的 cabin
export async function createEditCabin(newCabin, id) {
  const hasImagePath = newCabin.image?.startsWith?.(supabaseUrl);

  const imageName = `${Math.random()}-${newCabin.image.name}`.replaceAll(
    "/",
    ""
  );
  const imagePath = hasImagePath
    ? newCabin.image
    : `${supabaseUrl}/storage/v1/object/public/cabin-images//${imageName}`;

  // https://nvbtdjgdbhgsgccuwkap.supabase.co/storage/v1/object/public/cabin-images//cabin-001.jpg

  // 1. Create cabin
  let query = supabase.from("cabins");
  // A) Create
  if (!id) query = query.insert([{ ...newCabin, image: imagePath }]);
  // B)Edit
  if (id)
    query = query
      .update({ ...newCabin, image: imagePath })
      .eq("id", id)
      .select();

  const { data, error } = await query.select().single();

  if (error) {
    console.error(error);
    throw new Error("Cabin could not be created");
  }

  // 2. Upload image, uploading file using standard upload
  if (hasImagePath) return data; // 如果已有图片就不用上传了，直接返回

  const { error: storageError } = await supabase.storage
    .from("cabin-images")
    .upload(imageName, newCabin.image);

  // 3. Delele the cabin IF there was an error
  if (storageError) {
    await supabase.from("cabins").delete().eq("id", data.id);
    console.error(storageError);
    throw new Error(
      "Cabin image could not be uploaded and the cabin was not created"
    );
  }

  return data;
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
