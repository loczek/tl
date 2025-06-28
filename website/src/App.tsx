import { useState } from "react";
import Footer from "./components/Footer";
import Header from "./components/Header";
import Loader from "./components/Loader";

function App() {
  const [url, setUrl] = useState("");
  const [shortCode, setShortCode] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleAction = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!url) {
      console.error("ERROR: missing form key 'url'");
      return;
    }

    try {
      setIsLoading(true);
      const res = await fetch("http://localhost:3000/api/add", {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url }),
      });

      const data = await res.json();

      setShortCode(data.short_code);
    } catch (e) {
      console.error(e);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container h-screen mx-auto space-y-4 py-8">
      <Header />
      <div className="relative h-4/5 flex items-center justify-center flex-col rounded-2xl bg-[#141414]">
        <div className="absolute inset-0 border-ring blur-xs brightness-75" />
        <div className="absolute inset-0 border-ring-glow brightness-75" />
        {/* <img src={bglow} className="mix-blend-lighten absolute inset-0" /> */}
        {/* <div className="absolute inset-0 outer rounded-2xl mix-blend-lighten">
          <div className="temp h-full mix-blend-color-dodge"></div> */}
        {/* <div className="absolute inset-0 outer rounded-2xl mix-blend-plus-lighter"> */}
        {/* <div className="absolute inset-0 rounded-2xl border-8 border-[#808080] blur-lg mix-blend-color-dodge"></div> */}
        {/* </div> */}
        <h1 className="text-6xl tracking-tight font-bold my-8 text-white z-10 anim-from-bottom">
          Simple Link Shortener
        </h1>
        <form
          onSubmit={handleAction}
          // className="space-x-4 *:focus-visible:outline-2 *:outline-offset-2 *:focus-visible:outline-zinc-600 font-semibold"
          className="gap-4 *:focus-visible:outline-2 *:outline-offset-2 *:focus-visible:outline-[#303030] font-semibold z-10 flex"
        >
          {shortCode ? (
            <div
              onClick={() => {
                navigator.clipboard.writeText(
                  "http://localhost:3000/" + shortCode
                );
              }}
              className="px-6 py-4 bg-[#1C1C1C] rounded-md w-96 inline-block"
            >
              http://localhost:3000/{shortCode}
            </div>
          ) : (
            <input
              className="px-6 py-4 bg-[#1C1C1C] rounded-md w-96 anim-from-bottom"
              onChange={(e) => setUrl(e.target.value)}
              value={url}
              placeholder="Enter url here"
              maxLength={2048}
              name="url"
              type="url"
              autoFocus
              autoComplete="off"
            />
          )}
          <button
            type="submit"
            className="bg-[#242424] px-6 py-4 rounded-md anim-from-bottom"
          >
            {isLoading ? <Loader /> : "Shorten"}
          </button>
        </form>
      </div>
      <Footer />
    </div>
  );
}

export default App;
