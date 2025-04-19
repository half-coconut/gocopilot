import styled from "styled-components";

import Button from "../../ui/Button";
import Modal from "../../ui/Modal";
import EditJob from "./EditJob";

const StyledButton = styled.div`
  &:has(button) {
    display: flex;
    justify-content: flex-end;
    /* align-items: flex-start; */
    gap: 1.2rem;
  }
`;

function AddJob() {
  return (
    <Modal>
      <Modal.Open opens="task-form">
        <StyledButton>
          <Button>Add new job</Button>
        </StyledButton>
      </Modal.Open>
      <Modal.Window name="task-form">
        <EditJob />
      </Modal.Window>
    </Modal>
  );
}

export default AddJob;
