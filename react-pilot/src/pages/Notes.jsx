import Heading from "../ui/Heading";
import AddNote from "../features/worknotes/AddNote";
import Row from "../ui/Row";
import NotesSummary from "../features/worknotes/NotesSummary";

function Notes() {
  return (
    <>
      <Row type="horizontal">
        <Heading as="h1">Work notes</Heading>
        <AddNote />
      </Row>

      <Row>
        <NotesSummary />
      </Row>
    </>
  );
}

export default Notes;
