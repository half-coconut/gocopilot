import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { login as loginApi } from "../../services/apiAuth";
import toast from "react-hot-toast";

export function useLogin() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { mutate: login, isLoading: isLoginLoading } = useMutation({
    mutationFn: ({ email, password }) => loginApi({ email, password }),
    onSuccess: (user) => {
      // set 'user' data into react query cache
      queryClient.setQueryData(["user"], user);
      navigate("/dashboard", { replace: true });
      toast.success("Login successful!");
    },
    onError: (err) => {
      console.log("ERROR", err);
      // toast.error("Provided email or password incorrect");
      toast.error(err);
    },
  });

  return { login, isLoginLoading };
}
