import { useQueryClient, useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { createEditTask } from "../../services/apiTasks";

export function useEditTask() {
  const queryClient = useQueryClient();

  const { mutate: editTask, isLoading: isEditing } = useMutation({
    mutationFn: ({ newTaskData }) => createEditTask(newTaskData),
    onSuccess: () => {
      toast.success("Task successfully edited");
      queryClient.invalidateQueries({
        queryKey: ["tasks"],
      });
    },
    onError: (err) => toast.error(err.message),
  });
  return { isEditing, editTask };
}
