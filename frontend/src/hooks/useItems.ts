import { useState, useEffect } from "react";
import { itemsApi, type Item } from "@/api/items";

interface UseItemsReturn {
    items: Item[];
    loading: boolean;
    error: string;
    createItem: (title: string, description: string) => Promise<void>;
    updateItem: (
        id: string,
        title: string,
        description: string,
    ) => Promise<void>;
    deleteItem: (id: string) => Promise<void>;
    refreshItems: () => Promise<void>;
}

export function useItems(): UseItemsReturn {
    const [items, setItems] = useState<Item[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const refreshItems = async () => {
        setLoading(true);
        setError("");
        try {
            const data = await itemsApi.getItems();
            setItems(data);
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Failed to fetch items",
            );
        } finally {
            setLoading(false);
        }
    };

    const createItem = async (title: string, description: string) => {
        setError("");
        try {
            const newItem = await itemsApi.createItem(title, description);
            setItems([newItem, ...items]);
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Failed to create item",
            );
            throw err;
        }
    };

    const updateItem = async (
        id: string,
        title: string,
        description: string,
    ) => {
        setError("");
        try {
            const updated = await itemsApi.updateItem(id, title, description);
            setItems(items.map((item) => (item.id === id ? updated : item)));
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Failed to update item",
            );
            throw err;
        }
    };

    const deleteItem = async (id: string) => {
        setError("");
        try {
            await itemsApi.deleteItem(id);
            setItems(items.filter((item) => item.id !== id));
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Failed to delete item",
            );
            throw err;
        }
    };

    useEffect(() => {
        refreshItems();
    }, []);

    return {
        items,
        loading,
        error,
        createItem,
        updateItem,
        deleteItem,
        refreshItems,
    };
}
