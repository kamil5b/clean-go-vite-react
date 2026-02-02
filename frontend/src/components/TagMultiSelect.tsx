import { useState, useEffect, useCallback, useRef } from "react";
import { tagApi } from "@/api/tag";
import { TagResponse } from "@/types/response/tag";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Check, ChevronsUpDown, X } from "lucide-react";
import { cn } from "@/lib/utils";

interface TagMultiSelectProps {
  value: number[];
  onChange: (value: number[]) => void;
  placeholder?: string;
  disabled?: boolean;
}

export function TagMultiSelect({
  value,
  onChange,
  placeholder = "Select tags...",
  disabled = false,
}: TagMultiSelectProps) {
  const [open, setOpen] = useState(false);
  const [tags, setTags] = useState<TagResponse[]>([]);
  const [selectedTags, setSelectedTags] = useState<TagResponse[]>([]);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const observerTarget = useRef<HTMLDivElement>(null);

  // Fetch selected tags details
  useEffect(() => {
    const fetchSelectedTags = async () => {
      if (value.length === 0) {
        setSelectedTags([]);
        return;
      }

      try {
        const promises = value.map((id) => tagApi.getById(id));
        const results = await Promise.all(promises);
        setSelectedTags(results);
      } catch (error) {
        console.error("Failed to fetch selected tags:", error);
      }
    };

    fetchSelectedTags();
  }, [value]);

  // Fetch tags with pagination
  const fetchTags = useCallback(
    async (pageNum: number, searchQuery: string, reset: boolean = false) => {
      if (loading) return;

      setLoading(true);
      try {
        const response = await tagApi.getAll(pageNum, 20, searchQuery);

        if (reset) {
          setTags(response.data);
        } else {
          setTags((prev) => [...prev, ...response.data]);
        }

        setHasMore(pageNum < response.meta.totalPage);
        setPage(pageNum);
      } catch (error) {
        console.error("Failed to fetch tags:", error);
      } finally {
        setLoading(false);
      }
    },
    [loading]
  );

  // Initial fetch and search
  useEffect(() => {
    if (open) {
      fetchTags(1, search, true);
    }
  }, [open, search]);

  // Infinite scroll observer
  useEffect(() => {
    if (!open) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loading) {
          fetchTags(page + 1, search, false);
        }
      },
      { threshold: 0.1 }
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
  }, [open, hasMore, loading, page, search, fetchTags]);

  const handleSelect = (tag: TagResponse) => {
    const newValue = value.includes(tag.id)
      ? value.filter((id) => id !== tag.id)
      : [...value, tag.id];
    onChange(newValue);
  };

  const handleRemove = (tagId: number) => {
    onChange(value.filter((id) => id !== tagId));
  };

  const handleSearchChange = (searchValue: string) => {
    setSearch(searchValue);
    setPage(1);
    setTags([]);
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
              <span className="text-muted-foreground">{placeholder}</span>
            ) : (
              <span>{value.length} tag(s) selected</span>
            )}
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-full p-0" align="start">
          <Command shouldFilter={false}>
            <CommandInput
              placeholder="Search tags..."
              value={search}
              onValueChange={handleSearchChange}
            />
            <CommandList>
              <CommandEmpty>
                {loading ? "Loading..." : "No tags found."}
              </CommandEmpty>
              <CommandGroup>
                <ScrollArea className="h-[300px]">
                  {tags.map((tag) => (
                    <CommandItem
                      key={tag.id}
                      value={tag.id.toString()}
                      onSelect={() => handleSelect(tag)}
                      className="cursor-pointer"
                    >
                      <Check
                        className={cn(
                          "mr-2 h-4 w-4",
                          value.includes(tag.id) ? "opacity-100" : "opacity-0"
                        )}
                      />
                      <div
                        className="w-4 h-4 rounded mr-2 border"
                        style={{ backgroundColor: tag.color_hex }}
                      />
                      <div className="flex-1">
                        <div className="font-medium">{tag.name}</div>
                      </div>
                    </CommandItem>
                  ))}
                  {hasMore && (
                    <div ref={observerTarget} className="py-2 text-center text-sm text-muted-foreground">
                      {loading ? "Loading more..." : "Scroll for more"}
                    </div>
                  )}
                </ScrollArea>
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      {/* Selected tags badges */}
      {selectedTags.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {selectedTags.map((tag) => (
            <Badge
              key={tag.id}
              className="gap-1 text-white"
              style={{ backgroundColor: tag.color_hex }}
            >
              {tag.name}
              <button
                type="button"
                onClick={() => handleRemove(tag.id)}
                className="ml-1 hover:bg-black/20 rounded-full"
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
