import { motion, type Variants } from "motion/react";
import { type FC } from "react";
import { twMerge } from "tailwind-merge";
import Button from "./Button";
import ModalBackdrop from "./ModalBackdrop";

const variants: Variants = {
  hidden: {
    y: "var(--modal-y-from)",
    scale: "var(--modal-scale-from)",
    originY: "var(--modal-origin-y)",
    originX: "var(--modal-origin-x)",
    opacity: 0,
  },
  visible: {
    y: "var(--modal-y-to)",
    scale: "var(--modal-scale-to)",
    opacity: 1,
    transition: {
      duration: 0.1,
      type: "spring",
      damping: 25,
      stiffness: 500,
    },
  },
  exit: {
    y: "var(--modal-y-from)",
    scale: "var(--modal-scale-from)",
    opacity: 0,
    transition: {
      duration: 0.1,
      type: "spring",
      damping: 25,
      stiffness: 500,
    },
  },
};

interface IModel {
  handleClose: () => void;
  hasCloseButton?: boolean;
  className?: string;
  children: React.ReactNode;
}

const Modal: FC<IModel> = ({
  className,
  handleClose,
  hasCloseButton,
  children,
}) => {
  return (
    <ModalBackdrop onClick={handleClose}>
      <motion.div
        className={twMerge(
          "rounded-2xl p-6 [--modal-origin-x:0.5] [--modal-origin-y:0] [--modal-scale-from:1] [--modal-scale-to:1] [--modal-y-from:64px] [--modal-y-to:0px] bg-dark-900 sm:p-8 sm:[--modal-origin-y:0.5] sm:[--modal-scale-from:0.8] sm:[--modal-scale-to:1] sm:[--modal-y-from:0]",
          className
        )}
        variants={variants}
        initial="hidden"
        animate="visible"
        exit="exit"
        onClick={(e) => e.stopPropagation()}
      >
        {children}
        {hasCloseButton && (
          <div className="mt-4 flex justify-center">
            <Button onClick={handleClose}>Close</Button>
          </div>
        )}
      </motion.div>
    </ModalBackdrop>
  );
};

export default Modal;
