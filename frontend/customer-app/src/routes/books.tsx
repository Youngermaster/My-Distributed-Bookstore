import {
  createFileRoute,
  Outlet,
  useRouterState,
} from "@tanstack/react-router";
import BookList from "@/pages/BookList";

function BooksRouteComponent() {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const isDetailRoute = pathname !== "/books" && pathname.startsWith("/books/");

  if (isDetailRoute) {
    return <Outlet />;
  }

  return <BookList />;
}

export const Route = createFileRoute("/books")({
  component: BooksRouteComponent,
});
