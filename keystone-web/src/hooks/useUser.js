import { useEffect, useState } from 'react'
import useStore from '../utils/store'
import { KEYSTONE_WEB } from '@keystone/core/dist/constants'

export default () => {
  const userSession = useStore(s => s.userSession)
  const userData = useStore(s => s.userData)
  const setUserData = useStore(s => s.setUserData)
  const loggedIn = useStore(s => s.loggedIn)
  const setLoggedIn = useStore(s => s.setLoggedIn)

  const [signinPending, setSigninPending] = useState(false)

  useEffect(() => {
    if (userSession && userSession.isUserSignedIn()) {
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

        //initWorkspace({ userSession, userData })
        /* save Data */
        // setUserData(userData)
        /* set Name */
        // setUserName(userData.username)
        /* set Public Key  */
        // setUserPublicKey(
        //   blockstack.getPublicKeyFromPrivate(userData.appPrivateKey)
        // )

        setLoggedIn(true)

        setSigninPending(false)

        /* --END LOADING-- */
        // setLoading(false)
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
