import { useParams } from "react-router-dom";
import Heading from "../../ui/Heading";
import DebugResultBox from "./DebugResultBox";

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
  return (
    <>
      <Heading as="h3">Respone #{interfaceId}</Heading>
      <DebugResultBox data={data} />
    </>
  );
}

export default InterfaceDetail;
