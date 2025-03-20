import { useParams } from "react-router-dom";
import styled from "styled-components";

import Heading from "../../ui/Heading";
import DebugResultBox from "./DebugResultBox";
import { useMoveBack } from "../../hooks/useMoveBack";
import ButtonText from "../../ui/ButtonText";
import Row from "../../ui/Row";

const HeadingGroup = styled.div`
  display: flex;
  gap: 2.4rem;
  align-items: center;
`;

const data = `
+++ Requests +++
[total 总请求数: 1]
[rate 请求速率: 2.50]
[throughput 吞吐量: ...]
+++ Duration +++
[total 总持续时间: 399.29ms]
...
+++ Success +++
[ratio 成功率: 0.00%]
[status codes:  400...:1]
`;

function InterfaceDetail() {
  const { interfaceId } = useParams();
  const moveBack = useMoveBack();

  return (
    <>
      <Row type="horizontal">
        <HeadingGroup>
          <Heading as="h3">Respone #{interfaceId}</Heading>
        </HeadingGroup>
        <ButtonText onClick={moveBack}>&larr; Back</ButtonText>
      </Row>

      <DebugResultBox data={data} />
    </>
  );
}

export default InterfaceDetail;
