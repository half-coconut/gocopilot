import styled from "styled-components";
import PropTypes from "prop-types";
import { JsonEditor, githubLightTheme } from "json-edit-react";
// import { JsonEditor, monoLightTheme } from "json-edit-react";

import {
  HiOutlineClock,
  HiOutlineUserPlus,
  // HiOutlineChatBubbleBottomCenterText,
  // HiOutlineCheckCircle,
  HiOutlineRocketLaunch,
  HiOutlineClipboardDocumentList,
  HiOutlineArchiveBox,
  HiOutlineUser,
} from "react-icons/hi2";

import { TbReportAnalytics } from "react-icons/tb";

import {
  MdOutlineFaceUnlock,
  MdOutlineFaceRetouchingNatural,
} from "react-icons/md";

import DataItem from "../../ui/DataItem.jsx";
import { convertNanosecondsToHMS } from "../../utils/helpers.js";
import {
  useDebugInterfaces,
  useDebugTask,
  useExecuteTask,
} from "./useExecuteTask.js";

const StyledBookingDataBox = styled.section`
  /* Box */
  background-color: var(--color-grey-0);
  border: 1px solid var(--color-grey-100);
  border-radius: var(--border-radius-md);

  overflow: hidden;
`;

const Header = styled.header`
  background-color: var(--color-brand-500);
  padding: 2rem 4rem;
  color: #e0e7ff;
  font-size: 1.8rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: space-between;

  svg {
    height: 3.2rem;
    width: 3.2rem;
  }

  & div:first-child {
    display: flex;
    align-items: center;
    gap: 1.6rem;
    font-weight: 600;
    font-size: 1.8rem;
  }

  & span {
    font-family: "Sono";
    font-size: 2rem;
    margin-left: 4px;
  }
`;

const Section = styled.section`
  padding: 3.2rem 4rem 1.2rem;
`;

const Guest = styled.div`
  display: flex;
  align-items: center;
  gap: 1.2rem;
  margin-bottom: 1.6rem;
  color: var(--color-grey-500);

  & p:first-of-type {
    font-weight: 500;
    color: var(--color-grey-700);
  }
`;

const Result = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.6rem 3.2rem;
  border-radius: var(--border-radius-sm);
  margin-top: 2.4rem;

  background-color: ${(props) =>
    props.isPassed ? "var(--color-green-100)" : "var(--color-yellow-100)"};
  color: ${(props) =>
    props.isPassed ? "var(--color-green-700)" : "var(--color-yellow-700)"};

  & p:last-child {
    text-transform: uppercase;
    font-size: 1.4rem;
    font-weight: 600;
  }

  svg {
    height: 2.4rem;
    width: 2.4rem;
    color: currentColor !important;
  }
`;

const StyledNote = styled.div`
  /* Box */
  background-color: var(--color-grey-0);
  border: 1px solid var(--color-grey-100);
  border-radius: var(--border-radius-md);

  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.6rem 3.2rem;
  border-radius: var(--border-radius-sm);
  /* margin-top: 2.4rem; */
  /* height: 300px; */
  overflow: auto;

  flex-direction: column;
  gap: 2.4rem;
  grid-column: 1 / span 2;
  /* padding-top: 2.4rem; */
`;

const StyleUl = styled.ul`
  margin-top: 10px;
  color: var(--color-grey-500);
  gap: 2.4rem;
`;

const Footer = styled.footer`
  padding: 1.6rem 4rem;
  font-size: 1.2rem;
  color: var(--color-grey-500);
  text-align: right;
`;

function TaskDataBox({ taskItem, status }) {
  const { isLoading, debugInterfacesItem } = useDebugInterfaces();
  const { isLoading: isDebugTaskLoading, debugTaskItem } = useDebugTask();
  const { isLoading: isExecuteLoading, executeTaskItem } = useExecuteTask();

  // if (!isLoading && status == "api") {
  //   console.log(
  //     "box 页面显示的内容 debuggingApis",
  //     JSON.stringify(debugInterfacesItem, null, 2)
  //   );
  // }

  // if (!isDebugTaskLoading && status == "task") {
  //   console.log(
  //     "box 页面显示的内容 debuggingTask",
  //     JSON.stringify(debugTaskItem, null, 2)
  //   );
  // }

  // if (!isExecuteLoading && status == "execute") {
  //   console.log(
  //     "box 页面显示的内容 executeTaskItem",
  //     JSON.stringify(executeTaskItem, null, 2)
  //   );
  // }

  const {
    id,
    name,
    a_ids,
    // apis,
    durations,
    workers,
    max_workers,
    rate,
    creator,
    updater,
    ctime,
    utime,
  } = taskItem;

  return (
    <StyledBookingDataBox>
      <Header>
        <div>
          <HiOutlineClipboardDocumentList />
          <p>{name}</p>
        </div>
      </Header>

      <Section>
        <Guest>
          <HiOutlineArchiveBox />
          <p>{id}</p>
          <span>&bull;</span>
          <p>{name}</p>
          <HiOutlineUser />
          <p>creator</p>
          <span>&bull;</span>
          <p>{creator}</p>
        </Guest>

        <DataItem icon={<HiOutlineClock />} label="Durations">
          <span>&bull;</span>
          <p>{convertNanosecondsToHMS(durations)}</p>
        </DataItem>

        <DataItem icon={<HiOutlineUserPlus />} label="Workers Scope">
          <span>&bull;</span>
          workers:
          <p>{workers}</p>
          <span>&bull;</span>
          max workers:
          <p>{max_workers}</p>
        </DataItem>
        <DataItem icon={<HiOutlineRocketLaunch />} label="Rate">
          <span>&bull;</span>
          <p>{rate}</p>
          <span>&bull;</span>
          Number of interfaces:
          <p>{a_ids.length}</p>
        </DataItem>

        {/* <DataItem icon={<HiOutlineCheckCircle />} label="apis">
          <pre>
            {apis.map((item, index) => (
              <p key={index}>{item}</p>
            ))}
          </pre>
        </DataItem> */}

        <Result isPassed={status}>
          <DataItem
            icon={<TbReportAnalytics />}
            label={
              status == "api"
                ? `Interfaces Debug Result`
                : status == "task"
                ? `Task Debug Result`
                : `Task Execute Result`
            }
          ></DataItem>
        </Result>

        <StyledNote>
          <StyleUl>
            {!isLoading && status == "api" ? (
              <JsonEditor
                data={JSON.stringify(debugInterfacesItem, null, 2)}
                theme={githubLightTheme}
              />
            ) : (
              ""
            )}
            {!isDebugTaskLoading && status == "task" ? (
              <pre>{debugTaskItem?.data}</pre>
            ) : (
              ""
            )}
            {!isExecuteLoading && status == "execute" ? (
              <pre>{executeTaskItem?.data}</pre>
            ) : (
              ""
            )}
          </StyleUl>
        </StyledNote>
      </Section>

      <Footer>
        <p>
          <MdOutlineFaceUnlock /> {creator} Created on {ctime}
        </p>
        <p>
          <MdOutlineFaceRetouchingNatural /> {updater} Updated on {utime}
        </p>
      </Footer>
    </StyledBookingDataBox>
  );
}

TaskDataBox.propTypes = {
  taskItem: PropTypes.shape({
    id: PropTypes.number,
    name: PropTypes.string,
    a_ids: PropTypes.List,
    apis: PropTypes.List,
    durations: PropTypes.string,
    workers: PropTypes.number,
    max_workers: PropTypes.number,
    rate: PropTypes.string,
    creator: PropTypes.string,
    updater: PropTypes.string,
    ctime: PropTypes.string,
    utime: PropTypes.string,
  }).isRequired,
  status: PropTypes.string,
};

export default TaskDataBox;
