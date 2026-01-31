# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## Project Structure Rules

### Type Organization

**RULE: All types used in `src/api/` MUST be imported from `src/types/`**

- Define all request and response types in `src/types/request/` and `src/types/response/`
- API modules should only import types, never define them
- Use Zod schemas for validation in the types folder
- This ensures a single source of truth for all API contracts

Example:
```typescript
// ✅ CORRECT - types/response/user.ts
export const LoginResponseSchema = z.object({
  token: z.string(),
  user: GetUserSchema,
});
export type LoginResponse = z.infer<typeof LoginResponseSchema>;

// ✅ CORRECT - api/auth.ts
import { LoginResponse } from "@/types/response/user";
export const authApi = {
  login: async (): Promise<LoginResponse> => { ... }
}

// ❌ WRONG - defining types directly in api/auth.ts
export interface AuthResponse { ... }
```

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type aware lint rules:

- Configure the top-level `parserOptions` property like this:

```js
   parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    project: ['./tsconfig.json', './tsconfig.node.json'],
    tsconfigRootDir: __dirname,
   },
```

- Replace `plugin:@typescript-eslint/recommended` to `plugin:@typescript-eslint/recommended-type-checked` or `plugin:@typescript-eslint/strict-type-checked`
- Optionally add `plugin:@typescript-eslint/stylistic-type-checked`
- Install [eslint-plugin-react](https://github.com/jsx-eslint/eslint-plugin-react) and add `plugin:react/recommended` & `plugin:react/jsx-runtime` to the `extends` list