import * as React from "react";

import { cn } from "@/lib/utils";

export interface SeparatorProps extends React.HTMLAttributes<HTMLDivElement> {
  orientation?: "horizontal" | "vertical";
  decorative?: boolean;
}

const Separator = React.forwardRef<HTMLDivElement, SeparatorProps>(
  (
    {
      className,
      orientation = "horizontal",
      decorative = true,
      role,
      ...props
    },
    ref
  ) => {
    const baseRole = decorative ? "none" : role ?? "separator";
    return (
      <div
        ref={ref}
        role={baseRole}
        aria-orientation={orientation === "vertical" ? "vertical" : undefined}
        className={cn(
          "shrink-0 bg-border",
          orientation === "vertical" ? "w-px" : "h-px w-full",
          className
        )}
        {...props}
      />
    );
  }
);
Separator.displayName = "Separator";

export { Separator };
