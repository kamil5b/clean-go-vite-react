import { createBrowserRouter } from "react-router-dom";
import { HomePage, CounterPage } from "@/pages";

const router = createBrowserRouter([
    {
        path: "/",
        element: <HomePage />,
    },
    {
        path: "/counter",
        element: <CounterPage />,
    },
]);

export default router;
