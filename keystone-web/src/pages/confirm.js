import React, { useState } from 'react'
import useUser from '../hooks/useUser'
import { Link } from '@reach/router'
import queryString from 'query-string'
import KeystoneError from '@keystone/core/lib/error'
import { writeFileToGaia } from '@keystone/core/lib/file/gaia'

const connectTerminal = async ({
  location,
  userSession,
  setTerminalConnected,
}) => {
  const { token, id } = queryString.parse(location.search)
  const { username } = userSession.loadUserData()

  if (!token || !id)
    throw new KeystoneError(
      'MissingParams',
      'Missing token and blockstack id query strings',
      location
    )

  // check that the blockstack id sent is matching the one connected to the browser
  if (id !== username) {
    throw new KeystoneError(
      'AccountMismatch',
      `The blockstack account sent from the terminal (${id}) is not the same than the one connected to the browser (${username}).`,
      { terminalAccount: id, browserAccount: username }
    )
  }

  // Retrieve session
  const blockstackSessionStore = JSON.stringify(
    userSession.store.getSessionData()
  )

  // Upload and encrypt with the public key which is the token
  const file = await writeFileToGaia(userSession, {
    path: `${token}.json`,
    content: blockstackSessionStore,
    encrypt: token,
  })

  if (file) {
    setTerminalConnected(true)
  }
}

export default () => {
  const { loggedIn, redirectToSignIn, userSession } = useUser()
  const [terminalConnected, setTerminalConnected] = useState(false)
  const [missingParams, setMissingParams] = useState(false)
  const [error, setError] = useState(false)
  const [connecting, setConnecting] = useState(false)

  if (loggedIn && !connecting) {
    setConnecting(true)
    connectTerminal({
      location: window.location,
      userSession,
      setTerminalConnected,
      // userData,
    }).catch(error => {
      switch (error.code) {
        case 'MissingParams':
          setMissingParams(true)
          break
        case 'AccountMismatch':
          setError(error.message)
          break
        default:
          setError(error.message)
          throw error
      }
    })
  }

  return (
    <div className="flex flex-col items-center ">
      {error && (
        <>
          <h2 className="text-xl text-red-600">
            <span
              role="img"
              aria-label="A cartoon-styled representation of a collision"
            >
              💥
            </span>
            {error}
          </h2>
        </>
      )}

      {loggedIn && missingParams && (
        <>
          <h2 className="text-xl mb-4 text-red-600">
            <span
              role="img"
              aria-label="A cartoon-styled representation of a collision"
            >
              💥
            </span>
            Your link is malformed. Please open an issue on GitHub.
          </h2>
          <div>
            Or check that the link in your browser is the same than the one
            provided by your terminal with the command `ks login`.
          </div>
        </>
      )}

      {!error && !missingParams && (
        <div className="shadow-md rounded p-4 bg-white w-2/4">
          {loggedIn && (
            <>
              {terminalConnected && (
                <>
                  <h2 className="text-xl">
                    <span
                      role="img"
                      aria-label="A party popper, as explodes in a shower of confetti and streamers at a celebration"
                    >
                      🎉
                    </span>
                    Your terminal is connected. You can close this window.
                  </h2>
                  <div>
                    <Link to="/" className="text-blue-500 underline mr-1">
                      Read the documentation
                    </Link>
                    or type `ks --help` in your terminal to start with Keystone.
                  </div>
                </>
              )}

              {!terminalConnected && (
                <>
                  <h2 className="text-xl">
                    <span
                      role="img"
                      aria-label="A key, as opens a door or lock"
                    >
                      🔑
                    </span>
                    Connecting your terminal...
                  </h2>
                  <div>It should take less than a minute.</div>
                </>
              )}
            </>
          )}

          {!loggedIn && (
            <h2 className="text-xl">
              You need to sign in with your Blockstack account to connect your
              terminal.
            </h2>
          )}
        </div>
      )}

      {!missingParams && !error && !loggedIn && (
        <div className="my-4 flex flex-row w-2/4 justify-end">
          <div
            className="rounded font-bold text-white bg-primary py-1 px-4 shadow-md text-center cursor-pointer"
            onClick={() => redirectToSignIn('/confirm')}
          >
            Sign in with Blockstack
          </div>
        </div>
      )}
    </div>
  )
}
