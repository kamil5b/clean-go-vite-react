import { createBrowserRouter } from "react-router-dom";
import {
    HomePage,
    LoginPage,
    RegisterPage,
    NotFoundPage,
    ItemsPage,
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
                path: "/items",
                element: (
                    <ProtectedRoute>
                        <ItemsPage />
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
