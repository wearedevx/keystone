import React from 'react'

export default ({ onClick, children, type = 'primary', disabled = false }) => {
  let colors = ''
  let disabledStyle = disabled ? 'opacity-25 cursor-not-allowed' : ''
  let clickAction = disabled ? () => {} : onClick

  switch (type) {
    case 'secondary':
      colors = 'text-white bg-black'
      break
    case 'warning':
      colors = 'text-gray-700 bg-yellow-600'
      break
    case 'danger':
      colors = 'text-white bg-red-600'
      break
    default:
      colors = 'text-white bg-green-700'
  }

  return (
    <div
      className={`font-bold ${colors} py-1 px-4 text-center cursor-pointer ${disabledStyle}`}
      onClick={clickAction}
    >
      {children}
    </div>
  )
}
