import React, { Fragment } from "react";
import { Listbox, Transition } from "@headlessui/react";
import { ChevronUpDownIcon } from "@heroicons/react/20/solid";

type RequestEndpointDropdownProps = {
  value: string;
  onChange: (value: string) => void;
  method: string;
  className?: string;
};

const RequestEndpointDropdown = ({
  value,
  onChange,
  method,
  className,
}: RequestEndpointDropdownProps) => {
  type RequestOption = {
    label: string;
    endpoint: string;
  };

  const requestEndpointMap: { [method: string]: RequestOption[] } = {
    GET: [
      { label: "Get All Users", endpoint: "/api/users" },
      { label: "Get User by ID", endpoint: "/api/users/:id" },
      {
        label: "Get Skills by Frequency",
        endpoint: "/api/skills/?min_frequency=15&max_frequency=20",
      },
    ],
    POST: [],
    PUT: [{ label: "Update User by ID", endpoint: "/api/users/:id" }],
    DELETE: [],
  };

  return (
    <Listbox value={value} onChange={onChange}>
      <div className="relative w-full h-full">
        <div className="w-full h-full flex flex-row justify-end items-center">
          <Listbox.Button className="w-6 h-6 cursor-pointer">
            <ChevronUpDownIcon className="w-full h-full text-stone-400" />
          </Listbox.Button>
        </div>
        <Transition
          as={Fragment}
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <Listbox.Options className="absolute max-h-60 w-full mt-4 p-4 rounded bg-stone-900 border-2 border-stone-700 shadow-lg shadow-black focus:outline-none">
            {requestEndpointMap[method].map(({ label, endpoint }, index) => (
              <Listbox.Option key={index} value={endpoint}>
                <div className="py-2 px-4 flex items-center justify-start cursor-pointer">
                  <p className={`text-stone-100 font-bold leading-none`}>
                    {label}
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

export default RequestEndpointDropdown;
