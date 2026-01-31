import { z } from "zod";

export const GetEmailLogSchema = z.object({
    id: z.string().uuid(),
    to: z.string().email(),
    subject: z.string(),
    body: z.string(),
    status: z.enum(["pending", "sent", "failed"]),
    created_at: z.string(),
});

export type GetEmailLog = z.infer<typeof GetEmailLogSchema>;
