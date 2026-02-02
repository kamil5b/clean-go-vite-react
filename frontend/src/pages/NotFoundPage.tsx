import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";

export function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <div className="text-center space-y-6">
        <div className="space-y-2">
          <h1 className="text-9xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-400">
            404
          </h1>
          <h2 className="text-4xl font-bold text-white">Page Not Found</h2>
        </div>

        <p className="text-xl text-slate-300 max-w-md">
          Sorry, the page you're looking for doesn't exist. It might have been
          moved or deleted.
        </p>

        <div className="flex gap-4 justify-center pt-4">
          <Button
            onClick={() => navigate("/")}
            className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-2"
          >
            Go Home
          </Button>
          <Button
            onClick={() => navigate(-1)}
            variant="outline"
            className="text-white border-slate-500 hover:bg-slate-700 px-8 py-2"
          >
            Go Back
          </Button>
        </div>
      </div>
    </div>
  );
}
