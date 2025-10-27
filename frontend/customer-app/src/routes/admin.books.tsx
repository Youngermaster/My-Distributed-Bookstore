import { createFileRoute, redirect } from "@tanstack/react-router";
import ManageBooks from "@/pages/admin/ManageBooks";
import { useAuthStore } from "@/store/authStore";

export const Route = createFileRoute("/admin/books")({
  beforeLoad: () => {
    const { isAuthenticated, user } = useAuthStore.getState();
    if (!isAuthenticated) {
      throw redirect({ to: "/login" });
    }
    if (user?.role?.name !== "admin") {
      throw redirect({ to: "/" });
    }
  },
  component: ManageBooks,
});
