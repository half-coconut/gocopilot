import { useEffect } from "react";
import PropTypes from "prop-types";
import { useLocalStorageState } from "../hooks/useLocalStorageState";
import { DarkModeContext } from "./useDarkMode";

export function DarkModeProvider({ children }) {
  // const [isDarkMode, setIsDarkMode] = useLocalStorageState(false, "isDarkMode");
  const [isDarkMode, setIsDarkMode] = useLocalStorageState(
    window.matchMedia("(prefers-color-scheme: dark)").matches,
    "isDarkMode"
  );

  useEffect(
    function () {
      if (isDarkMode) {
        document.documentElement.classList.add("dark-mode");
        document.documentElement.classList.remove("light-mode");
      } else {
        document.documentElement.classList.add("light-mode");
        document.documentElement.classList.remove("dark-mode");
      }
    },
    [isDarkMode]
  );

  function toggleDarkMode() {
    setIsDarkMode((isDark) => !isDark);
  }
  return (
    <DarkModeContext.Provider value={{ isDarkMode, toggleDarkMode }}>
      {children}
    </DarkModeContext.Provider>
  );
}
DarkModeProvider.propTypes = {
  children: PropTypes.node,
};
