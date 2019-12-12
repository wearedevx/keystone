import React, { useState, useEffect } from "react"
import useStore from "../store"
import { checkInvitations } from "../core"

export const useProjects = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  const {
    userSession,
    userData,
    projects,
    setUserData,
    setProjects,
    loggedIn,
    setLoading,
  } = useStore(state => ({
    userSession: state.userSession,
    userData: state.userData,
    loggedIn: state.loggedIn,
    projects: state.projects,
    setUserData: state.setUserData,
    setProjects: state.setProjects,
    setInvitations: state.setInvitations,
    setLoading: state.setLoading,
  }))

  useEffect(() => {
    if (loggedIn) {
      const getFiles = async () => {
        setLoading(true)
        const projectsUpdate = await checkInvitations(userSession)
        console.log("projectsUpdate", projectsUpdate)
        setProjects(projectsUpdate)
        setLoading(false)
      }
      getFiles()
    }
  }, [loggedIn])

  return { setProjects, projects }
}
