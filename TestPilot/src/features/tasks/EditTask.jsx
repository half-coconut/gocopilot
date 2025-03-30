import { useForm, Controller } from "react-hook-form";
import PropTypes from "prop-types";
import Form from "../../ui/Form.jsx";
import Heading from "../../ui/Heading";
import Button from "../../ui/Button.jsx";
import FormRow from "../../ui/FormRow.jsx";
import Input from "../../ui/Input.jsx";
import { Select } from "antd";

function EditTask({ onCloseModal }) {
  const initialData = {
    name: "task-01",
    a_ids: [],
    durations: "10m0s",
    workers: "5",
    max_workers: "200",
    timeout: "30s",
  };

  const {
    register,
    handleSubmit,
    control,
    // formState,
  } = useForm({
    defaultValues: initialData,
  });

  const interfaces = [
    { value: "1", label: "Login" },
    { value: "2", label: "Login2" },
    { value: "3", label: "Login3" },
  ];

  const handleChange = (value) => {
    console.log(`Selected: ${value}`);
  };

  function onSubmit(data) {
    console.log("onsubmit 的 data:", data);
  }

  return (
    <>
      <Heading as="h2">Create a new Task</Heading>
      <Form
        onSubmit={handleSubmit(onSubmit)}
        type={onCloseModal ? "modal" : "regular"}
      >
        <FormRow label="Task name">
          <Input
            type="text"
            id="name"
            {...register("name", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Interfaces">
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

        <FormRow label="Duration">
          <Input
            type="text"
            id="durations"
            {...register("durations", { required: "This field is required" })}
          />
        </FormRow>
        <FormRow label="Workers">
          <Input
            type="text"
            id="workers"
            {...register("workers", { required: "This field is required" })}
          />
        </FormRow>
        <FormRow label="Max workers">
          <Input
            type="text"
            id="max_workers"
            {...register("max_workers", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Timeout">
          <Input
            type="text"
            id="timeout"
            {...register("timeout", { required: "This field is required" })}
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
