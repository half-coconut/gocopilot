// import TaskTable from "../features/bookings/TaskTable";
// import TaskTableOperations from "../features/bookings/TaskTableOperations";

import AddTask from "../features/tasks/AddTask";
import Heading from "../ui/Heading";
import Row from "../ui/Row";

function Tasks() {
  return (
    <>
      <Row type="horizontal">
        <Heading as="h1">All Tasks</Heading>
        {/* <TaskTableOperations /> */}
        <AddTask />
      </Row>

      {/* <TaskTable /> */}
    </>
  );
}

export default Tasks;
