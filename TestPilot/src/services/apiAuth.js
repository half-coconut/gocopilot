import supabase from "./supabase";
import PropTypes from "prop-types";
import { supabaseUrl } from "../services/supabase";
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
// 基于 supabase
export async function loginV1({ email, password }) {
  const { data, error } = await supabase.auth.signInWithPassword({
    email: email,
    password: password,
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

export async function getCurrentUser() {
  const { data: session } = await supabase.auth.getSession();
  if (!session.session) return null;

  const { data, error } = await supabase.auth.getUser();

  if (error) throw new Error(error.message);
  return data?.user;
}

export async function logoutV1() {
  const { error } = await supabase.auth.signOut();
  if (error) throw new Error(error.message);
}

// 基于后端服务
export async function logout() {
  const resp = await userService.logout();
  return resp;
}

export async function updateCurrentUser({ password, fullName, avatar }) {
  // 1. Update password OR fullName
  let updateData;
  if (password) updateData = { password };
  if (fullName) updateData = { data: { fullName } };

  const { data, error } = await supabase.auth.updateUser(updateData);

  if (error) throw new Error(error.message);
  if (!avatar) return data;

  // 2. Upload the avatar image
  const fileName = `avatar-${data.user.id}-${Math.random()}`;

  const { error: uploadError } = await supabase.storage
    .from("avatars")
    .upload(fileName, avatar);

  if (uploadError) throw new Error(uploadError.message);

  // 3. Update avatar in the user
  const { data: updatedUser, error: updatedUserError } =
    await supabase.auth.updateUser({
      data: {
        avatar: `${supabaseUrl}/storage/v1/object/public/avatars//${fileName}`,
      },
    });

  if (updatedUserError) throw new Error(updatedUserError.message);

  return updatedUser;
}

updateCurrentUser.propTypes = {
  fullName: PropTypes.string,
  password: PropTypes.string,
  avatar: PropTypes.string,
};
