import { useQuery, useQueryClient } from "@tanstack/react-query";
import { getBookings } from "../../services/apiTasks";
import { useParams, useSearchParams } from "react-router-dom";
import { PAGE_SIZE } from "../../utils/constants";

import getTasks, { getTask } from "../../services/apiTasks";

export function useTasks() {
  const {
    isLoading,
    data: taskData,
    error,
  } = useQuery({
    queryKey: ["tasks"],
    queryFn: getTasks,
  });

  const taskItems = taskData?.data?.data?.tasks;
  const total = taskData?.data?.data?.total;

  return { isLoading, error, taskItems, total };
}

export function useTask() {
  const { taskId } = useParams();
  const {
    isLoading,
    data: taskData,
    error,
  } = useQuery({
    queryKey: ["task", taskId],
    queryFn: () => getTask(taskId),
    retry: false,
  });

  const taskItem = taskData?.data;

  // console.log("taskItem: ", taskItem);
  // console.log(JSON.stringify(taskItem, null, 2));

  return { isLoading, error, taskItem };
}

export function useBookings() {
  const queryClient = useQueryClient();
  const [searchParams] = useSearchParams();

  // FILTER
  const filterValue = searchParams.get("status");
  const filter =
    !filterValue || filterValue === "all"
      ? null
      : { field: "status", value: filterValue };
  // { field: "totalPrice", value: 5000, method: "gte" };

  // SORT
  const sortByRaw = searchParams.get("sortBy") || "startDate-desc";
  const [field, direction] = sortByRaw.split("-");
  const sortBy = { field, direction };

  // PAGINATION
  const page = !searchParams.get("page") ? 1 : Number(searchParams.get("page"));

  // QUERY
  const {
    isLoading,
    data: { data: bookings, count } = {},
    error,
  } = useQuery({
    // 注意这里，使用 filter 作为 key 时，如果变化了，这里也会重新加载数据的
    queryKey: ["bookings", filter, sortBy, page],
    queryFn: () => getBookings({ filter, sortBy, page }),
  });

  // PRE-FETCHING 提前 fetch 当前页面，之前和之后的数据
  const pageCount = Math.ceil(count / PAGE_SIZE);
  if (page < pageCount)
    queryClient.prefetchQuery({
      queryKey: ["bookings", filter, sortBy, page + 1],
      queryFn: () => getBookings({ filter, sortBy, page: page + 1 }),
    });

  if (page > 1)
    queryClient.prefetchQuery({
      queryKey: ["bookings", filter, sortBy, page - 1],
      queryFn: () => getBookings({ filter, sortBy, page: page - 1 }),
    });

  return { isLoading, error, bookings, count };
}
