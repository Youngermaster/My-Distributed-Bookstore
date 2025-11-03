import { createFileRoute, redirect } from "@tanstack/react-router";
import AdminDashboard from "@/pages/admin/Dashboard";
import { useAuthStore } from "@/store/authStore";

export const Route = createFileRoute("/admin/dashboard")({
  beforeLoad: () => {
    const { isAuthenticated, user } = useAuthStore.getState();
    if (!isAuthenticated) {
      throw redirect({ to: "/login" });
    }
    if (user?.role?.name !== "admin") {
      throw redirect({ to: "/" });
    }
  },
  component: AdminDashboard,
});
