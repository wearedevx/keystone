import React from 'react'
import Button from './button'
import useUser from '../hooks/useUser'

export default ({ disabled = false, onClick, children, type = 'primary' }) => {
  const { signUserOut, userSession } = useUser()
  const { username = 'username missing' } = userSession.loadUserData()

  return (
    <div className="my-4 flex flex-row w-2/4 justify-center">
      <div className="flex flex-col items-center">
        <Button disabled={disabled} onClick={onClick}>
          {children}
        </Button>
        <div className="text-xs text-gray-600 mt-1 flex flex-col items-center">
          <div>connected with {username}</div>
          <div
            className="underline font-bold cursor-pointer"
            onClick={() => signUserOut()}
          >
            Log out
          </div>
        </div>
      </div>
    </div>
  )
}
