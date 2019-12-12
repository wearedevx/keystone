import { useEffect, useState } from 'react'
import useStore from '../utils/store'

console.log('fifjfjkefjhdjhddfbjksbfhh')

export default () => {
  const userSession = useStore(s => s.userSession)
  const userData = useStore(s => s.userData)
  const setUserData = useStore(s => s.setUserData)
  const loggedIn = useStore(s => s.loggedIn)
  const setLoggedIn = useStore(s => s.setLoggedIn)

  const [signinPending, setSigninPending] = useState(false)

  useEffect(() => {
    if (
      userSession &&
      userSession.isUserSignedIn &&
      userSession.isUserSignedIn()
    ) {
      setLoggedIn(true)
      setUserData(userSession.loadUserData())
    } else {
      setLoggedIn(false)

      if (userSession.isSignInPending) {
        setSigninPending(userSession.isSignInPending())
      }
    }
  }, [])

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

  return { loggedIn, userData }
}
