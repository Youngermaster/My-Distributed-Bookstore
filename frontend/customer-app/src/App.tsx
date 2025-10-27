import { BrowserRouter, Routes, Route, Navigate, Link } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useEffect } from "react";
import { useAuthStore } from "./store/authStore";
import { Button } from "./components/ui/button";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Home from "./pages/Home";
import BookList from "./pages/BookList";
import BookDetail from "./pages/BookDetail";
import Genres from "./pages/Genres";
import GenreDetail from "./pages/GenreDetail";
import AuthorDetail from "./pages/AuthorDetail";
import Wishlist from "./pages/Wishlist";
import ManageBooks from "./pages/admin/ManageBooks";
import { Heart, BookOpen, LogOut, User, Shield } from "lucide-react";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function ProtectedRoute({
  children,
  requireAdmin = false,
}: {
  children: React.ReactNode;
  requireAdmin?: boolean;
}) {
  const { isAuthenticated, user } = useAuthStore();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requireAdmin && user?.role?.name !== "admin") {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

function Navigation() {
  const { isAuthenticated, user, logout } = useAuthStore();
  const isAdmin = user?.role?.name === "admin";

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center gap-8">
            <Link
              to="/"
              className="flex items-center gap-2 text-xl font-bold text-gray-900 hover:text-gray-700"
            >
              <BookOpen className="h-6 w-6" />
              Bookstore
            </Link>
            <div className="hidden md:flex items-center gap-4">
              <Link to="/books" className="text-gray-700 hover:text-gray-900 font-medium">
                Books
              </Link>
              <Link to="/genres" className="text-gray-700 hover:text-gray-900 font-medium">
                Genres
              </Link>
            </div>
          </div>

          <div className="flex items-center gap-2">
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
                <div className="flex items-center gap-2 px-3 py-1 bg-gray-100 rounded-md">
                  <User className="h-4 w-4 text-gray-600" />
                  <span className="text-sm font-medium text-gray-700">
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

function AppContent() {
  const { loadUser, isAuthenticated } = useAuthStore();

  useEffect(() => {
    if (isAuthenticated) {
      loadUser();
    }
  }, [isAuthenticated, loadUser]);

  return (
    <>
      <Navigation />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/books" element={<BookList />} />
        <Route path="/books/:id" element={<BookDetail />} />
        <Route path="/genres" element={<Genres />} />
        <Route path="/genres/:slug" element={<GenreDetail />} />
        <Route path="/authors/:id" element={<AuthorDetail />} />
        <Route
          path="/wishlist"
          element={
            <ProtectedRoute>
              <Wishlist />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/books"
          element={
            <ProtectedRoute requireAdmin>
              <ManageBooks />
            </ProtectedRoute>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AppContent />
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
