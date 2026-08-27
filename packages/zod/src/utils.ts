import { z } from "zod";

interface PaginatedResponse<ItemType> {
  data: ItemType[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

const schemaWithPagination = <ItemType>(
  schema: z.ZodType<ItemType>,
): z.ZodType<PaginatedResponse<ItemType>> =>
  z.object({
    data: z.array(schema),
    total: z.number(),
    page: z.number(),
    limit: z.number(),
    totalPages: z.number(),
  });

export default schemaWithPagination;
