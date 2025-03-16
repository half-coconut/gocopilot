import { useQuery } from "@tanstack/react-query";
import getInterfaces from "../../services/apiInterfaces";

export function useInterfaces() {
  const {
    isLoading,
    data: interfaceData,
    error,
  } = useQuery({
    queryKey: ["interfaces"],
    queryFn: getInterfaces,
  });

  const interfaceItems = interfaceData?.data?.data?.interfaces;
  const total = interfaceData?.data?.data?.total;

  return { isLoading, error, interfaceItems, total };
}
