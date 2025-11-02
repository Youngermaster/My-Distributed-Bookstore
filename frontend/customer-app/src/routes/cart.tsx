import { createFileRoute } from "@tanstack/react-router";
import Cart from "@/pages/Cart";

export const Route = createFileRoute("/cart")({
  component: Cart,
});
