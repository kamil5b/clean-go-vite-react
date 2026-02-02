import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { invoiceApi } from "@/api/invoice";
import { itemApi } from "@/api/item";
import { InvoiceItemInput } from "@/types/request/invoice";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ItemMultiSelect } from "@/components/ItemMultiSelect";
import { TagMultiSelect } from "@/components/TagMultiSelect";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Trash2 } from "lucide-react";

interface InvoiceItemForm {
    item_id: number;
    itemName?: string;
    quantity: number;
    unit_price: number;
    total_price: number;
}

export default function InvoiceFormPage() {
    const navigate = useNavigate();
    const { id } = useParams();
    const isEditMode = !!id;

    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState("");

    // Form state
    const [selectedItemIds, setSelectedItemIds] = useState<number[]>([]);
    const [selectedTagIds, setSelectedTagIds] = useState<number[]>([]);
    const [invoiceItems, setInvoiceItems] = useState<InvoiceItemForm[]>([]);

    // Fetch invoice data if editing
    useEffect(() => {
        if (isEditMode && id) {
            fetchInvoice(parseInt(id));
        }
    }, [id, isEditMode]);

    const fetchInvoice = async (invoiceId: number) => {
        setLoading(true);
        try {
            const invoice = await invoiceApi.getById(invoiceId);

            // Set items
            const itemForms: InvoiceItemForm[] = invoice.items.map((item) => ({
                item_id: item.item_id,
                itemName: item.item.name,
                quantity: item.quantity,
                unit_price: item.unit_price,
                total_price: item.total_price,
            }));
            setInvoiceItems(itemForms);
            setSelectedItemIds(itemForms.map((item) => item.item_id));

            // Set tags
            setSelectedTagIds(invoice.tags.map((tag) => tag.id));
        } catch (err: any) {
            setError(err.message || "Failed to load invoice");
        } finally {
            setLoading(false);
        }
    };

    // Handle item selection changes
    const handleItemsChange = async (newItemIds: number[]) => {
        const addedIds = newItemIds.filter(
            (id) => !selectedItemIds.includes(id),
        );
        const removedIds = selectedItemIds.filter(
            (id) => !newItemIds.includes(id),
        );

        setSelectedItemIds(newItemIds);

        // Remove items that were deselected
        if (removedIds.length > 0) {
            setInvoiceItems((prev) =>
                prev.filter((item) => !removedIds.includes(item.item_id)),
            );
        }

        // Add new items
        if (addedIds.length > 0) {
            try {
                const promises = addedIds.map((id) => itemApi.getById(id));
                const items = await Promise.all(promises);

                const newInvoiceItems: InvoiceItemForm[] = items.map(
                    (item) => ({
                        item_id: item.id,
                        itemName: item.name,
                        quantity: 1,
                        unit_price: 0,
                        total_price: 0,
                    }),
                );

                setInvoiceItems((prev) => [...prev, ...newInvoiceItems]);
            } catch (err) {
                console.error("Failed to fetch item details:", err);
            }
        }
    };

    // Update item quantity
    const updateItemQuantity = (itemId: number, quantity: number) => {
        setInvoiceItems((prev) =>
            prev.map((item) =>
                item.item_id === itemId
                    ? {
                          ...item,
                          quantity: quantity || 0,
                          total_price: (quantity || 0) * item.unit_price,
                      }
                    : item,
            ),
        );
    };

    // Update item unit price
    const updateItemUnitPrice = (itemId: number, unitPrice: number) => {
        setInvoiceItems((prev) =>
            prev.map((item) =>
                item.item_id === itemId
                    ? {
                          ...item,
                          unit_price: unitPrice || 0,
                          total_price: item.quantity * (unitPrice || 0),
                      }
                    : item,
            ),
        );
    };

    // Remove item
    const removeItem = (itemId: number) => {
        setInvoiceItems((prev) =>
            prev.filter((item) => item.item_id !== itemId),
        );
        setSelectedItemIds((prev) => prev.filter((id) => id !== itemId));
    };

    // Calculate grand total
    const calculateGrandTotal = () => {
        return invoiceItems.reduce((sum, item) => sum + item.total_price, 0);
    };

    // Handle form submission
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (invoiceItems.length === 0) {
            setError("Please add at least one item");
            return;
        }

        // Validate all items have quantity and unit price
        const invalidItems = invoiceItems.filter(
            (item) => item.quantity <= 0 || item.unit_price < 0,
        );
        if (invalidItems.length > 0) {
            setError(
                "All items must have a quantity greater than 0 and a valid unit price",
            );
            return;
        }

        const grandTotal = calculateGrandTotal();

        const invoiceData = {
            grand_price: grandTotal,
            items: invoiceItems.map(
                (item): InvoiceItemInput => ({
                    item_id: item.item_id,
                    quantity: item.quantity,
                    unit_price: item.unit_price,
                }),
            ),
            tags: selectedTagIds,
        };

        setSubmitting(true);
        try {
            if (isEditMode && id) {
                await invoiceApi.update(parseInt(id), invoiceData);
            } else {
                await invoiceApi.create(invoiceData);
            }
            navigate("/invoices");
        } catch (err: any) {
            setError(err.message || "Failed to save invoice");
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <div className="container mx-auto py-8">
                <div className="text-center">Loading...</div>
            </div>
        );
    }

    return (
        <div className="container mx-auto py-8">
            <Card>
                <CardHeader>
                    <div className="flex justify-between items-center">
                        <CardTitle>
                            {isEditMode
                                ? `Edit Invoice #${id}`
                                : "Create Invoice"}
                        </CardTitle>
                        <Button
                            variant="outline"
                            onClick={() => navigate("/invoices")}
                        >
                            Cancel
                        </Button>
                    </div>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {error && (
                            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                                {error}
                            </div>
                        )}

                        {/* Items Selection */}
                        <div className="space-y-2">
                            <Label>Items *</Label>
                            <ItemMultiSelect
                                value={selectedItemIds}
                                onChange={handleItemsChange}
                                placeholder="Select items for this invoice..."
                                disabled={submitting}
                            />
                        </div>

                        {/* Items Table */}
                        {invoiceItems.length > 0 && (
                            <div className="border rounded-lg">
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>Item</TableHead>
                                            <TableHead className="w-32">
                                                Quantity
                                            </TableHead>
                                            <TableHead className="w-32">
                                                Unit Price
                                            </TableHead>
                                            <TableHead className="w-32">
                                                Total
                                            </TableHead>
                                            <TableHead className="w-16"></TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {invoiceItems.map((item) => (
                                            <TableRow key={item.item_id}>
                                                <TableCell className="font-medium">
                                                    {item.itemName ||
                                                        `Item #${item.item_id}`}
                                                </TableCell>
                                                <TableCell>
                                                    <Input
                                                        type="number"
                                                        min="1"
                                                        value={item.quantity}
                                                        onChange={(e) =>
                                                            updateItemQuantity(
                                                                item.item_id,
                                                                parseInt(
                                                                    e.target
                                                                        .value,
                                                                ) || 0,
                                                            )
                                                        }
                                                        disabled={submitting}
                                                        required
                                                    />
                                                </TableCell>
                                                <TableCell>
                                                    <Input
                                                        type="number"
                                                        min="0"
                                                        step="0.01"
                                                        value={item.unit_price}
                                                        onChange={(e) =>
                                                            updateItemUnitPrice(
                                                                item.item_id,
                                                                parseFloat(
                                                                    e.target
                                                                        .value,
                                                                ) || 0,
                                                            )
                                                        }
                                                        disabled={submitting}
                                                        required
                                                    />
                                                </TableCell>
                                                <TableCell>
                                                    $
                                                    {item.total_price.toFixed(
                                                        2,
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    <Button
                                                        type="button"
                                                        variant="ghost"
                                                        size="sm"
                                                        onClick={() =>
                                                            removeItem(
                                                                item.item_id,
                                                            )
                                                        }
                                                        disabled={submitting}
                                                    >
                                                        <Trash2 className="h-4 w-4" />
                                                    </Button>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                        <TableRow>
                                            <TableCell
                                                colSpan={3}
                                                className="text-right font-bold"
                                            >
                                                Grand Total:
                                            </TableCell>
                                            <TableCell className="font-bold">
                                                $
                                                {calculateGrandTotal().toFixed(
                                                    2,
                                                )}
                                            </TableCell>
                                            <TableCell></TableCell>
                                        </TableRow>
                                    </TableBody>
                                </Table>
                            </div>
                        )}

                        {/* Tags Selection */}
                        <div className="space-y-2">
                            <Label>Tags</Label>
                            <TagMultiSelect
                                value={selectedTagIds}
                                onChange={setSelectedTagIds}
                                placeholder="Select tags for this invoice..."
                                disabled={submitting}
                            />
                        </div>

                        {/* Submit Button */}
                        <div className="flex justify-end gap-2 pt-4">
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => navigate("/invoices")}
                                disabled={submitting}
                            >
                                Cancel
                            </Button>
                            <Button type="submit" disabled={submitting}>
                                {submitting
                                    ? "Saving..."
                                    : isEditMode
                                      ? "Update Invoice"
                                      : "Create Invoice"}
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
}
