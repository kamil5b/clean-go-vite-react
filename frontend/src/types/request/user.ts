import { z } from "zod";

export const RegisterUserRequestSchema = z.object({
    email: z.string().email("Invalid email address"),
    password: z.string().min(8, "Password must be at least 8 characters"),
    name: z.string().min(1, "Name is required"),
});

export const LoginRequestSchema = z.object({
    email: z.string().email("Invalid email address"),
    password: z.string().min(1, "Password is required"),
});

export type RegisterUserRequest = z.infer<typeof RegisterUserRequestSchema>;
export type LoginRequest = z.infer<typeof LoginRequestSchema>;
