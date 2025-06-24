import logo from "../assets/logo.svg";
import github from "../assets/icons/github.svg";

function Header() {
  return (
    <header className="w-full py-4 px-8 rounded-2xl bg-dark-900 flex justify-between">
      <img src={logo} alt="logo" className="size-10" />
      <a
        href="https://github.com/42fm/42fm"
        className="p-2 bg-dark-700 rounded-lg hover:bg-dark-600 transition-colors "
      >
        <img src={github} alt="github repo" className="size-6" />
      </a>
    </header>
  );
}

export default Header;
