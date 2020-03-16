import React from 'react'

export default ({ children }) => (
  <p className="font-mono text-blue-700 mt-6 text-xs">
    <span className="font-bold">Project ID:</span> {children}
  </p>
)
