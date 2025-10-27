import { createFileRoute } from "@tanstack/react-router";
import Genres from "@/pages/Genres";

export const Route = createFileRoute("/genres")({
  component: Genres,
});
