import { useForm } from "react-hook-form";
import PropTypes from "prop-types";
import styled from "styled-components";
import HelpIcon from "@mui/icons-material/Help";
import { Tooltip } from "antd";
import { IconButton } from "@mui/material";

import Form from "../../ui/Form.jsx";
import Heading from "../../ui/Heading.jsx";
import Button from "../../ui/Button.jsx";
import FormRow from "../../ui/FormRow.jsx";
import Input from "../../ui/Input.jsx";
import { useCreateJob } from "./useCreateJob.js";
import { useTasks } from "../tasks/useTasks.js";

import Textarea from "../../ui/Textarea.jsx";
import { useState } from "react";

const StyledSelect = styled.select`
  font-size: 1.4rem;
  padding: 0.8rem 1.2rem;
  border: 1px solid
    ${(props) =>
      props.type === "white"
        ? "var(--color-grey-100)"
        : "var(--color-grey-300)"};
  border-radius: var(--border-radius-sm);
  background-color: var(--color-grey-0);
  font-weight: 500;
  box-shadow: var(--shadow-sm);
  width: 230px;
`;

function EditJob({ onCloseModal }) {
  const { isLoading, taskItems } = useTasks();

  const { isCreating, createJob } = useCreateJob();

  const isWorking = isCreating || isLoading;

  const tasks = !isLoading
    ? taskItems.map((item) => ({
        value: Number(item.id),
        label: item.name,
      }))
    : [];

  const [selectedTaskId, setSelectedTaskId] = useState(1);

  const handleTaskChange = (event) => {
    setSelectedTaskId(Number(event.target.value)); // Convert to number here
  };

  const types = [
    { value: "scheduled", label: "Scheduled" },
    { value: "http", label: "HTTP" },
  ];

  const initialData = {
    name: "Internal Task",
    description: "This is a sample task description.",
    type: "scheduled",
    cron: "*/15 * * * *",
    http_cfg: "",
    task_id: 1,
    timezone: "UTC",
    durations: "1m0s",
    retry: false,
    max_retries: 0,
  };

  const { register, handleSubmit, reset, formState } = useForm({
    defaultValues: initialData,
  });

  const { errors } = formState;

  function onSubmit(data) {
    createJob(
      { ...data, task_id: selectedTaskId },
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
      <Heading as="h2">Create a new job</Heading>
      <Form
        onSubmit={handleSubmit(onSubmit, onError)}
        type={onCloseModal ? "modal" : "regular"}
      >
        <FormRow label="Job name" error={errors?.name?.message}>
          <Input
            type="text"
            id="name"
            {...register("name", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Description" error={errors?.description?.message}>
          <Textarea
            type="text"
            id="description"
            disabled={isWorking}
            {...register("description", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Type" error={errors?.type?.message}>
          <StyledSelect {...register("type")} id="type">
            {types.map((type) => (
              <option key={type.value} value={type.value}>
                {type.label}
              </option>
            ))}
          </StyledSelect>
        </FormRow>

        <FormRow label="Cron" error={errors?.cron?.message}>
          <Input
            type="text"
            id="cron"
            {...register("cron", { required: "This field is required" })}
          />
        </FormRow>

        <Tooltip
          title={
            <span style={{ fontSize: "1.5rem" }}>
              Please enter a cron expression
            </span>
          }
        >
          <IconButton>
            <HelpIcon style={{ fontSize: "2rem" }} />
          </IconButton>
        </Tooltip>

        <FormRow label="HTTP Config" error={errors?.http_cfg?.message}>
          <Input type="text" id="http_cfg" {...register("http_cfg")} />
        </FormRow>

        <FormRow label="Task" error={errors?.task_id?.message}>
          <StyledSelect
            {...register("task_id")}
            id="task_id"
            onChange={handleTaskChange}
          >
            {tasks.map((task) => (
              <option key={task.value} value={task.value}>
                {task.label}
              </option>
            ))}
          </StyledSelect>
        </FormRow>

        <FormRow label="Duration" error={errors?.durations?.message}>
          <Input
            type="text"
            id="durations"
            {...register("durations", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Timezone" error={errors?.timezone?.message}>
          <Input
            type="text"
            id="timezone"
            {...register("timezone", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Retry" error={errors?.retry?.message}>
          <Input
            type="checkbox"
            id="retry"
            {...register("retry", {
              required: "This field is required",
            })}
          />
        </FormRow>
        <FormRow label="Max Retries" error={errors?.max_retries?.message}>
          <Input
            type="number"
            id="max_retries"
            {...register("max_retries", {
              required: "This field is required",
              valueAsNumber: true,
              min: {
                value: 0,
                message: "Min value is 0",
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

EditJob.propTypes = {
  onCloseModal: PropTypes.func,
};

export default EditJob;
