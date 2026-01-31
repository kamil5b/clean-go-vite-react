import { z } from "zod";

export const GetMessageSchema = z.object({
    content: z.string(),
});

export type GetMessage = z.infer<typeof GetMessageSchema>;
