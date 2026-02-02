import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";

interface UseLoginReturn {
    email: string;
    setEmail: (email: string) => void;
    password: string;
    setPassword: (password: string) => void;
    error: string;
    loading: boolean;
    handleSubmit: (e: React.FormEvent) => Promise<void>;
    reset: () => void;
}

export function useLogin(): UseLoginReturn {
    const navigate = useNavigate();
    const { login } = useAuth();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setLoading(true);

        try {
            await login(email, password);
            navigate("/");
        } catch (err) {
            setError(err instanceof Error ? err.message : "Login failed");
        } finally {
            setLoading(false);
        }
    };

    const reset = () => {
        setEmail("");
        setPassword("");
        setError("");
        setLoading(false);
    };

    return {
        email,
        setEmail,
        password,
        setPassword,
        error,
        loading,
        handleSubmit,
        reset,
    };
}
