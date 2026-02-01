import { useEffect, useState } from "react";
import reactLogo from "../assets/react.svg";
import viteLogo from "/vite.svg";
import { messageApi } from "@/api/message";

export function HomePage() {
    const [messageFromServer, setMessageFromServer] = useState("");

    useEffect(() => {
        const fetchData = async () => {
            const data = await messageApi.getMessage();
            setMessageFromServer(data?.content || "No message received");
        };

        fetchData().catch((e) => console.error(e));
    }, []);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen max-w-5xl mx-auto px-8 py-8 text-center">
            <div className="flex gap-8 mb-8">
                <a
                    href="https://vitejs.dev"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="transition-all duration-300 hover:drop-shadow-lg hover:drop-shadow-[0_0_2em_#646cffaa]"
                >
                    <img src={viteLogo} className="h-24 p-6" alt="Vite logo" />
                </a>
                <a
                    href="https://react.dev"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="transition-all duration-300 hover:drop-shadow-lg hover:drop-shadow-[0_0_2em_#61dafbaa] animate-spin-slow"
                >
                    <img
                        src={reactLogo}
                        className="h-24 p-6"
                        alt="React logo"
                    />
                </a>
            </div>
            <h1 className="text-5xl font-bold mb-4">Golang + Vite + React</h1>
            <h2 className="text-2xl text-gray-400">{messageFromServer}</h2>
        </div>
    );
}
