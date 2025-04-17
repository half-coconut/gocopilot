// import PropTypes from "prop-types";
import jobService from "./apiService/jobService.js";

export async function createEditJob(newJob) {
  return jobService.edit(newJob);
}

export default async function openJob(id) {
  const resp = await jobService.open(id);
  return resp;
}

export async function closeJob(id) {
  const resp = await jobService.close(id);
  return resp;
}
