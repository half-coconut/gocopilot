import Button from "../../ui/Button";
import Modal from "../../ui/Modal";

import EditNotes from "./EditNotes";

function AddNote() {
  return (
    <Modal>
      <Modal.Open opens="note-form">
        <Button>Add new note</Button>
      </Modal.Open>
      <Modal.Window name="note-form">
        <EditNotes />
      </Modal.Window>
    </Modal>
  );
}

export default AddNote;
