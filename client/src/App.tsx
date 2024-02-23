import React, { useState } from "react";
import RequestMethodDropdown from "./components/RequestMethodDropdown";
import RequestEndpointDropdown from "./components/RequestEndpointDropdown";
import { ChevronRightIcon } from "@heroicons/react/20/solid";

function App() {
  const [requestMethod, setRequestMethod] = useState("GET");
  const [requestEndpoint, setRequestEndpoint] = useState("");
  const [response, setResponse] = useState<any>(null);
  const [showBody, setShowBody] = useState(false);
  const [body, setBody] = useState<any>("");

  const fetchRequest = async () => {
    const adr = `${requestEndpoint}`;
    console.log(adr);
    try {
      const requestOptions: RequestInit = {
        method: requestMethod,
        headers: {
          Accept: "*/*",
          "Accept-Encoding": "gzip, deflate, br",
          Connection: "keep-alive",
          "Content-Type": "application/json",
        },
      };
      // Attach the body to the request, must not be GET request and must be a valid JSON
      if (requestMethod !== "GET") {
        try {
          const obj = JSON.parse(body);
          requestOptions.body = JSON.stringify(obj);
        } catch (error) {
          console.error(error);
        }
      }

      const response = await fetch(adr, requestOptions);
      try {
        const text = await response.text();
        // setResponse(text);
        const data = await JSON.parse(text);
        setResponse(data);
      } catch (error: any) {
        console.error(error);
        setResponse(error?.message || error?.toString() || error);
      }
    } catch (error: any) {
      setResponse(error);
      console.error(error);
    }
  };

  const handleSubmit = () => {
    try {
      const obj = JSON.parse(body);
      const pretty = JSON.stringify(obj, null, 2);
      setBody(pretty);
    } catch (error) {}
    fetchRequest();
    setShowBody(false);
  };

  return (
    <div className="App">
      <header className="App-header">
        <div className="relative w-screen h-screen bg-stone-900 font-mono">
          <div className="absolute left-0 top-0 w-full h-full opacity-[0.03] bg-[linear-gradient(to_right,gray_1px,transparent_1px),linear-gradient(gray_1px,transparent_1px)] bg-[length:16px_16px] z-[0]" />
          <div className="relative w-full h-full pt-24 flex flex-col justify-start items-center z-[1]">
            <div className="w-[1000px] flex flex-row justify-start items-center">
              <div
                className="rounded h-12 bg-stone-900 ring-stone-700 ring-2 flex justify-start shadow-lg shadow-black cursor-pointer select-none"
                onClick={() => setShowBody((p) => !p)}
              >
                <p
                  className={`text-stone-100 font-bold h-full w-32 flex items-center justify-center p-2 ${
                    !showBody && `bg-stone-700`
                  }`}
                >
                  Response
                </p>
                <p
                  className={`text-stone-100 font-bold h-full w-32 flex items-center justify-center p-2 ${
                    showBody && `bg-stone-700`
                  }`}
                >
                  Body
                </p>
              </div>
            </div>

            <div className="mt-4 mx-auto w-[1000px] h-12 rounded-full border-stone-700 border-2 flex flex-row p-2 shadow-lg shadow-black">
              <div className="w-32 h-full flex gap-2 items-center justify-center">
                <RequestMethodDropdown
                  value={requestMethod}
                  onChange={setRequestMethod}
                />
              </div>
              {/* Vertical Line */}
              <div className="w-1 h-full bg-stone-700" />
              {/* Endpoint Input */}
              <div className="relative flex flex-row flex-1">
                <div className="absolute w-[calc(100%-24px)] h-full z-[1]">
                  <input
                    className="w-full h-full bg-stone-900 text-stone-100 px-4 focus:outline-none"
                    value={requestEndpoint}
                    onChange={(e) => setRequestEndpoint(e.target.value)}
                  />
                </div>
                {/* Dropdown for Endpoint Selection */}
                <div className="absolute w-full h-full flex justify-end item-center z-[0]">
                  <RequestEndpointDropdown
                    value={requestEndpoint}
                    onChange={setRequestEndpoint}
                    method={requestMethod}
                  />
                </div>
              </div>
              <ChevronRightIcon
                className="w-12 h-full text-stone-400 cursor-pointer"
                onClick={(e) => handleSubmit()}
              />
            </div>

            <div className="w-[1000px] h-[400px] overflow-auto mt-4 p-4 rounded bg-stone-900 border-2 border-stone-700 shadow-lg shadow-black focus:outline-none">
              {showBody ? (
                <textarea
                  className="w-full h-[95%] bg-stone-900 text-stone-100 px-4 focus:outline-none"
                  value={body}
                  onChange={(e) => setBody(e.target.value)}
                />
              ) : (
                <pre className="text-stone-100">
                  {JSON.stringify(response, null, 2)}
                </pre>
              )}
            </div>
          </div>
        </div>
      </header>
    </div>
  );
}

export default App;
