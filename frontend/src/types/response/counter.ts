import { z } from "zod";

export const GetCounterSchema = z.object({
    value: z.number().int(),
});

export type GetCounter = z.infer<typeof GetCounterSchema>;
