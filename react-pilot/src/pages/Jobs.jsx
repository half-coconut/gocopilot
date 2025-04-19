import Heading from "../ui/Heading";
import Row from "../ui/Row";
import AddJob from "../features/jobs/AddJob";

function Jobs() {
  return (
    <>
      <Row type="horizontal">
        <Heading as="h1">All Jobs</Heading>
      </Row>

      <Row>
        {/* <InterfaceTable />*/}
        <AddJob />
      </Row>
    </>
  );
}

export default Jobs;
