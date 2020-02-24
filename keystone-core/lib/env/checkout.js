const debug = require('debug')('keystone:core:env')
const fs = require('fs')
const path = require('path')
const KeystoneError = require('../error')
const {
  getLatestProjectDescriptor,
  getLatestEnvDescriptor,
} = require('../descriptor')

const {
  changeEnvConfig,
  getModifiedFilesFromCacheFolder,
  getCacheFolder,
} = require('../file/disk')

const checkoutEnv = async (
  userSession,
  { project, env, absoluteProjectPath, force = false }
) => {
  const modifiedFiles = await getModifiedFilesFromCacheFolder(
    getCacheFolder(absoluteProjectPath),
    absoluteProjectPath
  ).filter(f => f.status !== 'ok')
  if (modifiedFiles.length > 0 && !force) {
    throw new KeystoneError('PendingModification', '', modifiedFiles)
  }
  const projectDescriptor = await getLatestProjectDescriptor(userSession, {
    project,
    type: 'project',
  })

  // Retrieve updated project descriptor
  const envFound = projectDescriptor.content.env.find(
    envObject => envObject === env
  )

  if (envFound) {
    await getLatestEnvDescriptor(userSession, {
      project,
      env,
      type: 'env',
    })

    return await changeEnvConfig({
      env,
      absoluteProjectPath,
    })
  }
  throw new Error(`The environment ${env} is not defined in this project`)
}

module.exports = {
  checkoutEnv,
}
