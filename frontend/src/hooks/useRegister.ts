import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";

interface UseRegisterReturn {
    name: string;
    setName: (name: string) => void;
    email: string;
    setEmail: (email: string) => void;
    password: string;
    setPassword: (password: string) => void;
    confirmPassword: string;
    setConfirmPassword: (confirmPassword: string) => void;
    error: string;
    loading: boolean;
    handleSubmit: (e: React.FormEvent) => Promise<void>;
    reset: () => void;
}

export function useRegister(): UseRegisterReturn {
    const navigate = useNavigate();
    const { register } = useAuth();
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (password !== confirmPassword) {
            setError("Passwords do not match");
            return;
        }

        if (password.length < 8) {
            setError("Password must be at least 8 characters");
            return;
        }

        setLoading(true);

        try {
            await register(email, password, name);
            navigate("/");
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Registration failed",
            );
        } finally {
            setLoading(false);
        }
    };

    const reset = () => {
        setName("");
        setEmail("");
        setPassword("");
        setConfirmPassword("");
        setError("");
        setLoading(false);
    };

    return {
        name,
        setName,
        email,
        setEmail,
        password,
        setPassword,
        confirmPassword,
        setConfirmPassword,
        error,
        loading,
        handleSubmit,
        reset,
    };
}
