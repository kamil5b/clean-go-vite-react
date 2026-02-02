import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import {
    NavigationMenu,
    NavigationMenuList,
    NavigationMenuItem,
    NavigationMenuLink,
} from "@/components/ui/navigation-menu";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function NavBar() {
    const location = useLocation();
    const navigate = useNavigate();
    const { user, isAuthenticated, logout } = useAuth();

    const navItems = [
        { label: "Home", href: "/" },
        ...(isAuthenticated ? [{ label: "Items", href: "/items" }] : []),
    ];

    const isActive = (href: string) => location.pathname === href;

    const handleLogout = async () => {
        try {
            await logout();
            navigate("/login");
        } catch (error) {
            console.error("Logout failed:", error);
        }
    };

    return (
        <nav className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center">
                        <span className="text-xl font-bold">App</span>
                    </div>

                    <div className="flex items-center gap-4">
                        <NavigationMenu>
                            <NavigationMenuList>
                                {navItems.map((item) => (
                                    <NavigationMenuItem key={item.href}>
                                        <NavigationMenuLink
                                            href={item.href}
                                            className={cn(
                                                "inline-flex items-center justify-center px-4 py-2 rounded-md text-sm font-medium transition-colors",
                                                isActive(item.href)
                                                    ? "bg-accent text-accent-foreground"
                                                    : "hover:bg-accent hover:text-accent-foreground",
                                            )}
                                        >
                                            {item.label}
                                        </NavigationMenuLink>
                                    </NavigationMenuItem>
                                ))}
                            </NavigationMenuList>
                        </NavigationMenu>

                        {isAuthenticated && user ? (
                            <div className="flex items-center gap-4 ml-4 pl-4 border-l">
                                <div className="text-sm">
                                    <p className="font-medium">{user.name}</p>
                                    <p className="text-xs text-gray-500">
                                        {user.email}
                                    </p>
                                </div>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={handleLogout}
                                >
                                    Logout
                                </Button>
                            </div>
                        ) : (
                            <div className="flex gap-2 ml-4 pl-4 border-l">
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => navigate("/login")}
                                >
                                    Login
                                </Button>
                                <Button
                                    size="sm"
                                    onClick={() => navigate("/register")}
                                >
                                    Register
                                </Button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </nav>
    );
}
