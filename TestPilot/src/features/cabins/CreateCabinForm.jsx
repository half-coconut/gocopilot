import { useForm } from "react-hook-form";
import PropTypes from "prop-types";

import Input from "../../ui/Input";
import Form from "../../ui/Form.jsx";
import Button from "../../ui/Button.jsx";
import FileInput from "../../ui/FileInput.jsx";
import Textarea from "../../ui/Textarea.jsx";
import FormRow from "../../ui/FormRow.jsx";
import { useCreateCabiin } from "./useCreateCabin.js";
import { useEditCabiin } from "./useEditCabin.js";
import Switch from "../../ui/Switch.jsx";
import { useState } from "react";

function CreateCabinForm({ cabinToEdit = {}, onCloseModal }) {
  const { isCreating, createCabin } = useCreateCabiin();
  const { isEditing, editCabin } = useEditCabiin();
  const isWorking = isCreating || isEditing;

  const { id: editId, ...editValue } = cabinToEdit;
  const isEditSession = Boolean(editId);

  const { register, handleSubmit, reset, formState } = useForm({
    defaultValues: isEditSession ? editValue : {},
  });

  const { errors } = formState;

  const [isOn, setIsOn] = useState(false);
  const handleChange = () => {
    setIsOn(!isOn);
  };

  function onSubmit(data) {
    // const image = typeof data.image === "string" ? data.image : data.image[0];

    if (isEditSession)
      editCabin(
        { newCabinData: { ...data, id: editId, debug: isOn } },
        {
          onSuccess: () => {
            reset(), onCloseModal?.();
          },
        }
      );
    else
      createCabin(
        { ...data, debug: isOn },
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

  const data = `
  +++ Requests +++
  [total 总请求数: 1]
  [rate 请求速率: 2.50]
  [throughput 吞吐量: ...]
  +++ Duration +++
  [total 总持续时间: 399.29ms]
  ...
  +++ Success +++
  [ratio 成功率: 0.00%]
  [status codes:  400...:1]
  `;

  return (
    <Form
      onSubmit={handleSubmit(onSubmit, onError)}
      type={onCloseModal ? "modal" : "regular"}
    >
      <FormRow label="Interface name" error={errors?.name?.message}>
        <Input
          type="text"
          id="name"
          disabled={isWorking}
          {...register("name", { required: "This field is required" })}
        />
      </FormRow>

      <FormRow label="URL" error={errors?.name?.message}>
        <Input
          type="text"
          id="url"
          disabled={isWorking}
          {...register("url", { required: "This field is required" })}
        />
      </FormRow>

      <FormRow label="Project" error={errors?.maxCapacity?.message}>
        <Input
          type="text"
          id="project"
          disabled={isWorking}
          {...register("project", {
            required: "This field is required",
            min: {
              value: 1,
              message: "project should be at least one",
            },
          })}
        />
      </FormRow>

      <FormRow label="Type" error={errors?.regularPrice?.message}>
        <Input
          type="text"
          id="type"
          disabled={isWorking}
          {...register("type", {
            required: "This field is required",
            min: {
              value: 1,
              message: "type should be at least one",
            },
          })}
        />
      </FormRow>

      <FormRow label="Method" error={errors?.discount?.message}>
        <Input
          type="text"
          id="method"
          disabled={isWorking}
          {...register("method", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Params" error={errors?.discount?.message}>
        <Input
          type="text"
          id="params"
          disabled={isWorking}
          {...register("params", {
            required: "This field is required",
          })}
        />
      </FormRow>

      <FormRow label="Body" error={errors?.description?.message}>
        <Textarea
          type="text"
          id="body"
          disabled={isWorking}
          {...register("body", { required: "This field is required" })}
        />
      </FormRow>

      <FormRow label="Debug" error={errors?.discount?.message}>
        <Switch onChange={handleChange} />
      </FormRow>

      <FormRow label="Upload">
        <FileInput
          id="image"
          accept="image/*"
          type="file"
          // {...register("image", {
          //   required: isEditSession ? false : "This field is required",
          // })}
        />
      </FormRow>
      {/* 
      <FormRow label="Respone" error={errors?.discount?.message}>
        <Input disabled />
      </FormRow> */}
      <div>
        Respone
        <pre>{data}</pre>
      </div>

      <FormRow>
        {/* type is an HTML attribute! */}
        <Button
          variation="secondary"
          type="reset"
          onClick={() => onCloseModal?.()}
        >
          Cancel
        </Button>
        <Button disabled={isWorking}>
          {isEditSession ? "Edit interface" : "Create new interface"}
        </Button>
        {isOn ? <Button>Debug</Button> : ""}
      </FormRow>
    </Form>
  );
}

CreateCabinForm.propTypes = {
  cabinToEdit: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    maxCapacity: PropTypes.number.isRequired,
    regularPrice: PropTypes.number.isRequired,
    discount: PropTypes.number.isRequired,
    image: PropTypes.string.isRequired,
  }),
  onCloseModal: PropTypes.func,
};

export default CreateCabinForm;
