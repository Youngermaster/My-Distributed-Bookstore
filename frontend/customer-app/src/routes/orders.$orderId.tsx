import { createFileRoute } from "@tanstack/react-router";
import OrderDetail from "@/pages/OrderDetail";

export const Route = createFileRoute("/orders/$orderId")({
  component: OrderDetail,
});
