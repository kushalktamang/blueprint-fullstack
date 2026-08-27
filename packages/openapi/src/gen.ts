import OpenAPI from "./index.js";
import fs from "fs";

const JSON_INDENT_SPACES = 2;

const replaceCustomFileTypesToOpenApiCompatible = (jsonString: string): string => {
  const searchPattern =
    /\{"type":"object","properties":\{"type":\{"type":"string","enum":\["file"\]\}\},\s*"required":\["type"\]\}/gu;
  const replacement = `{"type":"string","format":"binary"}`;

  return jsonString.replace(searchPattern, replacement);
};

const filteredDoc = replaceCustomFileTypesToOpenApiCompatible(JSON.stringify(OpenAPI));
const formattedDoc: unknown = JSON.parse(filteredDoc);
const filePaths = ["./openapi.json", "../../apps/server/static/openapi.json"];

filePaths.forEach((filePath) => {
  fs.writeFile(filePath, JSON.stringify(formattedDoc, null, JSON_INDENT_SPACES), (err) => {
    if (err !== null) {
      console.error(`Error writing to ${filePath}:`, err);
    }
  });
});
