import { z } from "zod";

export const CommonIDResponseSchema = z.object({
    value: z.string().uuid(),
});

export type CommonIDResponse = z.infer<typeof CommonIDResponseSchema>;
