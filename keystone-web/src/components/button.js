import React from 'react'

export default ({ onClick, children, type = 'primary' }) => {
  let colors = ''
  switch (type) {
    case 'secondary':
      colors = 'text-white bg-secondary'
      break
    case 'warning':
      colors = 'text-gray-700 bg-yellow-600'
      break
    case 'danger':
      colors = 'text-white bg-red-600'
      break
    default:
      colors = 'text-white bg-primary'
  }

  return (
    <div
      className={`rounded font-bold ${colors} py-1 px-4 shadow-md text-center cursor-pointer ml-2`}
      onClick={onClick}
    >
      {children}
    </div>
  )
}
