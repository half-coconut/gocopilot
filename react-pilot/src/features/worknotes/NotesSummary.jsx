import styled from "styled-components";
import PropTypes from "prop-types";
import { MdFavorite, MdOutlineMenuBook } from "react-icons/md";
import { FaStar } from "react-icons/fa";

import Heading from "../../ui/Heading";

const StyleUl = styled.ul`
  margin-top: 10px;
  color: var(--color-grey-500);
`;

const ButtonContainer = styled.div`
  display: flex;
  gap: 20px; /* 按钮之间的间距 */
  padding: 2.5px; /* 容器的内边距 */
  background-color: var(--color-grey-0); /* 背景色 */
  border-radius: 8px; /* 圆角 */
  justify-content: flex-end; /* 关键：使用 flex-end 将按钮右对齐 */
`;

const Button = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent; /* 背景透明 */
  border: none; /* 去掉边框 */
  color: var(--color-grey-500); /* 按钮文字颜色 */
  font-size: 16px; /* 字体大小 */
  cursor: pointer; /* 鼠标悬停时变成手型 */
  padding: 10px; /* 按钮的内边距 */

  &:hover {
    color: var(--color-grey-700); /* 悬停时改变颜色 */
  }

  & svg {
    margin-right: 5px; /* 图标与文字之间的间距 */
  }
`;

const IconButton = ({ icon, text }) => (
  <Button>
    {icon}
    {text}
  </Button>
);

IconButton.propTypes = {
  icon: PropTypes.string,
  text: PropTypes.string,
};

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

const counts = {
  read: 18,
  like: 17,
  collect: 16,
};

const notesList = [
  {
    title: "Is vaping as bad as smoking cigarettes?",
    content:
      "53% of current vapers in the UK, that's around 3 million people, used to smoke cigarettes. Many believe that vaping is less harmful than smoking cigarettes. A new study, which will be published soon, questions this.The study, by Manchester Metropolitan University, has found that e-cigarettes and vapes are just as bad for health as cigarettes....",
  },
  {
    title: "Woolly mice: Are woolly mammoths next?",
    content:
      "Scientists have created a genetically modified mouse that's woolly.The researchers plan to use their woolly mouse to test out other genetic changes before they try to create genetically-altered, mammoth-like elephants in the future.The company, Colossal Biosciences, hope to use the new mammoths in the fight against global warming....",
  },
  {
    title: "50% adults overweight or obese by 2050: Global study",
    content:
      "More than half of all adults and a third of children, teenagers and young adults around the world are predicted to be overweight or obese by 2050. The findings come from a new study of global data, covering more than 200 countries, published in The Lancet, a well-known British medical journal....",
  },
];

function NotesSummary() {
  console.log(notesList);
  return (
    <StyledNote>
      {notesList.map((note) => (
        <span key={note.title}>
          <Heading as="h2">{note.title}</Heading>
          <StyleUl>{note.content}</StyleUl>

          <ButtonContainer>
            <IconButton icon={<MdOutlineMenuBook />} text={counts.read} />
            <IconButton icon={<MdFavorite />} text={counts.like} />
            <IconButton icon={<FaStar />} text={counts.collect} />
            {/* <IconButton icon={<FaRetweet />} text="419" /> */}
          </ButtonContainer>
        </span>
      ))}
    </StyledNote>
  );
}

export default NotesSummary;
