import { useForm, Controller } from "react-hook-form";
import Form from "../ui/Form.jsx";
import Heading from "../ui/Heading";
import Button from "../ui/Button.jsx";
import FormRow from "../ui/FormRow.jsx";
import Input from "../ui/Input.jsx";
import { Select } from "antd";

function History() {
  const resp = {
    code: 1,
    message: "OK",
    data: "\n+++ Requests +++\n[total 总请求数: 300]\n[rate 请求速率: 27.80]\n[throughput 吞吐量: 27.80]\n\n+++ Duration +++\n[total 总持续时间: 10.791s]\n\n+++ Latencies +++\n[min 最小响应时间: 274.508ms]\n[mean 平均响应时间: 35.97ms]\n[max 最大响应时间: 1.319s]\n[P50 百分之50 响应时间 (中位数): 283.015ms]\n[P90 百分之90 响应时间: 290.74ms]\n[P95 百分之95 响应时间: 952.224ms]\n[P99 百分之99 响应时间: 1.09s]\n\n+++ Success +++\n[ratio 成功率: 100.00%]\n[status codes:  200:300]\n[passed: 300]\n[failed: 0]\n",
  };

  const dataItem = resp?.data;
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

        <div>
          <pre>{dataItem}</pre>
        </div>
      </Form>
    </>
  );
}

export default History;
