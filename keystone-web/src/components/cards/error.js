import React from 'react'

export default ({ title, children }) => (
  <div>
    <h2 className="text-xl mb-4 text-red-600">
      <span
        role="img"
        aria-label="A cartoon-styled representation of a collision"
      >
        💥
      </span>
      {title}
    </h2>
    <div>{children}</div>
  </div>
)
