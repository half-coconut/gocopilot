import styled from "styled-components";
// import { useDarkMode } from "../context/useDarkMode";

const StyledLogo = styled.div`
  text-align: center;
  height: 3.6rem;
  width: auto;
`;

const Img = styled.img`
  height: 7.2rem; /* 高度 */
  width: 7.2rem; /* 宽度，设置为与高度相同 */
  border-radius: 50%; /* 圆形边框 */
  object-fit: cover; /* 确保图片填充整个圆形区域 */
`;

function Logo() {
  // const { isDarkMode } = useDarkMode();

  // const src = isDarkMode ? "/logo-dark.png" : "/logo-light.png";
  const src = "/kitty.jpg";

  return (
    <StyledLogo>
      <Img src={src} alt="Logo" />
    </StyledLogo>
  );
}

export default Logo;
