import { useQueryClient, useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { deleteInterface as deleteInterfaceApi } from "../../services/apiInterfaces";

export function useDeleteInterface() {
  const queryClient = useQueryClient();

  const { isLoading: isDeleting, mutate: deleteInterface } = useMutation({
    // mutationFn: (id) => deleteCabin(id),
    mutationFn: deleteInterfaceApi,
    onSuccess: () => {
      toast.success("Interface successfully deleted");
      queryClient.invalidateQueries({
        queryKey: ["interfaces"],
      });
    },
    onError: (err) => toast.error(err.message),
  });

  return { isDeleting, deleteInterface };
}
