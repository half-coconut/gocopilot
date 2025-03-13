import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { startAskOpenai } from "../../services/apiOpenai"; // 假设这是你的 API 函数

export function useOpenai(ask) {
  // 接收 ask 作为参数
  const queryClient = useQueryClient();
  // 读取初始数据，如果需要
  const { data: initialData, isLoading: isQueryLoading } = useQuery({
    //如果存在，读取缓存
    queryKey: ["openai", ask], // 包含了 ask 的 queryKey
    queryFn: () => null, //  不需要请求数据, 我们只关心缓存是否存在.  可以返回 null.
    enabled: false, // 不要自动执行, 只需要读取之前缓存的数据即可
  });
  const { mutate: askAI, isLoading: isMutationLoading } = useMutation({
    mutationFn: (ask) => startAskOpenai(ask), //在这里传递ask
    onSuccess: (data, variables) => {
      //在变量中找到原来的参数
      toast.success("Ask openai successfully");
      queryClient.setQueryData(["openai", variables], (oldData) => ({
        ...oldData,
        ...data,
      })); //使用函数式更新，避免闭包问题
    },
    onError: (err) => toast.error(err.message),
  });

  return { isLoading: isQueryLoading || isMutationLoading, askAI, initialData }; //initialData 包含了初始数据和 isLoading状态
}
