import { useState } from "react";
import Header from "./components/Header";
import Footer from "./components/Footer";

function App() {
  const [shortCode, setShortCode] = useState<string>("");

  const handleAction = async (formData: FormData) => {
    const url = formData.get("url");

    if (!url) {
      console.error("ERROR: missing form key 'url'");
      return;
    }

    try {
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
        <h1 className="text-6xl tracking-tight font-bold my-8 text-white z-10">
          Simple Link Shortener
        </h1>
        <form
          action={handleAction}
          // className="space-x-4 *:focus-visible:outline-2 *:outline-offset-2 *:focus-visible:outline-zinc-600 font-semibold"
          className="space-x-4 *:focus-visible:outline-2 *:outline-offset-2 *:focus-visible:outline-[#303030] font-semibold z-10"
        >
          {shortCode ? (
            <div
              onClick={() => {
                navigator.clipboard.writeText(
                  "http://localhost:3000/" + shortCode,
                );
              }}
              className="px-6 py-4 bg-[#1C1C1C] rounded-md w-96 inline-block"
            >
              http://localhost:3000/{shortCode}
            </div>
          ) : (
            <input
              placeholder="Enter url here"
              maxLength={2048}
              name="url"
              type="url"
              className="px-6 py-4 bg-[#1C1C1C] rounded-md w-96"
              autoFocus
            />
          )}
          <button className="bg-[#242424] px-6 py-4 rounded-md">Shorten</button>
        </form>
      </div>
      <Footer />
    </div>
  );
}

export default App;
