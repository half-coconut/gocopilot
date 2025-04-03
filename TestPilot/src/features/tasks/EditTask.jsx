import { useForm, Controller } from "react-hook-form";
import PropTypes from "prop-types";
import Form from "../../ui/Form.jsx";
import Heading from "../../ui/Heading";
import Button from "../../ui/Button.jsx";
import FormRow from "../../ui/FormRow.jsx";
import Input from "../../ui/Input.jsx";
import { Select } from "antd";
import { useCreateTask } from "./useCreateTask.js";
import { useInterfaces } from "../interfaces/useInterfaces.js";

function EditTask({ onCloseModal }) {
  const { isLoading, interfaceItems } = useInterfaces();

  const { isCreating, createTask } = useCreateTask();
  const isWorking = isCreating || isLoading;

  const interfaces = !isLoading
    ? interfaceItems.map((item) => ({
        value: item.id,
        label: item.name,
      }))
    : [];

  const initialData = {
    name: "task-01",
    a_ids: [],
    durations: "10m0s",
    workers: 5,
    max_workers: 100,
    rate: 10,
  };

  const { register, handleSubmit, control, reset, formState } = useForm({
    defaultValues: initialData,
  });

  const { errors } = formState;

  const handleChange = (value) => {
    console.log(`Selected: ${value}`);
  };

  function onSubmit(data) {
    console.log("onsubmit 的 data:", data);
    createTask(
      { ...data },
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

  return (
    <>
      <Heading as="h2">Create a new Task</Heading>
      <Form
        onSubmit={handleSubmit(onSubmit, onError)}
        type={onCloseModal ? "modal" : "regular"}
      >
        <FormRow label="Task name" error={errors?.name?.message}>
          <Input
            type="text"
            id="name"
            {...register("name", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Interfaces" error={errors?.a_ids?.message}>
          <Controller
            id="a_ids"
            name="a_ids"
            {...register("a_ids", {
              required: "This field is required",
            })}
            control={control}
            render={({ field }) => (
              <Select
                {...field}
                showSearch
                mode="tags"
                disabled={isWorking}
                style={{ width: 300 }}
                placeholder="Tags Mode"
                onChange={(value) => {
                  field.onChange(value); // 更新 Controller 的值
                  handleChange(value); // 进行其他处理
                }}
                options={interfaces}
              />
            )}
          />
        </FormRow>

        <FormRow label="Duration" error={errors?.durations?.message}>
          <Input
            type="text"
            id="durations"
            {...register("durations", { required: "This field is required" })}
          />
        </FormRow>
        <FormRow label="Workers" error={errors?.workers?.message}>
          <Input
            type="number"
            id="workers"
            {...register("workers", {
              required: "This field is required",
              valueAsNumber: true,
              min: {
                value: 1,
                message: "Min value is 1",
              },
            })}
          />
        </FormRow>
        <FormRow label="Max workers" error={errors?.max_workers?.message}>
          <Input
            type="number"
            id="max_workers"
            {...register("max_workers", {
              required: "This field is required",
              valueAsNumber: true,
              min: {
                value: 1,
                message: "Min value is 1",
              },
            })}
          />
        </FormRow>

        <FormRow label="Rate" error={errors?.rate?.message}>
          <Input
            type="number"
            id="rate"
            disabled={isWorking}
            {...register("rate", {
              required: "This field is required",
              valueAsNumber: true,
              min: {
                value: 1,
                message: "Min value is 1",
              },
            })}
          />
        </FormRow>

        <FormRow>
          <Button
            variation="secondary"
            type="reset"
            onClick={() => onCloseModal?.()}
          >
            Cancel
          </Button>
          <Button type="submit">submit</Button>
        </FormRow>
      </Form>
    </>
  );
}

EditTask.propTypes = {
  onCloseModal: PropTypes.func,
};

export default EditTask;
