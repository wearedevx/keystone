import React from 'react'

export default ({ title, children }) => {
  return (
    <div className="p-4 bg-white w-2/4 text-center">
      <h2 className="text-xl mb-4 font-mono">{title}</h2>
      <div>{children}</div>
    </div>
  )
}
