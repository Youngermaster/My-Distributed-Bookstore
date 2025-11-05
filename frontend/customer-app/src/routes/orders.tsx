import { createFileRoute } from "@tanstack/react-router";
import Orders from "@/pages/Orders";

export const Route = createFileRoute("/orders")({
  component: Orders,
});
