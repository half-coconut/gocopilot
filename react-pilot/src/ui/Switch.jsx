import PropTypes from "prop-types";
import styled from "styled-components";

// 定义开关容器
const SwitchContainer = styled.label`
  position: relative;
  display: inline-block;
  width: 50px;
  height: 27px;
`;

// 定义隐藏的复选框（checkbox），用于保持状态
const HiddenCheckbox = styled.input.attrs({ type: "checkbox" })`
  height: 0;
  width: 0;
  visibility: hidden;
`;

// 定义滑轨样式
const Slider = styled.span`
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--color-grey-300);
  transition: 0.4s;
  border-radius: 34px;

  // 如果禁用，则改变背景色并禁止光标变为手型
  ${({ disabled }) =>
    disabled &&
    `
    background-color: var(--color-grey-200); // 更浅的灰色
    cursor: not-allowed; // 禁止点击
  `}

  &:before {
    position: absolute;
    content: "";
    height: 22px;
    width: 22px;
    left: 2px;
    bottom: 2px;
    background-color: white;
    transition: 0.4s;
    border-radius: 50%;
  }
`;

// 定义被选中时的滑轨样式
const SwitchStyled = styled(HiddenCheckbox)`
  &:checked + ${Slider} {
    background-color: var(--color-brand-600);
  }

  &:focus + ${Slider} {
    box-shadow: 0 0 1px var(--color-brand-50);
  }

  &:checked + ${Slider}:before {
    transform: translateX(26px);
  }
`;

function Switch({ checked, onChange, disabled }) {
  return (
    <SwitchContainer>
      <SwitchStyled checked={checked} onChange={onChange} disabled={disabled} />
      <Slider />
    </SwitchContainer>
  );
}

Switch.propTypes = {
  checked: PropTypes.node.isRequired,
  onChange: PropTypes.func.isRequired,
  disabled: PropTypes.bool,
};

Switch.defaultProps = {
  disabled: false, // 默认不禁用
};

export default Switch;
