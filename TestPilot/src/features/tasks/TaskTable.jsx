import Table from "../../ui/Table.jsx";
import Menus from "../../ui/Menus.jsx";
import Empty from "../../ui/Empty.jsx";

import Spinner from "../../ui/Spinner.jsx";
import TaskRow from "./TaskRow.jsx";
import Pagination from "../../ui/Pagination.jsx";
import { useTasks } from "./useTasks.js";

function TaskTable() {
  const { isLoading, taskItems, total } = useTasks();

  if (isLoading) return <Spinner />;
  if (!total) return <Empty resourceName="tasks" />;

  return (
    <Menus>
      <Table columns="1fr 1fr 1fr 1.3fr 1fr 1fr 1fr 1fr">
        <Table.Header>
          {/* <div></div> */}
          <div>Name</div>
          <div>Interfaces</div>
          <div>Duration</div>
          <div>Workers Scope</div>
          <div>Rate</div>
          <div>Creator</div>
          <div>Ctime</div>
          <div></div>
        </Table.Header>

        <Table.Body
          data={taskItems}
          render={(taskItem) => (
            <TaskRow key={taskItem.id} taskItem={taskItem} />
          )}
        />
        <Table.Footer>
          <Pagination count={total} />
        </Table.Footer>
      </Table>
    </Menus>
  );
}

export default TaskTable;
