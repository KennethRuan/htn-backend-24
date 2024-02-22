import React, { useState } from "react";
import RequestMethodDropdown from "./components/RequestMethodDropdown";
import RequestEndpointDropdown from "./components/RequestEndpointDropdown";
import { ChevronRightIcon } from "@heroicons/react/20/solid";

function App() {
  const [requestMethod, setRequestMethod] = useState("GET");
  const [requestEndpoint, setRequestEndpoint] = useState("");

  const API_BASE = "http://localhost:3000";

  const fetchRequest = async () => {
    const adr = `${API_BASE}${requestEndpoint}`;
    console.log(adr);
    try {
      const response = await fetch(adr, {
        method: requestMethod,
      });
      const data = await response.json();
      console.log(data);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <div className="relative w-screen h-screen bg-stone-900 font-mono">
          <div className="absolute left-0 top-0 w-full h-full opacity-[0.03] bg-[linear-gradient(to_right,gray_1px,transparent_1px),linear-gradient(gray_1px,transparent_1px)] bg-[length:16px_16px] z-[0]" />
          <div className="relative w-full h-full pt-48 flex flex-col justify-start items-center z-[1]">
            <div className="mx-auto w-[1000px] h-12 rounded-full border-stone-700 border-2 flex flex-row p-2 shadow-lg shadow-black">
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
                onClick={(e) => fetchRequest()}
              />
            </div>
          </div>
        </div>
      </header>
    </div>
  );
}

export default App;
