import { match } from "ts-pattern";

const getSecurityMetadata = ({
  security = true,
  securityType = "bearer",
}: Readonly<{
  security?: boolean;
  securityType?: "bearer" | "service";
}> = {}) => {
  const openApiSecurity = match(securityType)
    .with("bearer", () => [
      {
        bearerAuth: [],
      },
    ])
    .with("service", () => [
      {
        "x-service-token": [],
      },
    ])
    .exhaustive();

  return {
    ...(security && { openApiSecurity }),
  };
};

export default getSecurityMetadata;
