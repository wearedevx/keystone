const blockstack = require('blockstack')
const { SessionData } = require('blockstack/lib/auth/sessionData')
const fs = require('fs')
const pathUtil = require('path')
const {
  KEYSTONE_WEB,
  SESSION_FILENAME,
  KEYSTONE_ENV_CONFIG_PATH,
  KEYSTONE_HIDDEN_FOLDER,
} = require('@keystone.sh/core/lib/constants')

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
  try {
    const profile = await blockstack.lookupProfile(id)
    if (profile && profile.apps) {
      return profile.apps[KEYSTONE_WEB]
    }
    return false
  } catch (err) {
    // ignore error 400
    return false
  }
}

const getFilepath = ({ filename, apphub }) => {
  return `${apphub}${filename}`
}

const getSession = async path => {
  const session = await read({
    path: `${path}${pathUtil.sep}`,
    filename: SESSION_FILENAME,
  })
  const userSession = createUserSession(session)
  if (userSession && userSession.isUserSignedIn()) {
    return userSession
  }
  return false
}

const getProjectConfigFolderPath = (configFileName, currentPath = '.') => {
  if (fs.existsSync(pathUtil.join(currentPath, configFileName))) {
    return currentPath
  }

  if (fs.existsSync(pathUtil.join(currentPath, '..'))) {
    return getProjectConfigFolderPath(
      configFileName,
      pathUtil.join(currentPath, '..')
    )
  }

  if (!process.env.KEYSTONE_SHARED) {
    throw new Error('no ksconfig found')
  } else {
    return '.'
  }
}

const getProjectConfig = async (projectFilename = '.ksconfig') => {
  const projectConfigFolderPath = getProjectConfigFolderPath(projectFilename)
  let config
  try {
    config = await read({
      filename: pathUtil.join(projectConfigFolderPath, projectFilename),
    })
  } catch (err) {
    if (!process.env.KEYSTONE_SHARED) throw err
  }
  let envConfig
  try {
    envConfig = await getEnvConfig(projectConfigFolderPath)
  } catch (err) {}
  return {
    config: { ...config, ...envConfig },
    absoluteProjectPath: pathUtil.resolve(projectConfigFolderPath),
  }
}

const getEnvConfig = projectConfigFolderPath => {
  return read({
    path: pathUtil.join(
      projectConfigFolderPath,
      KEYSTONE_HIDDEN_FOLDER,
      KEYSTONE_ENV_CONFIG_PATH
    ),
    filename: '',
  })
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
