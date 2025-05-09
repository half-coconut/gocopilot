import { BrowserRouter, Route, Routes, Navigate } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Toaster } from "react-hot-toast";

import Dashboard from "./pages/Dashboard";
import Tasks from "./pages/Tasks";
import Interfaces from "./pages/Interfaces";
import Users from "./pages/Users";
import Settings from "./pages/Settings";
import Account from "./pages/Account";
import Login from "./pages/Login";
import PageNotFound from "./pages/PageNotFound";
import GlobalStyles from "./styles/GlobalStyle";
import AppLayout from "./ui/AppLayout";
import Task from "./pages/Task";
import Checkin from "./pages/Checkin";
import { ProtectedRoute } from "./ui/ProtectedRoute";
import { DarkModeProvider } from "./context/DarkModeContext";
import AICopilot from "./pages/AICopilot";
import Interface from "./pages/Interface";
import History from "./pages/History";
import Environments from "./pages/Environments";
import Notes from "./pages/Notes";
import Jobs from "./pages/Jobs";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // 数据自动刷新的时间 5min
      staleTime: 1000 * 60 * 5,
      // staleTime: 1000,
      // 缓存时间 10min
      cacheTime: 1000 * 60 * 10,
    },
  },
});

function App() {
  return (
    <DarkModeProvider>
      <QueryClientProvider client={queryClient}>
        <ReactQueryDevtools initialIsOpen={false} />

        <GlobalStyles />
        <BrowserRouter>
          <Routes>
            <Route
              element={
                <ProtectedRoute>
                  <AppLayout />
                </ProtectedRoute>
              }
            >
              <Route index element={<Navigate replace to="dashboard" />} />
              <Route path="dashboard" element={<Dashboard />} />
              <Route path="jobs" element={<Jobs />} />
              <Route path="tasks" element={<Tasks />} />
              <Route path="tasks/:taskId" element={<Task />} />
              <Route path="checkin/:bookingId" element={<Checkin />} />
              <Route path="interfaces" element={<Interfaces />} />
              <Route path="interfaces/:interfaceId" element={<Interface />} />
              <Route path="aicopilot" element={<AICopilot />} />
              <Route path="history" element={<History />} />
              <Route path="environments" element={<Environments />} />
              <Route path="notes" element={<Notes />} />
              <Route path="users" element={<Users />} />
              <Route path="settings" element={<Settings />} />
              <Route path="account" element={<Account />} />
            </Route>

            <Route path="login" element={<Login />} />
            <Route path="*" element={<PageNotFound />} />
          </Routes>
        </BrowserRouter>

        <Toaster
          position="top-center"
          gutter={12}
          containerStyle={{ margin: "8px" }}
          toastOptions={{
            success: {
              duration: 3000,
            },
            error: {
              duration: 5000,
            },
            style: {
              fontSize: "16px",
              maxWidth: "500px",
              padding: "16px",
              backgroundColor: "var(--color-grey-0)",
              color: "var(--color-grey-700)",
            },
          }}
        />
      </QueryClientProvider>
    </DarkModeProvider>
  );
}

export default App;
