import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateBooking } from "../../services/apiTasks";
import toast from "react-hot-toast";

export function useCheckout() {
  const quentClient = useQueryClient();

  const { mutate: checkout, isLoading: isCheckingOut } = useMutation({
    mutationFn: (bookingId) =>
      updateBooking(bookingId, {
        status: "checked-out",
      }),
    onSuccess: (data) => {
      toast.success(`Booking #${data.id} successfully checked out`);
      quentClient.invalidateQueries({ active: true });
    },
    onError: () => toast.error("There was an error while checking in"),
  });

  return { checkout, isCheckingOut };
}
