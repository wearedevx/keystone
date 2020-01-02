import { useEffect, useState } from 'react'
import useStore from '../utils/store'
import { KEYSTONE_WEB, PUBKEY } from '@keystone.sh/core/lib/constants'
import * as blockstack from 'blockstack'
import { readFileFromGaia, writeFileToGaia } from '@keystone.sh/core/lib/file/gaia'

const savePublicKey = async userSession => {
  const userData = userSession.loadUserData()
  const publicKey = blockstack.getPublicKeyFromPrivate(userData.appPrivateKey)
  await writeFileToGaia(userSession, {
    path: PUBKEY,
    content: publicKey,
    sign: true,
    encrypt: false,
    json: false,
  })
}

export default () => {
  const userSession = useStore(s => s.userSession)
  const userData = useStore(s => s.userData)
  const setUserData = useStore(s => s.setUserData)
  const loggedIn = useStore(s => s.loggedIn)
  const setLoggedIn = useStore(s => s.setLoggedIn)

  const [signinPending, setSigninPending] = useState(false)

  useEffect(() => {
    if (userSession && userSession.isUserSignedIn()) {
      // makes sure the public key is available to others.
      const setPublicKey = async () => {
        let remotePubKey = undefined
        try {
          remotePubKey = await readFileFromGaia(userSession, {
            path: PUBKEY,
            decrypt: false,
            json: false,
            verify: true,
          })
        } catch (error) {
          console.error(error)
        } finally {
          if (!remotePubKey) {
            await savePublicKey(userSession)
          }
        }
      }
      setPublicKey()
      setLoggedIn(true)
    } else {
      setLoggedIn(false)
      setSigninPending(userSession.isSignInPending())
    }
  }, [])

  useEffect(() => {
    if (loggedIn) setUserData(userSession.loadUserData())
  }, [loggedIn])

  useEffect(() => {
    if (signinPending) {
      /* --START LOADING-- */
      userSession.handlePendingSignIn().then(async data => {
        /* initialize workspace, create required db files for keystone */

        setLoggedIn(true)

        setSigninPending(false)
      })
    }
  }, [signinPending])

  return {
    loggedIn,
    userData,
    redirectToSignIn: path => {
      userSession.redirectToSignIn(
        `${KEYSTONE_WEB}${path}`,
        `${KEYSTONE_WEB}/manifest.json`,
        ['email', 'publish_data', 'store_write']
      )

      const intervalId = setInterval(() => {
        const isUserSignedIn = userSession.isUserSignedIn()
        if (isUserSignedIn) {
          setSigninPending(false)
          clearInterval(intervalId)
        }
      }, 2000)
    },
    signUserOut: userSession.signUserOut,
    userSession,
  }
}
