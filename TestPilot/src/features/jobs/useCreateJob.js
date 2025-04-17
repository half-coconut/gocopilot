import { useQueryClient, useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { createEditJob } from "../../services/apiJobs";

export function useCreateJob() {
  const queryClient = useQueryClient();

  const { mutate: createJob, isLoading: isCreating } = useMutation({
    mutationFn: createEditJob,
    onSuccess: () => {
      toast.success("New job successfully created");
      queryClient.invalidateQueries({
        queryKey: ["job"],
      });
    },
    onError: (err) => toast.error(err.message),
  });

  return { isCreating, createJob };
}
