import { useState, useCallback, useEffect } from "react";
import { counterApi } from "@/api/counter";

interface UseCounterReturn {
    count: number;
    increment: () => Promise<void>;
    loading: boolean;
    error: string | null;
}

export function useCounter(): UseCounterReturn {
    const [count, setCount] = useState<number>(0);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchCounter = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await counterApi.getCounter();
                setCount(response.value);
            } catch (err) {
                setError(
                    err instanceof Error
                        ? err.message
                        : "Failed to fetch counter",
                );
            } finally {
                setLoading(false);
            }
        };

        fetchCounter();
    }, []);

    const increment = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await counterApi.incrementCounter();
            setCount(response.value);
        } catch (err) {
            setError(
                err instanceof Error
                    ? err.message
                    : "Failed to increment counter",
            );
        } finally {
            setLoading(false);
        }
    }, []);

    return { count, increment, loading, error };
}
