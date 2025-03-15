import { useQueryClient, useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { createEditInterface } from "../../services/apiInterfaces";

export function useCreateInterface() {
  const queryClient = useQueryClient();

  const { mutate: createInterface, isLoading: isCreating } = useMutation({
    mutationFn: createEditInterface,
    onSuccess: () => {
      toast.success("New interface successfully created");
      queryClient.invalidateQueries({
        queryKey: ["interfaces"],
      });
    },
    onError: (err) => toast.error(err.message),
  });

  return { isCreating, createInterface };
}
