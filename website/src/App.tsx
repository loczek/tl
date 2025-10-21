import { useState, useTransition } from "react";
import * as z from "zod/mini";
import Button from "./components/Button";
import Footer from "./components/Footer";
import Header from "./components/Header";
import Loader from "./components/Loader";
import { shortPrefix } from "./constants/http";
import { writeToClipboard } from "./utils/clipboard";
import { postData } from "./utils/fetch";
import { joinPaths } from "./utils/join";
import { sleep } from "./utils/sleep";

const apiAddResponseSchema = z.object({ short_code: z.string() });

function App() {
  const [url, setUrl] = useState("");
  const [shortCode, setShortCode] = useState("");
  const [isLoading, startTransition] = useTransition();

  const handleAction = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!url) {
      return;
    }

    try {
      startTransition(async () => {
        await sleep(300);

        const data = await postData(
          shortPrefix + "/api/add",
          { url },
          apiAddResponseSchema,
        );

        setShortCode(data.short_code);
      });
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
        <div className="anim-from-bottom flex flex-col items-center">
          <h1 className="text-6xl tracking-tight font-bold mb-8 text-white z-10">
            Simple Link Shortener
          </h1>
          <form
            onSubmit={handleAction}
            className="gap-4 *:outline-2 *:outline-offset-2 *:focus-visible:outline-[#303030] *:outline-transparent font-semibold z-10 flex"
          >
            {shortCode ? (
              <>
                <div className="px-6 py-4 bg-[#1C1C1C] rounded-xl w-96 inline-block">
                  {joinPaths(shortPrefix, shortCode)}
                </div>
                <Button
                  type="button"
                  onClick={() =>
                    void writeToClipboard(joinPaths(shortPrefix, shortCode))
                  }
                >
                  Copy
                </Button>
              </>
            ) : (
              <>
                <input
                  className="px-6 py-4 bg-[#1C1C1C] rounded-xl w-96 outline-1! outline-dark-600! outline-offset-0! focus:outline-dark-500"
                  onChange={(e) => setUrl(e.target.value)}
                  value={url}
                  placeholder="Enter url here"
                  maxLength={2048}
                  name="url"
                  type="url"
                  autoFocus
                  autoComplete="off"
                  required
                />
                <Button type="submit">
                  {isLoading ? <Loader /> : "Shorten"}
                </Button>
              </>
            )}
          </form>
        </div>
      </div>

      <Footer />
    </div>
  );
}

export default App;
