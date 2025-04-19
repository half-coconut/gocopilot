import supabase, { supabaseUrl } from "./supabase";
import PropTypes from "prop-types";
import userService from "./apiService/userService";

export async function signup({ fullName, email, password }) {
  const { data, error } = await supabase.auth.signUp({
    email: email,
    password: password,
    options: {
      data: { fullName: fullName, avatar: "" },
    },
  });

  if (error) throw new Error(error.message);

  return data;
}

// 基于后端服务
export async function login({ email, password }) {
  const resp = await userService.login({ email, password });
  return resp;
}

login.propTypes = {
  email: PropTypes.string.isRequired,
  password: PropTypes.string.isRequired,
};

export async function getUserProfile() {
  const resp = await userService.profile();

  return resp;
}

// 基于后端服务
export async function logout() {
  const resp = await userService.logout();
  return resp;
}

// export async function updateCurrentUserV1({ password, fullName, avatar }) {
//   // 1. Update password OR fullName
//   let updateData;
//   if (password) updateData = { password };
//   if (fullName) updateData = { data: { fullName } };

//   const { data, error } = await supabase.auth.updateUser(updateData);

//   if (error) throw new Error(error.message);
//   if (!avatar) return data;

// 2. Upload the avatar image
// const fileName = `avatar-${data.user.id}-${Math.random()}`;

// const { error: uploadError } = await supabase.storage
//   .from("avatars")
//   .upload(fileName, avatar);

// if (uploadError) throw new Error(uploadError.message);

// // 3. Update avatar in the user
// const { data: updatedUser, error: updatedUserError } =
//   await supabase.auth.updateUser({
//     data: {
//       avatar: `${supabaseUrl}/storage/v1/object/public/avatars//${fileName}`,
//     },
//   });

// if (updatedUserError) throw new Error(updatedUserError.message);

// return updatedUser;
// }

export async function updateCurrentUser(updateData) {
  // 1. Update password OR fullName
  console.log("updateCurrentUser接口入参: ", updateData);
  if (!updateData.avatar) return await userService.editUser(updateData);

  // 2. Upload the avatar image
  const fileName = `avatar-${updateData.id}-${Math.random()}`;

  const { error: uploadError } = await supabase.storage
    .from("avatars")
    .upload(fileName, updateData.avatar);

  if (uploadError) throw new Error(uploadError.message);

  // 3. Update avatar in the user

  return await userService.editUser({
    ...updateData,
    avatar: `${supabaseUrl}/storage/v1/object/public/avatars//${fileName}`,
  });
}

updateCurrentUser.propTypes = {
  updateData: PropTypes.node,
};
