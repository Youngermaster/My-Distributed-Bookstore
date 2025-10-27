import { createFileRoute } from "@tanstack/react-router";
import AuthorDetail from "@/pages/AuthorDetail";

export const Route = createFileRoute("/authors/$id")({
  component: AuthorDetail,
});
