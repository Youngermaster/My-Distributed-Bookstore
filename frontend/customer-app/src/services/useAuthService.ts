import { useAuthStore } from "@/store/authStore";

/**
 * Hook for authentication operations
 * Wraps authStore for cleaner API usage
 */
export function useAuthService() {
  const {
    user,
    isAuthenticated,
    isLoading,
    error,
    login,
    register,
    logout,
    loadUser,
    clearError,
  } = useAuthStore();

  return {
    // State
    user,
    isAuthenticated,
    isLoading,
    error,

    // Check if user has specific role
    hasRole: (roleName: string) => {
      return user?.role?.name === roleName || false;
    },

    // Check if user is admin
    isAdmin: () => {
      return user?.role?.name === "admin" || false;
    },

    // Check if user has permission
    hasPermission: (permission: string) => {
      return user?.role?.permissions?.includes(permission) ?? false;
    },

    // Actions
    login,
    register,
    logout,
    loadUser,
    clearError,
  };
}
