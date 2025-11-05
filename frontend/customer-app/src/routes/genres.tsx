import {
  createFileRoute,
  Outlet,
  useRouterState,
} from "@tanstack/react-router";
import Genres from "@/pages/Genres";

function GenresRouteComponent() {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const isDetailRoute =
    pathname !== "/genres" && pathname.startsWith("/genres/");

  if (isDetailRoute) {
    return <Outlet />;
  }

  return <Genres />;
}

export const Route = createFileRoute("/genres")({
  component: GenresRouteComponent,
});
