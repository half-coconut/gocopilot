import { useForm, Controller } from "react-hook-form";
import Form from "../ui/Form.jsx";
import Heading from "../ui/Heading";
import Button from "../ui/Button.jsx";
import FormRow from "../ui/FormRow.jsx";
import Input from "../ui/Input.jsx";
import { Select } from "antd";

function History() {
  const initialData = {
    name: "Test",
    project: "gaia",
    header: '{"Content-Type": "text/plain"}',
  };

  const { register, handleSubmit, getValues, control } = useForm({
    defaultValues: initialData,
  });
  const projects = [
    { value: "bsc", label: "BSC" },
    { value: "core", label: "CORE" },
    { value: "eth", label: "ETH" },
    { value: "aether", label: "AETHER" },
    { value: "carv", label: "CARV" },
    { value: "gaia", label: "GAIA" },
  ];

  const headers = [
    {
      value: '{"Content-Type": "application/json"}',
      label: "Content-Type application/json",
    },
    {
      value: '{"Content-Type": "application/x-www-form-urlencoded"}',
      label: "Content-Type application/x-www-form-urlencoded",
    },
    {
      value: '{"Content-Type": "multipart/form-data"}',
      label: "Content-Type multipart/form-data",
    },
    {
      value: '{"Content-Type": "text/plain"}',
      label: "Content-Type text/plain",
    },
  ];

  function onSubmit(data) {
    console.log("onsubmit 的 data:", data);
    console.log("getValues 的 data:", getValues());
  }

  return (
    <>
      <Heading as="h1">History</Heading>
      <Form onSubmit={handleSubmit(onSubmit)}>
        <FormRow label="Interface name">
          <Input
            type="text"
            id="name"
            {...register("name", { required: "This field is required" })}
          />
        </FormRow>

        <FormRow label="Project">
          <Controller
            id="project"
            name="project"
            {...register("project", { required: "This field is required" })}
            control={control}
            render={({ field }) => (
              <Select
                {...field}
                showSearch
                style={
                  {
                    //   width: 200,
                  }
                }
                placeholder="Search to Select"
                optionFilterProp="label"
                filterSort={(optionA, optionB) =>
                  (optionA?.label ?? "")
                    .toLowerCase()
                    .localeCompare((optionB?.label ?? "").toLowerCase())
                }
                options={projects}
              />
            )}
          />
        </FormRow>

        <FormRow label="Header">
          <Controller
            id="header"
            name="header"
            {...register("header", { required: "This field is required" })}
            control={control}
            render={({ field }) => (
              <Select
                {...field}
                showSearch
                style={{
                  width: 300,
                }}
                placeholder="Search to Select"
                optionFilterProp="label"
                filterSort={(optionA, optionB) =>
                  (optionA?.label ?? "")
                    .toLowerCase()
                    .localeCompare((optionB?.label ?? "").toLowerCase())
                }
                options={headers}
              />
            )}
          />
        </FormRow>

        <FormRow>
          <Button>submit</Button>
        </FormRow>
      </Form>
    </>
  );
}

export default History;
