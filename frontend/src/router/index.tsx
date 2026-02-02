import { createBrowserRouter } from "react-router-dom";
import {
    HomePage,
    CounterPage,
    LoginPage,
    RegisterPage,
    NotFoundPage,
    ItemsPage,
    TagsPage,
    InvoicesPage,
    InvoiceFormPage,
    InvoiceDetailPage,
} from "@/pages";
import { RootLayout } from "@/layouts/RootLayout";
import { ProtectedRoute } from "@/components/ProtectedRoute";

const router = createBrowserRouter([
    {
        element: <RootLayout />,
        children: [
            {
                path: "/",
                element: <HomePage />,
            },
            {
                path: "/counter",
                element: (
                    <ProtectedRoute>
                        <CounterPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/items",
                element: (
                    <ProtectedRoute>
                        <ItemsPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/tags",
                element: (
                    <ProtectedRoute>
                        <TagsPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/invoices",
                element: (
                    <ProtectedRoute>
                        <InvoicesPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/invoices/create",
                element: (
                    <ProtectedRoute>
                        <InvoiceFormPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/invoices/:id",
                element: (
                    <ProtectedRoute>
                        <InvoiceDetailPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/invoices/:id/edit",
                element: (
                    <ProtectedRoute>
                        <InvoiceFormPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/login",
                element: <LoginPage />,
            },
            {
                path: "/register",
                element: <RegisterPage />,
            },
            {
                path: "*",
                element: <NotFoundPage />,
            },
        ],
    },
]);

export default router;
