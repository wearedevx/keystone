import React from 'react'

export default ({ title, children }) => (
  <div className="">
    <h2 className="text-xl font-mono font-bold">Something went wrong.</h2>
    <div className="text-lg mt-4">{title}</div>
    <div className="text-gray-700 text-sm">{children}</div>
  </div>
)
