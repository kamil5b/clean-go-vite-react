import { Button } from "@/components/ui/button";
import { useCounter } from "@/hooks/useCounter";

export function CounterPage() {
    const { count, increment, loading, error } = useCounter();

    return (
        <div className="flex items-center justify-center min-h-screen">
            <div className="flex flex-col items-center gap-8">
                <div
                    className="text-6xl font-bold"
                    style={{ fontSize: "64px" }}
                >
                    {count}
                </div>
                <Button onClick={increment} disabled={loading}>
                    {loading ? "Loading..." : "Increment"}
                </Button>
                {error && <p className="text-red-500 text-sm">{error}</p>}
            </div>
        </div>
    );
}
