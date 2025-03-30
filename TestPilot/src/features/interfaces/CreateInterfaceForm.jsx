import { useForm } from "react-hook-form";
import styled from "styled-components";
import PropTypes from "prop-types";

import Form from "../../ui/Form.jsx";
import Button from "../../ui/Button.jsx";
import FileInput from "../../ui/FileInput.jsx";
import Textarea from "../../ui/Textarea.jsx";
import FormRow from "../../ui/FormRow.jsx";
import { useCreateInterface } from "./useCreateInterface.js";
import { useEditInterface } from "./useEditInterface.js";
import Switch from "../../ui/Switch.jsx";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const Input = styled.input`
  border: 1px solid var(--color-grey-300);
  background-color: var(--color-grey-0);
  border-radius: var(--border-radius-sm);
  box-shadow: var(--shadow-sm);
  padding: 0.8rem 1.2rem;
  width: 300px;
`;

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
  width: 300px;
`;

function CreateInterfaceForm({ interfaceToEdit = {}, onCloseModal }) {
  const { isCreating, createInterface } = useCreateInterface();
  const { isEditing, editInterface } = useEditInterface();
  const isWorking = isCreating || isEditing;

  const { id: editId, ...editValue } = interfaceToEdit;
  const isEditSession = Boolean(editId);

  const { register, handleSubmit, reset, formState } = useForm({
    defaultValues: isEditSession ? editValue : {},
  });

  const { errors } = formState;

  const navigate = useNavigate();

  console.log("编辑状态获取的值：", editValue);

  const [isOn, setIsOn] = useState(false);
  const handleChange = () => {
    setIsOn(!isOn);
  };

  const projectsList = [
    { value: "BSC", label: "BSC" },
    { value: "CORE", label: "CORE" },
    { value: "ETH", label: "ETH" },
    { value: "AETHER", label: "AETHER" },
    { value: "CARV", label: "CARV" },
    { value: "GAIA", label: "GAIA" },
  ];

  // const headers = [
  //   {
  //     value: '{"Content-Type": "application/json"}',
  //     label: "Content-Type application/json",
  //   },
  //   {
  //     value: '{"User-Agent", "PostmanRuntime/7.39.0"}',
  //     label: "User-Agent, PostmanRuntime/7.39.0",
  //   },
  //   {
  //     value: '{"Content-Type": "application/x-www-form-urlencoded"}',
  //     label: "Content-Type application/x-www-form-urlencoded",
  //   },
  //   {
  //     value: '{"Content-Type": "multipart/form-data"}',
  //     label: "Content-Type multipart/form-data",
  //   },
  //   {
  //     value: '{"Content-Type": "text/plain"}',
  //     label: "Content-Type text/plain",
  //   },
  // ];
  const types = [
    { value: "http", label: "http" },
    { value: "websocket", label: "websocket" },
  ];

  const methods = [
    { value: "GET", label: "GET" },
    { value: "POST", label: "POST" },
    { value: "PUT", label: "PUT" },
    { value: "PATCH", label: "PATCH" },
    { value: "DELETE", label: "DELETE" },
    { value: "HEAD", label: "HEAD" },
    { value: "OPTIONS", label: "OPTIONS" },
  ];
  function onSubmit(data) {
    // const image = typeof data.image === "string" ? data.image : data.image[0];
    console.log("onsubmit 的 data:", data);

    if (isEditSession)
      editInterface(
        {
          newInterfaceData: {
            ...data,
            id: editId,
            debug: isOn,
          },
        },
        {
          onSuccess: () => {
            {
              isOn ? navigate(`/interfaces/${editId}`) : reset(),
                onCloseModal?.();
            }
          },
        }
      );
    else
      createInterface(
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

      <FormRow label="Project" error={errors?.project?.message}>
        <StyledSelect {...register("project")} id="project">
          {projectsList.map((project) => (
            <option key={project.value} value={project.value}>
              {project.label}
            </option>
          ))}
        </StyledSelect>
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

      <FormRow label="Method" error={errors?.method?.message}>
        <StyledSelect
          {...register("method", {
            required: "This field is required",
          })}
          id="method"
        >
          {methods.map((method) => (
            <option key={method.value} value={method.value}>
              {method.label}
            </option>
          ))}
        </StyledSelect>
      </FormRow>

      <FormRow label="URL" error={errors?.url?.message}>
        <Input
          type="text"
          id="url"
          disabled={isWorking}
          {...register("url", { required: "This field is required" })}
        />
      </FormRow>

      <FormRow
        label="Authorization(Under construction)"
        error={errors?.params?.message}
      >
        <Input type="text" id="authorization" disabled />
      </FormRow>

      <FormRow label="Headers" error={errors?.header?.message}>
        <Textarea
          type="text"
          id="header"
          disabled={isWorking}
          {...register("header", { required: "This field is required" })}
        />
      </FormRow>

      <FormRow label="Body" error={errors?.body?.message}>
        <Textarea
          type="text"
          id="body"
          disabled={isWorking}
          {...register("body", { required: "This field is required" })}
        />
      </FormRow>
      <FormRow label="Params" error={errors?.params?.message}>
        <Input
          type="text"
          id="params"
          disabled={isWorking}
          {...register("params", {
            required: "This field is required",
          })}
        />
      </FormRow>

      {isEditSession ? (
        <FormRow label="Debug" error={errors?.discount?.message}>
          <Switch onChange={handleChange} />
        </FormRow>
      ) : (
        ""
      )}

      <FormRow label="Upload(Under construction)">
        <FileInput
          id="image"
          accept="image/*"
          type="file"
          // {...register("image", {
          //   required: isEditSession ? false : "This field is required",
          // })}
          disabled
        />
      </FormRow>

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
        {isOn ? <Button variation="run">Debug</Button> : ""}
      </FormRow>
    </Form>
  );
}

CreateInterfaceForm.propTypes = {
  interfaceToEdit: PropTypes.shape({
    id: PropTypes.number,
    name: PropTypes.string,
    url: PropTypes.string,
    project: PropTypes.string,
    header: PropTypes.string,
    type: PropTypes.string,
    method: PropTypes.string,
    params: PropTypes.string,
    body: PropTypes.string,
  }),
  onCloseModal: PropTypes.func,
};

export default CreateInterfaceForm;
