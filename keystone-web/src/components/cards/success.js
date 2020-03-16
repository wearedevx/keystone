import React from 'react'

export default ({ title, children }) => (
  <div className="lg:py-4 mb-10 lg:w-2/4">
    <h2 className="text-xl font-mono font-bold">{title}</h2>
    <div className="text-lg mt-4">{children}</div>
  </div>
)
