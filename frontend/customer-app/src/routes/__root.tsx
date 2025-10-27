import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useAuthStore } from "@/store/authStore";
import { useThemeStore } from "@/store/themeStore";
import { Button } from "@/components/ui/button";
import ThemeToggle from "@/components/ThemeToggle";
import { Heart, BookOpen, LogOut, User, Shield } from "lucide-react";
import { useEffect } from "react";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function Navigation() {
  const { isAuthenticated, user, logout } = useAuthStore();
  const isAdmin = user?.role?.name === "admin";

  return (
    <nav className="bg-background border-b border-border sticky top-0 z-50">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center gap-8">
            <Link
              to="/"
              className="flex items-center gap-2 text-xl font-bold text-foreground hover:text-foreground/80"
            >
              <BookOpen className="h-6 w-6" />
              Ohara Bookstore
            </Link>
            <div className="hidden md:flex items-center gap-4">
              <Link
                to="/books"
                className="text-foreground/70 hover:text-foreground font-medium"
              >
                Books
              </Link>
              <Link
                to="/genres"
                className="text-foreground/70 hover:text-foreground font-medium"
              >
                Genres
              </Link>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <ThemeToggle />
            {isAuthenticated ? (
              <>
                <Link to="/wishlist">
                  <Button variant="ghost" size="sm">
                    <Heart className="h-4 w-4 mr-2" />
                    Wishlist
                  </Button>
                </Link>
                {isAdmin && (
                  <Link to="/admin/books">
                    <Button variant="ghost" size="sm">
                      <Shield className="h-4 w-4 mr-2" />
                      Manage Books
                    </Button>
                  </Link>
                )}
                <div className="flex items-center gap-2 px-3 py-1 bg-muted rounded-md">
                  <User className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium text-foreground">
                    {user?.full_name}
                  </span>
                </div>
                <Button variant="outline" size="sm" onClick={() => logout()}>
                  <LogOut className="h-4 w-4 mr-2" />
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="ghost" size="sm">
                    Login
                  </Button>
                </Link>
                <Link to="/register">
                  <Button size="sm">Register</Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}

function RootComponent() {
  const { loadUser, isAuthenticated } = useAuthStore();
  const { theme } = useThemeStore();

  useEffect(() => {
    if (isAuthenticated) {
      loadUser();
    }
  }, [isAuthenticated, loadUser]);

  // Apply dark class to document element
  useEffect(() => {
    const root = document.documentElement;
    if (theme === "dark") {
      root.classList.add("dark");
    } else {
      root.classList.remove("dark");
    }
  }, [theme]);

  return (
    <QueryClientProvider client={queryClient}>
      <Navigation />
      <Outlet />
    </QueryClientProvider>
  );
}

export const Route = createRootRoute({
  component: RootComponent,
});
