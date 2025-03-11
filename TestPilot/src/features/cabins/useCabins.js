import { useQuery } from "@tanstack/react-query";
import getCabins from "../../services/apiCabins";

export function useCabins() {
  let {
    isLoading,
    data: cabins,
    error,
  } = useQuery({
    queryKey: ["cabins"],
    queryFn: getCabins,
  });

  cabins = cabins?.data;
  console.log("cabins: ", cabins);
  console.log("cabins data: ", cabins?.data);

  const cabinsData = cabins?.data;

  return { isLoading, error, cabinsData };
}
