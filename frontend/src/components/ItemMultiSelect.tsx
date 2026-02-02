import { useState, useEffect, useCallback, useRef } from "react";
import { itemApi } from "@/api/item";
import { ItemResponse } from "@/types/response/item";
import { Button } from "@/components/ui/button";

import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command";
import { Check, ChevronsUpDown, X } from "lucide-react";
import { cn } from "@/lib/utils";

interface ItemMultiSelectProps {
    value: number[];
    onChange: (value: number[]) => void;
    placeholder?: string;
    disabled?: boolean;
}

export function ItemMultiSelect({
    value,
    onChange,
    placeholder = "Select items...",
    disabled = false,
}: ItemMultiSelectProps) {
    const [open, setOpen] = useState(false);
    const [items, setItems] = useState<ItemResponse[]>([]);
    const [selectedItems, setSelectedItems] = useState<ItemResponse[]>([]);
    const [search, setSearch] = useState("");
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);
    const [loading, setLoading] = useState(false);
    const observerTarget = useRef<HTMLDivElement>(null);

    // Fetch selected items details
    useEffect(() => {
        const fetchSelectedItems = async () => {
            if (value.length === 0) {
                setSelectedItems([]);
                return;
            }

            try {
                const promises = value.map((id) => itemApi.getById(id));
                const results = await Promise.all(promises);
                setSelectedItems(results);
            } catch (error) {
                console.error("Failed to fetch selected items:", error);
            }
        };

        fetchSelectedItems();
    }, [value]);

    // Fetch items with pagination
    const fetchItems = useCallback(
        async (
            pageNum: number,
            searchQuery: string,
            reset: boolean = false,
        ) => {
            if (loading) return;

            setLoading(true);
            try {
                const response = await itemApi.getAll(pageNum, 20, searchQuery);

                if (reset) {
                    setItems(response.data);
                } else {
                    setItems((prev) => [...prev, ...response.data]);
                }

                setHasMore(pageNum < response.meta.totalPage);
                setPage(pageNum);
            } catch (error) {
                console.error("Failed to fetch items:", error);
            } finally {
                setLoading(false);
            }
        },
        [loading],
    );

    // Initial fetch and search
    useEffect(() => {
        if (open) {
            fetchItems(1, search, true);
        }
    }, [open, search]);

    // Infinite scroll observer
    useEffect(() => {
        if (!open) return;

        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasMore && !loading) {
                    fetchItems(page + 1, search, false);
                }
            },
            { threshold: 0.1 },
        );

        const currentTarget = observerTarget.current;
        if (currentTarget) {
            observer.observe(currentTarget);
        }

        return () => {
            if (currentTarget) {
                observer.unobserve(currentTarget);
            }
        };
    }, [open, hasMore, loading, page, search, fetchItems]);

    const handleSelect = (item: ItemResponse) => {
        const newValue = value.includes(item.id)
            ? value.filter((id) => id !== item.id)
            : [...value, item.id];
        onChange(newValue);
    };

    const handleRemove = (itemId: number) => {
        onChange(value.filter((id) => id !== itemId));
    };

    const handleSearchChange = (searchValue: string) => {
        setSearch(searchValue);
        setPage(1);
        setItems([]);
        setHasMore(true);
    };

    return (
        <div className="space-y-2">
            <Popover open={open} onOpenChange={setOpen}>
                <PopoverTrigger asChild>
                    <Button
                        variant="outline"
                        role="combobox"
                        aria-expanded={open}
                        className="w-full justify-between"
                        disabled={disabled}
                    >
                        {value.length === 0 ? (
                            <span className="text-muted-foreground">
                                {placeholder}
                            </span>
                        ) : (
                            <span>{value.length} item(s) selected</span>
                        )}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                    </Button>
                </PopoverTrigger>
                <PopoverContent className="w-full p-0" align="start">
                    <Command shouldFilter={false}>
                        <CommandInput
                            placeholder="Search items..."
                            value={search}
                            onValueChange={handleSearchChange}
                        />
                        <CommandList>
                            <CommandEmpty>
                                {loading ? "Loading..." : "No items found."}
                            </CommandEmpty>
                            <CommandGroup>
                                <ScrollArea className="h-[300px]">
                                    {items.map((item) => (
                                        <CommandItem
                                            key={item.id}
                                            value={item.id.toString()}
                                            onSelect={() => handleSelect(item)}
                                            className="cursor-pointer"
                                        >
                                            <Check
                                                className={cn(
                                                    "mr-2 h-4 w-4",
                                                    value.includes(item.id)
                                                        ? "opacity-100"
                                                        : "opacity-0",
                                                )}
                                            />
                                            <div className="flex-1">
                                                <div className="font-medium">
                                                    {item.name}
                                                </div>
                                                {item.desc && (
                                                    <div className="text-xs text-muted-foreground truncate">
                                                        {item.desc}
                                                    </div>
                                                )}
                                            </div>
                                        </CommandItem>
                                    ))}
                                    {hasMore && (
                                        <div
                                            ref={observerTarget}
                                            className="py-2 text-center text-sm text-muted-foreground"
                                        >
                                            {loading
                                                ? "Loading more..."
                                                : "Scroll for more"}
                                        </div>
                                    )}
                                </ScrollArea>
                            </CommandGroup>
                        </CommandList>
                    </Command>
                </PopoverContent>
            </Popover>

            {/* Selected items badges */}
            {selectedItems.length > 0 && (
                <div className="flex flex-wrap gap-2">
                    {selectedItems.map((item) => (
                        <Badge
                            key={item.id}
                            variant="secondary"
                            className="gap-1"
                        >
                            {item.name}
                            <button
                                type="button"
                                onClick={() => handleRemove(item.id)}
                                className="ml-1 hover:bg-muted rounded-full"
                                disabled={disabled}
                            >
                                <X className="h-3 w-3" />
                            </button>
                        </Badge>
                    ))}
                </div>
            )}
        </div>
    );
}
