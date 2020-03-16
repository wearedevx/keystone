import React, { useState } from 'react'
import useUser from '../hooks/useUser'
import { Link } from '@reach/router'
import queryString from 'query-string'
import KeystoneError from '@keystone.sh/core/lib/error'
import { writeFileToGaia } from '@keystone.sh/core/lib/file/gaia'
import { LOGIN_KEY_PREFIX } from '@keystone.sh/core/lib/constants'
import ErrorCard from '../components/cards/error'
import SuccessCard from '../components/cards/success'
import WithLoggin from '../components/withLoggin'

const connectTerminal = async ({
  location,
  userSession,
  setTerminalConnected,
}) => {
  const { token, id } = queryString.parse(location.search)
  const { username } = userSession.loadUserData()

  // check that the blockstack id sent is matching the one connected to the browser
  if (!username) {
    throw new KeystoneError(
      'NoUsername',
      `Your account does not have a Blockstack ID.`,
      { terminalAccount: id, browserAccount: username }
    )
  }
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
    path: `${LOGIN_KEY_PREFIX}${token}.json`,
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
  const [NoUsername, setNoUsername] = useState(false)
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
        case 'NoUsername':
          setError(error.message)
          setNoUsername(true)
          break
        default:
          setError(error.message)
          throw error
      }
    })
  }

  return (
    <WithLoggin redirectURI="/invite">
      {error && <ErrorCard title={error} />}

      {NoUsername && (
        <ErrorCard
          title={`If you just created a new account on Blockstack, it can take up to a few hours for
         Blockstack to validate it.`}
        >
          Please, come back later or reach Blockstack directly.
        </ErrorCard>
      )}

      {missingParams && (
        <ErrorCard
          title={`Your link is malformed. Please open an issue on GitHub.`}
        >
          Or check that the link in your browser is the same than the one
          provided by your terminal with the command `ks login`.
        </ErrorCard>
      )}

      {!error && !missingParams && (
        <div className="p-4 bg-white w-2/4">
          {terminalConnected && (
            <SuccessCard
              title={`Your terminal is connected. You can close this window.`}
            >
              <div>
                <Link to="/" className="text-blue-500 underline mr-1">
                  Read the documentation
                </Link>
                or type{' '}
                <span className="font-mono text-sm font-bold">`ks --help`</span>{' '}
                in your terminal to start with Keystone.
              </div>
            </SuccessCard>
          )}

          {!terminalConnected && (
            <SuccessCard title={`Connecting your terminal...`}>
              It should take less than a minute.
            </SuccessCard>
          )}
        </div>
      )}
    </WithLoggin>
  )
}
