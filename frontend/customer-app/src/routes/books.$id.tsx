import { createFileRoute } from "@tanstack/react-router";
import BookDetail from "@/pages/BookDetail";

export const Route = createFileRoute("/books/$id")({
  component: BookDetail,
});
