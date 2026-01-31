import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

interface PasswordInputProps {
    id: string;
    value: string;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
}

export function PasswordInput({
    id,
    value,
    onChange,
    placeholder = "Enter your password",
    disabled = false,
    required = false,
}: PasswordInputProps) {
    const [showPassword, setShowPassword] = useState(false);

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

    return (
        <div className="relative w-full">
            <input
                id={id}
                type={showPassword ? "text" : "password"}
                value={value}
                onChange={onChange}
                required={required}
                disabled={disabled}
                className="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                placeholder={placeholder}
            />
            <button
                type="button"
                onClick={togglePasswordVisibility}
                disabled={disabled}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none"
                aria-label={showPassword ? "Hide password" : "Show password"}
            >
                {showPassword ? (
                    <EyeOff className="w-5 h-5" />
                ) : (
                    <Eye className="w-5 h-5" />
                )}
            </button>
        </div>
    );
}
