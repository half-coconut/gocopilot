import { useParams } from "react-router-dom";
import styled from "styled-components";
import Tooltip from "@mui/material/Tooltip";
import IconButton from "@mui/material/IconButton";
import HelpIcon from "@mui/icons-material/Help";

import Heading from "../../ui/Heading";
import { useMoveBack } from "../../hooks/useMoveBack";
import ButtonText from "../../ui/ButtonText";
import Row from "../../ui/Row";
// import { HiArrowUpOnSquare, HiOutlineCheckCircle } from "react-icons/hi2";
import { HiOutlineCheckCircle } from "react-icons/hi2";
import { useTask } from "./useTasks";
import Spinner from "../../ui/Spinner";
import Tag from "../../ui/Tag";
import TaskDataBox from "./TaskDataBox";

import DataItem from "../../ui/DataItem.jsx";
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
  width: 200px;
`;

const HeadingGroup = styled.div`
  display: flex;
  gap: 2.4rem;
  align-items: center;
`;

function TaskDetail() {
  const { taskId } = useParams();
  const moveBack = useMoveBack();
  const [status, setStatus] = useState("api");

  const { isLoading, taskItem } = useTask();

  if (isLoading) return <Spinner />;

  const statusList = [
    { value: "api", label: "Debugging Interfaces" },
    { value: "task", label: "Debugging Task" },
    { value: "execute", label: "Execute" },
  ];
  const statusToTagName = {
    api: "blue",
    task: "yellow",
    execute: "green",
    // failed: "silver",
  };

  return (
    <>
      <Row type="horizontal">
        <HeadingGroup>
          <Heading as="h1">Debug history #{taskId}</Heading>
          <Tag type={statusToTagName[status]}>{status.replace("-", " ")}</Tag>
        </HeadingGroup>
        <ButtonText onClick={moveBack}>&larr; Back</ButtonText>
      </Row>
      <DataItem icon={<HiOutlineCheckCircle />} label="status">
        <StyledSelect
          value={status}
          onChange={(e) => setStatus(e.target.value)}
        >
          {statusList.map((s) => (
            <option key={s.value} value={s.value}>
              {s.label}
            </option>
          ))}
        </StyledSelect>
        <Tooltip
          title={
            <span style={{ fontSize: "1.5rem" }}>
              Select different status to debug or run
            </span>
          }
        >
          <IconButton>
            <HelpIcon style={{ fontSize: "2rem" }} />
          </IconButton>
        </Tooltip>
      </DataItem>

      <TaskDataBox taskItem={taskItem?.data} status={status} />
    </>
  );
}

export default TaskDetail;
