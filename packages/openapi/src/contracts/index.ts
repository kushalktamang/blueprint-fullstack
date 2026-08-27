import healthContract from "./health.js";
import { initContract } from "@ts-rest/core";

const contract = initContract();

const apiContract = contract.router({
  Health: healthContract,
});

export default apiContract;
