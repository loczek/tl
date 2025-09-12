import { AnimatePresence } from "framer-motion";
import { type FC } from "react";

interface Props {
  children: React.ReactNode;
}

const ModalContainer: FC<Props> = ({ children }) => (
  <AnimatePresence initial={false} mode="wait">
    {children}
  </AnimatePresence>
);

export default ModalContainer;
