import React from 'react'

export default ({ title, children }) => (
  <div className="lg:py-4 mb-10 lg:w-2/4">
    <h2 className="text-xl font-mono font-bold">Something went wrong.</h2>
    <div className="text-lg mt-4">{title}</div>
    <div className="text-gray-700 text-sm">{children}</div>
  </div>
)
