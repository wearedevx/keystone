import React from 'react'
import useUser from '../hooks/useUser'
import Base from './cards/base'
import Button from './button'

export default ({
  children,
  loggedOutText = 'You need to sign in with your Blockstack account to continue.',
  loggedOutContent = null,
}) => {
  const { loggedIn, redirectToSignIn } = useUser()

  return (
    <>
      {loggedIn && children}
      {!loggedIn && (
        <>
          <Base title={loggedOutText}>{loggedOutContent}</Base>
          <div className="my-4 flex flex-row w-2/4 justify-center">
            <Button onClick={() => redirectToSignIn(window.location.href)}>
              Sign in with Blockstack
            </Button>
          </div>
        </>
      )}
    </>
  )
}
