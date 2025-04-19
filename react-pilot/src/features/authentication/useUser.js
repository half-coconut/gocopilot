import { useQuery } from "@tanstack/react-query";
import { getUserProfile } from "../../services/apiAuth";

export function useUser() {
  const { isLoading, data: user } = useQuery({
    queryKey: ["user"],
    queryFn: getUserProfile,
  });

  const userData = user?.data?.data;

  // 这里对是否为合法用户进行校验
  return { isLoading, userData, isAuthenticated: user?.status === 200 };
}
