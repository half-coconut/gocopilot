import Heading from "../ui/Heading";
import Row from "../ui/Row";

function Jobs() {
  return (
    <>
      <Row type="horizontal">
        <Heading as="h1">All Jobs</Heading>
      </Row>

      <Row>
        Here is your jobs list...
        {/* <InterfaceTable />
        <AddInterface /> */}
      </Row>
    </>
  );
}

export default Jobs;
