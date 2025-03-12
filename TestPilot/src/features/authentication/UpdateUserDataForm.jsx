import { useForm } from "react-hook-form";
import PropTypes from "prop-types";

import Button from "../../ui/Button.jsx";
import FileInput from "../../ui/FileInput.jsx";
import Form from "../../ui/Form.jsx";
import FormRow from "../../ui/FormRow";
import Input from "../../ui/Input";
import Textarea from "../../ui/Textarea.jsx";

import { useUpdateUser } from "./useUpdateUser.js";

function UpdateUserDataForm({ userToEidt = {}, onCloseModal }) {
  const { id: editId, ...userData } = userToEidt;

  const { updateUser, isUpdating } = useUpdateUser();
  const isEditSession = Boolean(editId);

  const { register, handleSubmit, reset, formState } = useForm({
    defaultValues: isEditSession ? userData : {},
  });

  const { errors } = formState;

  function onSubmit(data) {
    // const image = typeof data.image === "string" ? data.image : data.image[0];

    console.log("onSbmit data: ", data);

    if (isEditSession)
      updateUser(
        { newUpdateData: { ...data, id: editId } },
        {
          onSuccess: () => {
            reset(), onCloseModal?.();
          },
        }
      );
  }

  function onError(errors) {
    console.log(errors);
  }

  // const isPhoneAvailable =
  //   userData.phone !== undefined &&
  //   userData.phone !== null &&
  //   userData.phone !== "";

  return (
    <Form
      onSubmit={handleSubmit(onSubmit, onError)}
      type={onCloseModal ? "modal" : "regular"}
    >
      <FormRow label="Email address">
        <Input
          value={userData.email}
          disabled
          {...register("email", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Full name" error={errors?.name?.message}>
        <Input
          type="text"
          id="fullName"
          disabled={isUpdating}
          {...register("fullName", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Phone" error={errors?.name?.message}>
        <Input
          type="text"
          id="phone"
          // value={userData.phone}
          disabled={isUpdating}
          {...register("phone", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Department" error={errors?.name?.message}>
        <Input
          type="text"
          id="department"
          disabled={isUpdating}
          {...register("department", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Role" error={errors?.name?.message}>
        <Input
          type="text"
          id="role"
          disabled={isUpdating}
          {...register("role", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Description" error={errors?.name?.message}>
        <Textarea
          type="text"
          id="description"
          disabled={isUpdating}
          {...register("description", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Avatar image">
        <FileInput id="avatar" accept="image/*" disabled={isUpdating} />
      </FormRow>

      <FormRow>
        <Button
          type="reset"
          variation="secondary"
          onClick={() => onCloseModal?.()}
        >
          Cancel
        </Button>
        <Button disabled={isUpdating}>
          {isEditSession ? "Update account" : "Sign up a new users"}
        </Button>
      </FormRow>
    </Form>
  );
}

UpdateUserDataForm.propTypes = {
  userToEidt: PropTypes.shape({
    id: PropTypes.number.isRequired,
    email: PropTypes.string.isRequired,
    fullName: PropTypes.number.isRequired,
    phone: PropTypes.number.isRequired,
    department: PropTypes.number.isRequired,
    role: PropTypes.string.isRequired,
    description: PropTypes.string.isRequired,
  }),
  onCloseModal: PropTypes.func,
};

export default UpdateUserDataForm;
