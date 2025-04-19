import { Select } from "antd";

const option = [
  {
    value: "1",
    label: "Not Identified",
  },
  {
    value: "2",
    label: "Closed",
  },
  {
    value: "3",
    label: "Communicated",
  },
  {
    value: "4",
    label: "Identified",
  },
  {
    value: "5",
    label: "Resolved",
  },
  {
    value: "6",
    label: "Cancelled",
  },
];

function Environments() {
  return (
    <Select
      showSearch
      style={{ width: 200 }}
      placeholder="Search to Select"
      optionFilterProp="label"
      filterSort={(optionA, optionB) =>
        (optionA?.label ?? "")
          .toLowerCase()
          .localeCompare((optionB?.label ?? "").toLowerCase())
      }
      options={option}
    />
  );
}

export default Environments;
