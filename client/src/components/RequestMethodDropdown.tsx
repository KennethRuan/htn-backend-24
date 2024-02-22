import React, { Fragment } from "react";
import { Listbox, Transition } from "@headlessui/react";

type RequestMethodDropdownProps = {
  value: string;
  onChange: (value: string) => void;
  className?: string;
};

const RequestMethodDropdown = ({
  value,
  onChange,
  className,
}: RequestMethodDropdownProps) => {
  const methods = ["GET", "POST", "PUT", "DELETE"];
  const requestColourMap: { [method: string]: string } = {
    GET: "text-green-600",
    POST: "text-amber-200",
    PUT: "text-blue-400",
    DELETE: "text-rose-400",
  };

  return (
    <Listbox value={value} onChange={onChange}>
      <div className="relative w-full h-full">
        <Listbox.Button className="w-full h-full">
          <div className="w-full h-full flex flex-col items-center justify-center p-2">
            <p className={`${requestColourMap[value]} font-bold leading-none`}>
              {value}
            </p>
          </div>
        </Listbox.Button>
        <Transition
          as={Fragment}
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <Listbox.Options className="absolute max-h-60 w-full mt-4 p-4 rounded bg-stone-900 border-2 border-stone-700 shadow-lg shadow-black focus:outline-none">
            {methods.map((method, index) => (
              <Listbox.Option key={index} value={method}>
                <div className="py-2 px-4 flex items-center justify-start cursor-pointer">
                  <p
                    className={`${requestColourMap[method]} font-bold leading-none`}
                  >
                    {method}
                  </p>
                </div>
              </Listbox.Option>
            ))}
          </Listbox.Options>
        </Transition>
      </div>
    </Listbox>
  );
};

export default RequestMethodDropdown;
