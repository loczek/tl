import { motion } from "framer-motion";
import { type MouseEvent } from "react";

interface Props {
  children: React.ReactNode;
  onClick: (event: MouseEvent<HTMLDivElement> | undefined) => void;
}

const ModalBackdrop: React.FC<Props> = ({ children, onClick }) => {
  return (
    <motion.div
      className="fixed inset-0 z-[200] h-full flex items-end justify-center bg-black/60 sm:items-center"
      onClick={onClick}
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
    >
      {children}
    </motion.div>
  );
};

export default ModalBackdrop;
