import styled from "styled-components";
import PropTypes from "prop-types";
// import { JsonEditor, monoLightTheme } from "json-edit-react";
import { JsonEditor, githubLightTheme, monoLightTheme } from "json-edit-react";

import {
  HiOutlineChatBubbleBottomCenterText,
  HiOutlineCheckCircle,
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
    props.isPaid ? "var(--color-green-100)" : "var(--color-yellow-100)"};
  color: ${(props) =>
    props.isPaid ? "var(--color-green-700)" : "var(--color-yellow-700)"};

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
  height: 800px;
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

function DebugResultBox({ interfaceItem }) {
  console.log("interface item: ", interfaceItem);

  const {
    id,
    name,
    project,
    url,
    params,
    method,
    header,
    body,
    type,
    debug_result,
    creator,
    updater,
    ctime,
    utime,
  } = interfaceItem;

  console.log(
    id,
    name,
    project,
    url,
    params,
    method,
    header,
    body,
    type,
    debug_result,
    creator,
    updater,
    ctime,
    utime
  );

  return (
    <StyledBookingDataBox>
      <Header>
        <div>
          <HiOutlineClipboardDocumentList />
          <p>{project}</p>
        </div>

        <p></p>
      </Header>

      <Section>
        <Guest>
          <HiOutlineArchiveBox />
          <p>{name}</p>
          <span>&bull;</span>
          <HiOutlineUser />
          <p>creator</p>
          <span>&bull;</span>
          <p>{creator}</p>
        </Guest>

        <DataItem icon={<HiOutlineChatBubbleBottomCenterText />} label={method}>
          {url}
        </DataItem>

        <DataItem icon={<HiOutlineCheckCircle />} label="headers">
          <pre>
            <JsonEditor data={JSON.parse(header)} theme={monoLightTheme} />
          </pre>
        </DataItem>

        <Result isPaid="true">
          <DataItem
            icon={<TbReportAnalytics />}
            label={`Total Result`}
          ></DataItem>
        </Result>

        <StyledNote>
          <StyleUl>
            <JsonEditor
              data={JSON.parse(debug_result)}
              theme={githubLightTheme}
            />
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

DebugResultBox.propTypes = {
  interfaceItem: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    project: PropTypes.string.isRequired,
    url: PropTypes.string.isRequired,
    params: PropTypes.string.isRequired,
    method: PropTypes.string.isRequired,
    body: PropTypes.string.isRequired,
    header: PropTypes.string.isRequired,
    type: PropTypes.string.isRequired,
    debug_result: PropTypes.string.isRequired,
    creator: PropTypes.number.isRequired,
    updater: PropTypes.number.isRequired,
    ctime: PropTypes.string.isRequired,
    utime: PropTypes.string.isRequired,
  }).isRequired,
};

export default DebugResultBox;
