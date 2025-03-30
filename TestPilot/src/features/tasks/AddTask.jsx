import styled from "styled-components";

import Button from "../../ui/Button";
import Modal from "../../ui/Modal";
import EditTask from "./EditTask";

const StyledButton = styled.div`
  &:has(button) {
    display: flex;
    justify-content: flex-end;
    /* align-items: flex-start; */
    gap: 1.2rem;
  }
`;

function AddTask() {
  return (
    <Modal>
      <Modal.Open opens="task-form">
        <StyledButton>
          <Button>Add new task</Button>
        </StyledButton>
      </Modal.Open>
      <Modal.Window name="task-form">
        <EditTask />
      </Modal.Window>
    </Modal>
  );
}

export default AddTask;
