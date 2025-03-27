import { useParams } from "react-router-dom";
import styled from "styled-components";

import Heading from "../../ui/Heading";
import DebugResultBox from "./DebugResultBox";
import { useMoveBack } from "../../hooks/useMoveBack";
import ButtonText from "../../ui/ButtonText";
import Row from "../../ui/Row";
import Spinner from "../../ui/Spinner";
import { useInterface } from "./useInterfaces";

const HeadingGroup = styled.div`
  display: flex;
  gap: 2.4rem;
  align-items: center;
`;

function InterfaceDetail() {
  const { interfaceId } = useParams();
  const moveBack = useMoveBack();

  const { isLoading, interfaceItem } = useInterface();

  if (isLoading) return <Spinner />;

  return (
    <>
      <Row type="horizontal">
        <HeadingGroup>
          <Heading as="h3">Respone #{interfaceId}</Heading>
        </HeadingGroup>
        <ButtonText onClick={moveBack}>&larr; Back</ButtonText>
      </Row>

      <DebugResultBox interfaceItem={interfaceItem} />
    </>
  );
}

export default InterfaceDetail;
