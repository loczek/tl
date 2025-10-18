import React from "react";
import { twMerge } from "tailwind-merge";

interface Props extends React.ComponentPropsWithoutRef<"button"> {
  as?: React.ElementType<any>;
}

function Button({
  children,
  className,
  as: Component = "button",
  ...props
}: Props) {
  return (
    <Component
      type="button"
      {...props}
      className={twMerge(
        "bg-dark-700 px-6 py-4 rounded-xl transition-transform active:scale-95 cursor-pointer outline-dark-600 focus:outline-2 outline-offset-2",
        className,
      )}
    >
      {children}
    </Component>
  );
}

export default Button;
