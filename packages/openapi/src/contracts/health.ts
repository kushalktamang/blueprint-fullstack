import ZHealthResponse from "../../../zod/dist/health.js";
import { initContract } from "@ts-rest/core";

const contract = initContract();

const healthContract = contract.router({
  getHealth: {
    description: "Get health status",
    method: "GET",
    path: "/status",
    responses: {
      200: ZHealthResponse,
    },
    summary: "Get health",
  },
});

export default healthContract;
