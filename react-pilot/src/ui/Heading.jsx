import styled, { css } from "styled-components";

// 加上 css 之后，在 ${} 写的 js 就可以生效了
// const test = css`
//   text-align: center;
//   ${10 > 50 && "background-color: var(--color-brand-100)"}
// `;

const Heading = styled.h1`
  ${(props) =>
    props.as === "h1" &&
    css`
      font-size: 3rem;
      font-weight: 600;
    `}

  ${(props) =>
    props.as === "h2" &&
    css`
      font-size: 2rem;
      font-weight: 600;
    `}

    ${(props) =>
    props.as === "h3" &&
    css`
      font-size: 2rem;
      font-weight: 500;
    `}

    ${(props) =>
    props.as === "h4" &&
    css`
      font-size: 3rem;
      font-weight: 600;
      text-align: center;
    `}

  line-height:1.4
`;
export default Heading;
