import { useState } from "react";
import { useItems } from "@/hooks/useItems";
import {
    Card,
    CardHeader,
    CardTitle,
    CardContent,
    CardDescription,
} from "@/components/ui/card";
import {
    Table,
    TableHeader,
    TableBody,
    TableRow,
    TableHead,
    TableCell,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import {
    Field,
    FieldGroup,
    FieldLabel,
    FieldContent,
    FieldDescription,
} from "@/components/ui/field";

export function ItemsPage() {
    const { items, loading, error, createItem, updateItem, deleteItem } =
        useItems();
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [editingId, setEditingId] = useState<string | null>(null);
    const [editTitle, setEditTitle] = useState("");
    const [editDescription, setEditDescription] = useState("");
    const [submitError, setSubmitError] = useState("");

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitError("");
        if (!title.trim()) {
            setSubmitError("Title is required");
            return;
        }
        try {
            await createItem(title, description);
            setTitle("");
            setDescription("");
        } catch (err) {
            setSubmitError(
                err instanceof Error ? err.message : "Failed to create item",
            );
        }
    };

    const handleUpdate = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitError("");
        if (!editingId) return;
        if (!editTitle.trim()) {
            setSubmitError("Title is required");
            return;
        }
        try {
            await updateItem(editingId, editTitle, editDescription);
            setEditingId(null);
            setEditTitle("");
            setEditDescription("");
        } catch (err) {
            setSubmitError(
                err instanceof Error ? err.message : "Failed to update item",
            );
        }
    };

    const handleDelete = async (id: string) => {
        if (confirm("Are you sure you want to delete this item?")) {
            try {
                await deleteItem(id);
            } catch (err) {
                setSubmitError(
                    err instanceof Error
                        ? err.message
                        : "Failed to delete item",
                );
            }
        }
    };

    const startEdit = (
        id: string,
        itemTitle: string,
        itemDescription: string,
    ) => {
        setEditingId(id);
        setEditTitle(itemTitle);
        setEditDescription(itemDescription);
    };

    return (
        <div className="min-h-screen bg-linear-to-b from-slate-50 to-slate-100 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-6xl mx-auto space-y-8">
                <div>
                    <h1 className="text-4xl font-bold text-slate-900">Items</h1>
                    <p className="text-slate-600 mt-2">
                        Manage your items with ease
                    </p>
                </div>

                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
                        {error}
                    </div>
                )}

                {submitError && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
                        {submitError}
                    </div>
                )}

                {/* Create/Edit Form Card */}
                <Card>
                    <CardHeader>
                        <CardTitle>
                            {editingId ? "Edit Item" : "Create New Item"}
                        </CardTitle>
                        <CardDescription>
                            {editingId
                                ? "Update the item details below"
                                : "Add a new item to your collection"}
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form
                            onSubmit={editingId ? handleUpdate : handleCreate}
                            className="space-y-6"
                        >
                            <FieldGroup>
                                <Field>
                                    <FieldLabel htmlFor="title">
                                        Title
                                    </FieldLabel>
                                    <FieldContent>
                                        <input
                                            id="title"
                                            type="text"
                                            value={
                                                editingId ? editTitle : title
                                            }
                                            onChange={(e) =>
                                                editingId
                                                    ? setEditTitle(
                                                          e.target.value,
                                                      )
                                                    : setTitle(e.target.value)
                                            }
                                            placeholder="Enter item title"
                                            className="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                                        />
                                        <FieldDescription>
                                            Give your item a meaningful name
                                        </FieldDescription>
                                    </FieldContent>
                                </Field>
                            </FieldGroup>

                            <FieldGroup>
                                <Field>
                                    <FieldLabel htmlFor="description">
                                        Description
                                    </FieldLabel>
                                    <FieldContent>
                                        <textarea
                                            id="description"
                                            value={
                                                editingId
                                                    ? editDescription
                                                    : description
                                            }
                                            onChange={(e) =>
                                                editingId
                                                    ? setEditDescription(
                                                          e.target.value,
                                                      )
                                                    : setDescription(
                                                          e.target.value,
                                                      )
                                            }
                                            placeholder="Enter item description"
                                            rows={4}
                                            className="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none"
                                        />
                                        <FieldDescription>
                                            Optional: Add more details about
                                            this item
                                        </FieldDescription>
                                    </FieldContent>
                                </Field>
                            </FieldGroup>

                            <div className="flex gap-3 justify-end">
                                <Button
                                    type="submit"
                                    variant="default"
                                    size="default"
                                >
                                    {editingId ? "Update Item" : "Create Item"}
                                </Button>
                                {editingId && (
                                    <Button
                                        type="button"
                                        variant="outline"
                                        size="default"
                                        onClick={() => {
                                            setEditingId(null);
                                            setEditTitle("");
                                            setEditDescription("");
                                        }}
                                    >
                                        Cancel
                                    </Button>
                                )}
                            </div>
                        </form>
                    </CardContent>
                </Card>

                {/* Items List Card */}
                <Card>
                    <CardHeader>
                        <CardTitle>Items List</CardTitle>
                        <CardDescription>
                            {items.length === 0
                                ? "No items yet"
                                : `You have ${items.length} item${items.length !== 1 ? "s" : ""}`}
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        {loading ? (
                            <div className="text-center py-12">
                                <p className="text-slate-600">
                                    Loading items...
                                </p>
                            </div>
                        ) : items.length === 0 ? (
                            <div className="text-center py-12">
                                <p className="text-slate-600">
                                    No items yet. Create one to get started!
                                </p>
                            </div>
                        ) : (
                            <div className="overflow-x-auto">
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead className="w-1/3">
                                                Title
                                            </TableHead>
                                            <TableHead className="w-1/3">
                                                Description
                                            </TableHead>
                                            <TableHead className="w-1/6">
                                                Created
                                            </TableHead>
                                            <TableHead className="text-right w-1/6">
                                                Actions
                                            </TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {items.map((item) => (
                                            <TableRow key={item.id}>
                                                <TableCell className="font-semibold text-slate-900">
                                                    {item.title}
                                                </TableCell>
                                                <TableCell className="text-slate-600">
                                                    {item.description || "-"}
                                                </TableCell>
                                                <TableCell className="text-slate-500 text-sm">
                                                    {new Date(
                                                        item.created_at,
                                                    ).toLocaleDateString()}
                                                </TableCell>
                                                <TableCell className="text-right">
                                                    <div className="flex gap-2 justify-end">
                                                        <Button
                                                            onClick={() =>
                                                                startEdit(
                                                                    item.id,
                                                                    item.title,
                                                                    item.description,
                                                                )
                                                            }
                                                            variant="outline"
                                                            size="sm"
                                                        >
                                                            Edit
                                                        </Button>
                                                        <Button
                                                            onClick={() =>
                                                                handleDelete(
                                                                    item.id,
                                                                )
                                                            }
                                                            variant="destructive"
                                                            size="sm"
                                                        >
                                                            Delete
                                                        </Button>
                                                    </div>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
