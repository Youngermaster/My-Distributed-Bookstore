import { createFileRoute, redirect } from "@tanstack/react-router";
import Wishlist from "@/pages/Wishlist";
import { useAuthStore } from "@/store/authStore";

export const Route = createFileRoute("/wishlist")({
  beforeLoad: () => {
    const { isAuthenticated } = useAuthStore.getState();
    if (!isAuthenticated) {
      throw redirect({ to: "/login" });
    }
  },
  component: Wishlist,
});
