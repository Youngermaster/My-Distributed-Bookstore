import { createFileRoute } from "@tanstack/react-router";
import PublisherDetail from "@/pages/PublisherDetail";

export const Route = createFileRoute("/publishers/$id")({
  component: PublisherDetail,
});
