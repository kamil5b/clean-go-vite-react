import { createBrowserRouter } from "react-router-dom";
import { HomePage, CounterPage } from "@/pages";
import { RootLayout } from "@/layouts/RootLayout";

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
                element: <CounterPage />,
            },
        ],
    },
]);

export default router;
