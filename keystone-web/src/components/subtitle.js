import React from 'react'

export default ({ children }) => {
  return (
    <h3 className="text-xl font-bold mb-6 px-3 py-2 relative inline-block">
      <div className="z-10 relative">{children}</div>
      <div className="border-b-4 border-t-4 border-r-4 border-black transform skew-y-3 absolute z-0 left-0 right-0 top-0 bottom-0"></div>
    </h3>
  )
}
