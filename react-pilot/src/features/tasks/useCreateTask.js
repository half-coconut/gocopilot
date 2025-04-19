import { useQueryClient, useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { createEditTask } from "../../services/apiTasks";

export function useCreateTask() {
  const queryClient = useQueryClient();

  const { mutate: createTask, isLoading: isCreating } = useMutation({
    mutationFn: createEditTask,
    onSuccess: () => {
      toast.success("New task successfully created");
      queryClient.invalidateQueries({
        queryKey: ["tasks"],
      });
    },
    onError: (err) => toast.error(err.message),
  });

  return { isCreating, createTask };
}
