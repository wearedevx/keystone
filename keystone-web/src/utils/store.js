import create from 'zustand'
import deepCopy from './deepCopy'

// Declare global application state
const store = {
  appConfig: undefined,
  userSession: undefined,
  signinPending: false,
  userData: undefined,
  userName: undefined,
  loading: false,
  userPublicKey: undefined,
  //   projects: [],
  //   invitations: [],
  checkInvitation: false,
  error: undefined,
  loggedIn: false,
}

// Generate setter function names from props global state
const getSetterName = prop => {
  return `set${prop.charAt(0).toUpperCase()}${prop.slice(1)}`
}

// Add setters for each property in the state
const addSetters = (props, set) => {
  return Object.keys(props).reduce((state, key) => {
    const s = state
    s[key] = props[key]
    s[getSetterName(key)] = value => {
      const partialState = {}
      partialState[key] = value
      set(() => ({ ...deepCopy(partialState) }))
    }
    return s
  }, {})
}

// Create store
const [useStore] = create(set => {
  return addSetters(store, set)
})

export default useStore
