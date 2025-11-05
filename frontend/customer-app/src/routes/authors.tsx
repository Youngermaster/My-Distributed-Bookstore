import {
  createFileRoute,
  Outlet,
  useRouterState,
} from "@tanstack/react-router";
import Authors from "@/pages/Authors";

function AuthorsRouteComponent() {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const isDetailRoute =
    pathname !== "/authors" && pathname.startsWith("/authors/");

  if (isDetailRoute) {
    return <Outlet />;
  }

  return <Authors />;
}

export const Route = createFileRoute("/authors")({
  component: AuthorsRouteComponent,
});
