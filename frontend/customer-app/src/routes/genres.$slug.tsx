import { createFileRoute } from "@tanstack/react-router";
import GenreDetail from "@/pages/GenreDetail";

export const Route = createFileRoute("/genres/$slug")({
  component: GenreDetail,
});
