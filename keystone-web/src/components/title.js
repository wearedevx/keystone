import React from 'react'

export default ({ children, id = '' }) => {
  return (
    <h2
      id={id}
      className="text-2xl font-bold mb-6 px-3 py-2 relative text-white inline-block"
    >
      <div className="z-10 relative">{children}</div>
      <div className="bg-black transform skew-y-3 absolute z-0 left-0 right-0 top-0 bottom-0"></div>
    </h2>
  )
}
