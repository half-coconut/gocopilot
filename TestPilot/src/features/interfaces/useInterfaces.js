import { useQuery } from "@tanstack/react-query";
import getInterfaces from "../../services/apiInterfaces";

export function useInterfaces() {
  const {
    isLoading,
    data: cabins,
    error,
  } = useQuery({
    queryKey: ["interfaces"],
    queryFn: getInterfaces,
  });

  const cabinItems = cabins?.data?.data?.interfaces;
  const total = cabins?.data?.data?.total;
  // console.log("cabins: ", cabins);
  // console.log("cabinItems: ", cabinItems);
  // console.log("total: ", total);

  return { isLoading, error, cabinItems, total };
}
