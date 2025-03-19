import styled from "styled-components";
import { useState } from "react";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css"; // 引入 Quill 的样式

import Button from "../../ui/Button";
import ButtonGroup from "../../ui/ButtonGroup";

import Input from "../../ui/Input";
import Heading from "../../ui/Heading";

// import ButtonText from "../../ui/ButtonText";
import Row from "../../ui/Row";
// import { useNavigate } from "react-router-dom";

const HeadingGroup = styled.div`
  display: flex;
  gap: 2.4rem;
  align-items: center;
`;

const StyledNote = styled.div`
  /* Box */
  background-color: var(--color-grey-0);
  border: 1px solid var(--color-grey-100);
  border-radius: var(--border-radius-md);

  width: 800px;
  height: 800px;
  overflow: auto;

  padding: 3.2rem;
  display: flex;
  flex-direction: column;
  gap: 2.4rem;
  grid-column: 1 / span 2;
  padding-top: 2.4rem;
`;

function EditNotes() {
  const [inputText, setInputText] = useState("");
  const [content, setContent] = useState("");

  //   const [isSubmit, setIsSubmit] = useState(false);
  const [isPreview, setIsPreview] = useState(false);

  //   const navigate = useNavigate();

  const handleChange = (value) => {
    setContent(value);
  };

  function handlePreview() {
    setIsPreview(!isPreview);
  }

  //   function handleSubmit(e) {
  //     e.preventDefault();
  //     console.log(content);
  //     setIsSubmit(true);
  //     setIsPreview(false);
  //   }

  return (
    <StyledNote>
      <>
        <Row type="horizontal">
          <HeadingGroup>
            {/* <Heading as="h2">Editing a new note</Heading> */}
            {/* <Tag type={statusToTagName[status]}>{status.replace("-", " ")}</Tag> */}
          </HeadingGroup>
          {/* <ButtonText onClick={() => navigate(`/notes`)}>
            &larr; Back
          </ButtonText> */}
        </Row>
        <Input
          id="title"
          type="text"
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          placeholder="Type your title..."
        />
        <div>
          <ReactQuill
            id="content"
            value={content}
            onChange={handleChange}
            placeholder="Type your content..."
          />
        </div>

        <ButtonGroup>
          {/* <Button onClick={handleSubmit}>Submit</Button> */}
          <Button>Submit</Button>
          <Button onClick={handlePreview}>Preview</Button>
        </ButtonGroup>

        {isPreview ? (
          <span>
            <Heading as="h1">{inputText}</Heading>
            <div dangerouslySetInnerHTML={{ __html: content }} />
          </span>
        ) : (
          ""
        )}
      </>
    </StyledNote>
  );
}

export default EditNotes;
