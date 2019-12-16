import React from 'react'
import useUser from '../hooks/useUser'
import Base from './cards/base'

export default ({
  children,
  loggedOutText = 'You need to sign in with your Blockstack account continue.',
  loggedOutContent = null,
  redirectURI = '/',
}) => {
  const { loggedIn, redirectToSignIn } = useUser()

  return (
    <>
      {loggedIn && children}
      {!loggedIn && (
        <>
          <Base title={loggedOutText}>{loggedOutContent}</Base>
          <div className="my-4 flex flex-row w-2/4 justify-end">
            <div
              className="rounded font-bold text-white bg-primary py-1 px-4 shadow-md text-center cursor-pointer"
              onClick={() => redirectToSignIn(window.location.href)}
            >
              Sign in with Blockstack
            </div>
          </div>
        </>
      )}
    </>
  )
}
