import apiContract from "./contracts/index.js";
import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { generateOpenApi } from "@ts-rest/open-api";
import { z } from "zod";

extendZodWithOpenApi(z);

const Generate_Open_Api_Options_Index = 2;
type GenerateOpenApiOptionsIndex = typeof Generate_Open_Api_Options_Index;

type SecurityRequirementObject = Record<string, string[]>;
type OperationMapper = NonNullable<
  Parameters<typeof generateOpenApi>[GenerateOpenApiOptionsIndex]
>["operationMapper"];

const hasSecurity = (
  metadata: unknown,
): metadata is { openApiSecurity: SecurityRequirementObject[] } =>
  metadata !== null &&
  metadata !== undefined &&
  typeof metadata === "object" &&
  "openApiSecurity" in metadata;

const buildSecurityExtension = (
  metadata: unknown,
): Partial<{ security: SecurityRequirementObject[] }> => {
  if (hasSecurity(metadata)) {
    return { security: metadata.openApiSecurity };
  }
  return {};
};

const operationMapper: OperationMapper = (operation, appRoute) => ({
  ...operation,
  ...buildSecurityExtension(appRoute.metadata),
});

const OpenAPI = Object.assign(
  generateOpenApi(
    apiContract,
    {
      info: {
        description: "Blueprint REST API - Documentation",
        title: "Blueprint REST API - Documentation",
        version: "1.0.0",
      },
      openapi: "3.0.2",
      servers: [
        {
          description: "Local Server",
          url: "http://localhost:8080",
        },
      ],
    },
    {
      operationMapper,
      setOperationId: true,
    },
  ),
  {
    components: {
      securitySchemes: {
        bearerAuth: {
          bearerFormat: "JWT",
          scheme: "bearer",
          type: "http",
        },
        "x-service-token": {
          in: "header",
          name: "x-service-token",
          type: "apiKey",
        },
      },
    },
  },
);

export default OpenAPI;
