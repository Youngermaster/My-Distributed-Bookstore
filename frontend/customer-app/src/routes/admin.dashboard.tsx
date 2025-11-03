import { createFileRoute } from "@tanstack/react-router";
import AdminDashboard from "@/pages/admin/Dashboard";

export const Route = createFileRoute("/admin/dashboard")({
  component: AdminDashboard,
});
