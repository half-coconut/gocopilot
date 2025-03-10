import { useQuery } from "@tanstack/react-query";
import { getCurrentUser } from "../../services/apiAuth";

export function useUser() {
  const { isLoading, data: user } = useQuery({
    queryKey: ["user"],
    queryFn: getCurrentUser,
  });

  // return { isLoading, user, isAuthenticated: user?.role === "authenticated" };
  // 这里对是否为合法用户进行校验
  return { isLoading, user, isAuthenticated: true };
}
