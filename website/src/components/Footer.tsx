function Footer() {
  return (
    <footer className="w-full py-4 px-8 rounded-2xl bg-[#141414] flex justify-between items-center">
      <span className="py-2 block text-center">
        Made by <span className="opacity-60">loczek</span>
      </span>
      {import.meta.env.PROD ? "built and served (prod)" : "not built (dev)"}
    </footer>
  );
}

export default Footer;
