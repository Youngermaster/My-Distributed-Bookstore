import { createFileRoute } from "@tanstack/react-router";
import Publishers from "@/pages/Publishers";

export const Route = createFileRoute("/publishers")({
  component: Publishers,
});
