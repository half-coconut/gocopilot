import TaskTable from "../features/tasks/TaskTable";
// import TaskTableOperations from "../features/bookings/TaskTableOperations";

import AddTask from "../features/tasks/AddTask";
import Heading from "../ui/Heading";
import Row from "../ui/Row";

function Tasks() {
  return (
    <>
      <Row type="horizontal">
        <Heading as="h1">Performance testing tasks</Heading>
        {/* <TaskTableOperations /> */}
        <AddTask />
      </Row>

      <TaskTable />
    </>
  );
}

export default Tasks;
