import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import getInterfaces, {
  getInterfaceDetail,
} from "../../services/apiInterfaces";

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

export function useInterface() {
  const { interfaceId } = useParams();

  const {
    isLoading,
    data: interfaceData,
    error,
  } = useQuery({
    queryKey: ["interface", interfaceId],
    queryFn: () => getInterfaceDetail(interfaceId),
    retry: false, // react 失败了会 retry，但是有时候也不需要
  });

  const interfaceItem = interfaceData?.data?.data;
  const debug_result = interfaceData?.data?.data?.debug_result;

  return { isLoading, error, interfaceItem, debug_result };
}
