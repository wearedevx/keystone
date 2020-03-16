import React from 'react'

export default ({ title, children }) => (
  <div className="">
    <h2 className="text-xl font-mono font-bold">{title}</h2>
    <div className="text-lg mt-4">{children}</div>
  </div>
)
