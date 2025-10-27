import { createFileRoute } from "@tanstack/react-router";
import BookList from "@/pages/BookList";

export const Route = createFileRoute("/books")({
  component: BookList,
});
