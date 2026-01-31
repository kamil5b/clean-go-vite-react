import { z } from "zod";

export const SaveEmailRequestSchema = z.object({
    to: z.string().email("Invalid email address"),
    subject: z.string().min(1, "Subject is required"),
    body: z.string().min(1, "Body is required"),
});

export type SaveEmailRequest = z.infer<typeof SaveEmailRequestSchema>;
