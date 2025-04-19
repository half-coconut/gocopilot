import { useQuery } from "@tanstack/react-query";

import { useParams } from "react-router-dom";

import {
  debugInterfaces,
  debugTask,
  executeTask,
} from "../../services/apiTasks";

export function useDebugInterfaces() {
  const { taskId } = useParams();
  const {
    isLoading,
    data: taskData,
    error,
  } = useQuery({
    queryKey: ["task_interfaces_debug", taskId],
    queryFn: () => debugInterfaces(taskId),
    retry: false,
  });

  const debugInterfacesItem = taskData?.data;
  // console.log("debugInterfaces: ", debugInterfacesItem);
  // console.log(JSON.stringify(debugInterfacesItem, null, 2));
  return { isLoading, error, debugInterfacesItem };
}

export function useDebugTask() {
  const { taskId } = useParams();
  const {
    isLoading,
    data: taskData,
    error,
  } = useQuery({
    queryKey: ["task_debug", taskId],
    queryFn: () => debugTask(taskId),
    retry: false,
  });

  const debugTaskItem = taskData?.data;

  return { isLoading, error, debugTaskItem };
}

export function useExecuteTask() {
  const { taskId } = useParams();
  const {
    isLoading,
    data: taskData,
    error,
  } = useQuery({
    queryKey: ["task_execute", taskId],
    queryFn: () => executeTask(taskId),
    retry: false,
  });

  const executeTaskItem = taskData?.data;

  return { isLoading, error, executeTaskItem };
}
