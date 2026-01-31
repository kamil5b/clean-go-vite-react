import { createBrowserRouter } from "react-router-dom";
import { HomePage, CounterPage, LoginPage, RegisterPage } from "@/pages";
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
                path: "/login",
                element: <LoginPage />,
            },
            {
                path: "/register",
                element: <RegisterPage />,
            },
        ],
    },
]);

export default router;
