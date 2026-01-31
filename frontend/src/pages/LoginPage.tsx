import { Link } from "react-router-dom";
import { useLogin } from "@/hooks/useLogin";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
    Field,
    FieldGroup,
    FieldLabel,
    FieldContent,
} from "@/components/ui/field";
import { PasswordInput } from "@/components/PasswordInput";

export function LoginPage() {
    const {
        email,
        setEmail,
        password,
        setPassword,
        error,
        loading,
        handleSubmit,
    } = useLogin();

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-50">
            <Card className="w-full max-w-md p-8">
                <h1 className="text-2xl font-bold mb-6 text-center">Login</h1>

                {error && (
                    <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit}>
                    <FieldGroup>
                        <Field>
                            <FieldLabel htmlFor="email">Email</FieldLabel>
                            <FieldContent>
                                <input
                                    id="email"
                                    type="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    required
                                    disabled={loading}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                                    placeholder="you@example.com"
                                />
                            </FieldContent>
                        </Field>

                        <Field>
                            <FieldLabel htmlFor="password">Password</FieldLabel>
                            <FieldContent>
                                <PasswordInput
                                    id="password"
                                    value={password}
                                    onChange={(e) =>
                                        setPassword(e.target.value)
                                    }
                                    required
                                    disabled={loading}
                                />
                            </FieldContent>
                        </Field>
                    </FieldGroup>

                    <Button
                        type="submit"
                        disabled={loading}
                        className="w-full mt-6 mb-4"
                    >
                        {loading ? "Logging in..." : "Login"}
                    </Button>
                </form>

                <p className="text-center text-sm text-gray-600">
                    Don't have an account?{" "}
                    <Link
                        to="/register"
                        className="text-blue-600 hover:text-blue-700 font-medium"
                    >
                        Register here
                    </Link>
                </p>
            </Card>
        </div>
    );
}
