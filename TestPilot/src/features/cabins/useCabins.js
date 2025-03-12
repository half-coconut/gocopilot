import { useQuery } from "@tanstack/react-query";
import getCabins from "../../services/apiCabins";

export function useCabins() {
  const {
    isLoading,
    data: cabins,
    error,
  } = useQuery({
    queryKey: ["cabins"],
    queryFn: getCabins,
  });

  const cabinItems = cabins?.data?.data?.items;
  const total = cabins?.data?.data?.total;
  // console.log("cabins: ", cabins);
  // console.log("cabinItems: ", cabinItems);
  // console.log("total: ", total);

  return { isLoading, error, cabinItems, total };
}
