import styled from "styled-components";
import Heading from "../../ui/Heading";
import ButtonGroup from "../../ui/ButtonGroup";

const Button = styled.button`
  border: none;
  border-radius: var(--border-radius-sm);
  box-shadow: var(--shadow-sm);

  color: var(--color-grey-500);
  background-color: var(--color-grey-0);

  font-size: 1.2rem;
  padding: 0.4rem 0.8rem;
  /* text-transform: uppercase; */
  /* font-weight: 500; */
  text-align: center;

  &:hover {
    background-color: var(--color-grey-300);
  }
`;

const StyledNote = styled.div`
  /* Box */
  background-color: var(--color-grey-0);
  border: 1px solid var(--color-grey-100);
  border-radius: var(--border-radius-md);

  padding: 3.2rem;
  display: flex;
  flex-direction: column;
  gap: 2.4rem;
  grid-column: 1 / span 2;
  padding-top: 2.4rem;
`;

const Interactive = styled.div`
  display: flex;
  align-items: center;
  gap: 1.2rem;
  margin-bottom: 1.6rem;
  color: var(--color-grey-500);
  text-align: right;
  width: 100%;

  & p:first-of-type {
    font-weight: 200;
    /* color: var(--color-grey-700); */
    text-align: right;
    width: 100%;
  }
`;

const counts = {
  read: 11,
  like: 12,
  collect: 13,
};

const notesList = [
  {
    title: "Is vaping as bad as smoking cigarettes?",
    content: "Content Detail...",
  },
  {
    title: "Woolly mice: Are woolly mammoths next?",
    content: "Content Detail...",
  },
  {
    title: "50% adults overweight or obese by 2050: Global study",
    content: "Content Detail...",
  },
];

function NotesSummary() {
  console.log(notesList);
  return (
    <StyledNote>
      {notesList.map((note) => (
        <span key={note.title}>
          <Heading as="h2">{note.title}</Heading>
          <ul>{note.content}</ul>

          <Interactive>
            <p>
              <span>&bull;</span>
            </p>
            {/* <span>&bull;</span> */}

            <ButtonGroup>
              <Button>Read </Button>
              <p>{counts.read}</p>
              <span>&bull;</span>

              <Button>Like </Button>
              <p>{counts.like}</p>
              <span>&bull;</span>

              <Button>Collect </Button>
              <p>{counts.collect}</p>
            </ButtonGroup>
          </Interactive>
        </span>
      ))}
    </StyledNote>
  );
}

export default NotesSummary;
