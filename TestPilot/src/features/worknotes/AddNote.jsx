import styled from "styled-components";

import Button from "../../ui/Button";
import Modal from "../../ui/Modal";
import EditNotes from "./EditNotes";

const StyledButton = styled.div`
  &:has(button) {
    display: flex;
    justify-content: flex-end;
    /* align-items: flex-start; */
    gap: 1.2rem;
  }
`;

function AddNote() {
  return (
    <Modal>
      <Modal.Open opens="note-form">
        <StyledButton>
          <Button>Add new note</Button>
        </StyledButton>
      </Modal.Open>
      <Modal.Window name="note-form">
        <EditNotes />
      </Modal.Window>
    </Modal>
  );
}

export default AddNote;
