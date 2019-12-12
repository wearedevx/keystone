const blockstack = require('blockstack')
const { AppConfig, UserSession } = blockstack
const { SessionData } = require('blockstack/lib/auth/sessionData')

const {createLocalStorage, updateLocalStorage } = require('./localStorage')

const blockstackLoader = (userCredentials = {}) => {
//   blockstackData = blockstackData || {
//     'blockstack': process.env.BLOCKSTACK,
//     'blockstack-gaia-hub-config': process.env.BLOCKSTACK_GAIA_HUB_CONFIG,
//     'blockstack-transit-private-key': process.env.BLOCKSTACK_TRANSIT_PRIVATE_KEY
//   }
  localStorage = createLocalStorage(userCredentials)
  //create global window with localStorage and location
  window = {
    localStorage,
    location: {
      origin: ''
    }
  }
  return localStorage
}

const createUserSession = (userCredentials = {}) => {
    blockstackLoader()
    const config = new AppConfig(['store_write', 'publish_data'])
    const userSession = new UserSession(config)
    const sessionData = SessionData.fromJSON(userCredentials)
    userSession.store.setSessionData(sessionData)
    return userSession

}

module.exports = {
    blockstackLoader,
    createUserSession
}