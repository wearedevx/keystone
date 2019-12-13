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

  if (!token || !id)
    throw new KeystoneError(
      'MissingParams',
      'Missing token and blockstack id query strings',
      location
    )
  // Retrieve session
  const blockstackSessionStore = JSON.stringify(
    userSession.store.getSessionData()
  )

  // Upload and encrypt with the public key which is the token
  await writeFileToGaia(userSession, {
    path: `${token}.json`,
    content: blockstackSessionStore,
    encrypt: token,
  })

  setTerminalConnected(true)
}

export default () => {
  const { loggedIn, userData, redirectToSignIn, userSession } = useUser()
  const [terminalConnected, setTerminalConnected] = useState(false)
  const [missingParams, setMissingParams] = useState(false)
  const [error, setError] = useState(false)

  if (loggedIn && !missingParams) {
    try {
      connectTerminal({
        location: window.location,
        userSession,
        setTerminalConnected,
      })
    } catch (error) {
      switch (error.code) {
        case 'MissingParams':
          setMissingParams(true)
          break
        default:
          setError(error.message)
          throw error
      }
    }
  }

  return (
    <div className="flex flex-col items-center ">
      <div className="shadow-md rounded p-4 bg-white w-2/4">
        {loggedIn && missingParams && (
          <>
            <h2 className="text-xl mb-4">
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
        {loggedIn && terminalConnected && (
          <>
            <h2 className="text-xl">
              <span
                role="img"
                aria-label="A party popper, as explodes in a shower of confetti and streamers at a celebration"
              >
                🎉
              </span>
              Your terminal is connected.
            </h2>
            <div>
              <Link to="/" className="text-blue-500 underline mr-1">
                Read the documentation
              </Link>
              or type `ks --help` in your terminal to start with Keystone.
            </div>
          </>
        )}

        {!loggedIn && (
          <h2 className="text-xl">
            You need to sign in with your Blockstack account to connect your
            terminal.
          </h2>
        )}

        {error && <h2 className="text-xl">{error}</h2>}
      </div>

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

      {/* <InvitationBannerJoin project={project} from={from} id={id} /> */}
      {/* <InvitationBoard project={project} /> */}
    </div>
  )
}
