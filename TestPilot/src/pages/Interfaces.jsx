import AddInterface from "../features/interfaces/AddInterface";
import InterfaceTable from "../features/interfaces/InterfaceTable";
import InterfaceTableOperations from "../features/interfaces/InterfaceTableOperations";
import Heading from "../ui/Heading";
import Row from "../ui/Row";

function Cabins() {
  return (
    <>
      <Row type="horizontal">
        <Heading as="h1">All Interfaces</Heading>
        <InterfaceTableOperations />
      </Row>

      <Row>
        <InterfaceTable />
        <AddInterface />
      </Row>
    </>
  );
}

export default Cabins;
