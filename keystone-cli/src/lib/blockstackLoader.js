const blockstack = require('blockstack')
const { SessionData } = require('blockstack/lib/auth/sessionData')
const fs = require('fs')
const path = require('path')
const {
  KEYSTONE_WEB,
  SESSION_FILENAME,
} = require('@keystone/core/lib/constants')

const { read } = require('../lib/cliStorage')
const { createLocalStorage } = require('./localStorage')

const { AppConfig, UserSession } = blockstack

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
      origin: KEYSTONE_WEB,
    },
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

const getAppHub = async id => {
  const profile = await blockstack.lookupProfile(id)
  if (profile && profile.apps) {
    return profile.apps[KEYSTONE_WEB]
  }
  return false
}

const getFilepath = ({ filename, apphub }) => {
  return `${apphub}${filename}`
}

const getSession = async path => {
  const session = await read({ path: `${path}/`, filename: SESSION_FILENAME })
  const userSession = createUserSession(session)
  if (userSession && userSession.isUserSignedIn()) {
    return userSession
  }
  return false
}

const getProjectConfigFolderPath = (configFileName, currentPath = '.') => {
  if (fs.existsSync(path.join(currentPath, configFileName))) {
    return currentPath
  }

  if (fs.existsSync(path.join(currentPath, '..'))) {
    return getProjectConfigFolderPath(
      configFileName,
      path.join(currentPath, '..')
    )
  }

  throw new Error('no ksconfig found')
}

const getProjectConfig = async (projectfileName = '.ksconfig') => {
  const projectConfigFolderPath = getProjectConfigFolderPath(projectfileName)

  return {
    config: await read({
      filename: path.join(projectConfigFolderPath, projectfileName),
    }),
    absoluteProjectPath: path.resolve(projectConfigFolderPath),
  }
}

module.exports = {
  blockstackLoader,
  createUserSession,
  getAppHub,
  getFilepath,
  getSession,
  getProjectConfig,
  getProjectConfigFolderPath,
}
