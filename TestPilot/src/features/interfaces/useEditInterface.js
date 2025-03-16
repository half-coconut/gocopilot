import { useQueryClient, useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { createEditInterface } from "../../services/apiInterfaces";

export function useEditInterface() {
  const queryClient = useQueryClient();

  const { mutate: editInterface, isLoading: isEditing } = useMutation({
    mutationFn: ({ newInterfaceData }) => createEditInterface(newInterfaceData),
    onSuccess: () => {
      toast.success("Interface successfully edited");
      queryClient.invalidateQueries({
        queryKey: ["interfaces"],
      });
    },
    onError: (err) => toast.error(err.message),
  });
  return { isEditing, editInterface };
}
