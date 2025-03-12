import { NavLink } from "react-router-dom";
import styled from "styled-components";
import {
  HiOutlineCalendarDateRange,
  HiOutlineCog6Tooth,
  HiOutlineHome,
  HiOutlineChartPie,
  HiOutlineArchiveBox,
  // HiOutlineUser,
  HiOutlineMapPin,
  HiOutlineCommandLine,
} from "react-icons/hi2";

const NavList = styled.ul`
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
`;

const StyledNavLink = styled(NavLink)`
  &:link,
  &:visited {
    display: flex;
    align-items: center;
    gap: 1.2rem;

    color: var(--color-grey-600);
    font-size: 1.6rem;
    font-weight: 500;
    padding: 1.2rem 2.4rem;
    transition: all 0.3s;
  }

  /* This works because react-router places the active class on the active NavLink */
  &:hover,
  &:active,
  &.active:link,
  &.active:visited {
    color: var(--color-grey-800);
    background-color: var(--color-grey-50);
    border-radius: var(--border-radius-sm);
  }

  & svg {
    width: 2.4rem;
    height: 2.4rem;
    color: var(--color-grey-400);
    transition: all 0.3s;
  }

  &:hover svg,
  &:active svg,
  &.active:link svg,
  &.active:visited svg {
    color: var(--color-brand-600);
  }
`;

function MainNav() {
  return (
    <nav>
      <NavList>
        <li>
          <StyledNavLink to="/dashboard">
            <HiOutlineHome />
            <span>Home</span>
          </StyledNavLink>
        </li>
        <li>
          <StyledNavLink to="/bookings">
            <HiOutlineCalendarDateRange />
            <span>Plans</span>
          </StyledNavLink>
        </li>
        <li>
          <StyledNavLink to="/bookings">
            <HiOutlineChartPie />
            <span>Tasks</span>
          </StyledNavLink>
        </li>
        <li>
          <StyledNavLink to="/cabins">
            <HiOutlineArchiveBox />
            <span>Interfaces</span>
          </StyledNavLink>
        </li>

        <li>
          <StyledNavLink to="/">
            <HiOutlineCommandLine />
            <span>AI Copilot</span>
          </StyledNavLink>
        </li>
        <li>
          <StyledNavLink to="/">
            <HiOutlineMapPin />
            <span>Environments</span>
          </StyledNavLink>
        </li>
        <li>
          <StyledNavLink to="/settings">
            <HiOutlineCog6Tooth />
            <span>Settings</span>
          </StyledNavLink>
        </li>
        {/* 用户注册的页面 */}
        {/* <li>
          <StyledNavLink to="/users">
            <HiOutlineUser />
            <span>Users</span>
          </StyledNavLink>
        </li> */}
      </NavList>
    </nav>
  );
}

export default MainNav;
