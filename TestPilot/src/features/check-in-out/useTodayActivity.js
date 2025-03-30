import { useQuery } from "@tanstack/react-query";
import { getStaysTodayActivity } from "../../services/apiTasks";

export function useTodayActivity() {
  const { isLoading, data: activities } = useQuery({
    queryFn: getStaysTodayActivity,
    queryKey: ["today-activity"],
  });

  return { activities, isLoading };
}
