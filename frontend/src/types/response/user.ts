import { z } from "zod";

export const GetUserSchema = z.object({
    id: z.string().uuid(),
    email: z.string().email(),
    name: z.string(),
});

export type GetUser = z.infer<typeof GetUserSchema>;

export const LoginResponseSchema = z.object({
    token: z.string(),
    user: GetUserSchema,
});

export type LoginResponse = z.infer<typeof LoginResponseSchema>;

export const RegisterResponseSchema = z.object({
    token: z.string(),
    user: GetUserSchema,
});

export type RegisterResponse = z.infer<typeof RegisterResponseSchema>;

export const RefreshResponseSchema = z.object({
    token: z.string(),
});

export type RefreshResponse = z.infer<typeof RefreshResponseSchema>;

export const CSRFTokenResponseSchema = z.object({
    token: z.string(),
});

export type CSRFTokenResponse = z.infer<typeof CSRFTokenResponseSchema>;
