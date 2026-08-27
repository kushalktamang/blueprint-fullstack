import { API_URL } from "@/config/env.config";
import apiContract  from "@blueprint/openapi/contracts";
import { useAuth } from "@clerk/clerk-react";
import { initClient } from "@ts-rest/core";
import axios, { isAxiosError} from "axios";


type Headers = Awaited<ReturnType<NonNullable<Parameters<typeof initClient>[1]["api"]>>>["headers"];

interface ApiResponse {
  status: number;
  body: unknown;
  headers: Headers;
}

export type TApiClient = ReturnType<typeof useApiClient>;

export const useApiClient = ({ isBlob = false }: { isBlob?: boolean } = {}) => {
  const { getToken } = useAuth();
  return initClient(apiContract, {
    baseUrl: "",
    baseHeaders: {
      "Content-Type": "application/json",
    },
    api: async ({ path, method, headers, body }) => {
      const token = await getToken({ template: "custom" });

      const makeRequest = async (retryCount = 0): Promise<ApiResponse> => {
        try {
          const result = await axios.request({
            method,
            url: `${API_URL}/api${path}`,
            headers: {
              ...headers,
              ...(token !== null && token !== "" ? { Authorization: `Bearer ${token}` } : {}),
            },
            data: body,
            ...(isBlob ? { responseType: "blob" } : {}),
          });
          return {
            status: result.status,
            body: result.data,
            // The axios headers shape and the ts-rest fetcher's expected Headers shape
            // are structurally compatible but not provably identical to TS — no
            // runtime validation exists at this boundary either way.
            // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion
            headers: result.headers as unknown as Headers,
          };
        } catch (error: unknown) {
          if (isAxiosError(error)) {
            const { response } = error;
            // If unauthorized and we haven't retried yet, retry
            if (response?.status === 401 && retryCount < 2) {
              return makeRequest(retryCount + 1);
            }
            return {
              status: response?.status ?? 500,
              body: response?.data ?? { message: "Internal server error" },
              // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion
              headers: (response?.headers as unknown as Headers) ?? {},
            };
          }
          throw error;
        }
      };
      return makeRequest();
    },
  });
};
