import { useEffect } from 'react'
import blockstack, { UserSession, AppConfig } from 'blockstack'
import useStore from '../utils/store'
// import initWorkspace from '../core/initWorkspace'

window.testUserSession = new UserSession({
  appConfig: new AppConfig(['email', 'store_write', 'publish_data']),
})

export default () => {
  const appConfig = useStore(s => s.appConfig)
  const setAppConfig = useStore(s => s.setAppConfig)
  const userSession = useStore(s => s.userSession)
  const setUserSession = useStore(s => s.setUserSession)

  // const {
  //   appConfig,
  //   userSession,
  //   userData,
  //   loading,
  //   signinPending,
  //   loggedIn,
  //   userPublicKey,
  //   projects,
  //   invitations,
  //   checkInvitation,
  //   setUserSession,
  //   setUserData,
  //   setUserName,
  //   setLoading,
  //   setUserPublicKey,
  //   setProjects,
  //   setInvitations,
  //   setSigninPending,
  //   setLoggedIn,
  //   setAppConfig,
  // } = useStore(state => ({
  //   appConfig: state.appConfig,
  //   loading: state.loading,
  //   userSession: state.userSession,
  //   userData: state.userData,
  //   signinPending: state.signinPending,
  //   userPublicKey: state.userPublicKey,
  //   projects: state.projects,
  //   invitations: state.invitations,
  //   checkInvitation: state.checkInvitation,
  //   setUserSession: state.setUserSession,
  //   setUserData: state.setUserData,
  //   setUserName: state.setUserName,
  //   setLoading: state.setLoading,
  //   setUserPublicKey: state.setUserPublicKey,
  //   setProjects: state.setProjects,
  //   setInvitations: state.setInvitations,
  //   setSigninPending: state.setSigninPending,
  //   setLoggedIn: state.setLoggedIn,
  //   setAppConfig: state.setAppConfig,
  // }))

  useEffect(() => {
    if (!appConfig) {
      setAppConfig(new AppConfig(['store_write', 'publish_data']))
    }
  }, [])

  useEffect(() => {
    if (!userSession) {
      const session = new UserSession({
        appConfig,
      })
      console.log('TCL: session', session.isUserSignedIn)
      setUserSession(session)
    }
    // if ((!userSession || !userSession.isUserSignedIn()) && appConfig) {
    //   setUserSession(
    //     new UserSession({
    //       appConfig,
    //     })
    //   )
    // }
  }, [appConfig])

  return {
    userSession,
    // loading,
    // userSession,
    // userData,
    // userPublicKey,
    // projects,
    // invitations,
    // checkInvitation,
    // setUserSession,
    // setUserData,
    // setUserName,
    // setLoading,
    // setUserPublicKey,
    // setProjects,
    // setInvitations,
    // setSigninPending,
  }
}
