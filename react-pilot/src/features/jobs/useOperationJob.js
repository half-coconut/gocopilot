import { useQuery } from "@tanstack/react-query";
import openJob from "../../services/apiJobs";
import { useParams } from "react-router-dom";

export function useOpenJob() {
  const { jobId } = useParams();
  const {
    isLoading,
    data: jobData,
    error,
  } = useQuery({
    queryKey: ["job_open", jobId],
    queryFn: () => openJob(jobId),
    retry: false,
  });

  let res;

  if (jobData?.data == jobId) {
    res = true;
  }
  // console.log("debugInterfaces: ", debugInterfacesItem);
  // console.log(JSON.stringify(debugInterfacesItem, null, 2));
  return { isLoading, error, res };
}
