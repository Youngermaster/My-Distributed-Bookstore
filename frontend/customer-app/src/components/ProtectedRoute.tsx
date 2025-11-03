import { useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useAuthService } from "@/services";
import { toast } from "@/lib/toast";

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireAuth?: boolean;
  requireRole?: string;
  requirePermission?: string;
}

/**
 * Protected Route Component
 * Handles authentication and authorization checks
 */
export function ProtectedRoute({
  children,
  requireAuth = true,
  requireRole,
  requirePermission,
}: ProtectedRouteProps) {
  const navigate = useNavigate();
  const {
    isAuthenticated,
    user,
    hasRole,
    hasPermission: checkPermission,
    isLoading,
  } = useAuthService();

  useEffect(() => {
    if (isLoading) return;

    // Check authentication
    if (requireAuth && !isAuthenticated) {
      toast.error("Please login to access this page");
      navigate({ to: "/login" });
      return;
    }

    // Check role
    if (requireRole && !hasRole(requireRole)) {
      toast.error("You don't have permission to access this page");
      navigate({ to: "/" });
      return;
    }

    // Check specific permission
    if (requirePermission && !checkPermission(requirePermission)) {
      toast.error("You don't have permission to access this page");
      navigate({ to: "/" });
      return;
    }
  }, [
    isAuthenticated,
    user,
    isLoading,
    requireAuth,
    requireRole,
    requirePermission,
  ]);

  // Show loading state
  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  // Don't render children until auth check is complete
  if (requireAuth && !isAuthenticated) {
    return null;
  }

  if (requireRole && !hasRole(requireRole)) {
    return null;
  }

  if (requirePermission && !checkPermission(requirePermission)) {
    return null;
  }

  return <>{children}</>;
}
