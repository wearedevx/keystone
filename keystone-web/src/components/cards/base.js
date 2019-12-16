import React from 'react'

export default ({ title, children }) => {
  return (
    <div className="shadow-md rounded p-4 bg-white w-2/4">
      <h2 className="text-xl">{title}</h2>
      <div>{children}</div>
    </div>
  )
}
