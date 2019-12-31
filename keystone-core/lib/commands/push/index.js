const fs = require('fs')
const hash = require('object-hash')
const { getPath } = require('../../descriptor-path')

const {
  updateDescriptor,
  getLatestEnvDescriptor,
  updateFilesInEnvDesciptor,
  getLatestDescriptorByPath,
  getLatestMembersDescriptor,
  extractMembersByRole,
  getOwnDescriptorByPath,
} = require('../../descriptor')
const {
  writeFileToDisk,
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
} = require('../../file')
const KeystoneError = require('../../error')
const { ROLES } = require('../../constants')

const push = async (
  userSession,
  { project, env, files, absoluteProjectPath }
) => {
  const { username } = userSession.loadUserData()

  // create keystone cache folder
  const cacheFolder = getCacheFolder(absoluteProjectPath)

  const membersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
  })
  const members = extractMembersByRole(membersDescriptor, [
    ROLES.CONTRIBUTORS,
    ROLES.ADMINS,
  ])
  // Retrieve the latest version of the file from everyone.
  const envPath = getPath({ project, env, type: 'env', blockstackId: username })

  const envDescriptor = await getLatestDescriptorByPath(userSession, {
    descriptorPath: envPath,
    members,
  })

  const previousEnvDescriptor = await getOwnDescriptorByPath(userSession, {
    descriptorPath: envPath,
  })

  if (envDescriptor.checksum !== previousEnvDescriptor.checksum) {
    throw new KeystoneError(
      'PullBeforeYouPush',
      'A version of this file exist with another content.\nPlease pull before pushing your file.'
    )
  }

  const pushedFiles = await Promise.all(
    files.map(async file => {
      const { filename, fileContent } = file
      if (
        envDescriptor.content.files.find(
          f => f.name === filename && f.checksum === hash(fileContent)
        )
      ) {
        return { updated: false, filename }
      }

      const filePath = getPath({
        project,
        env,
        filename,
        blockstackId: username,
        type: 'file',
      })

      const fileDescriptor = await updateDescriptor(userSession, {
        project,
        env,
        type: 'file',
        content: fileContent,
        descriptorPath: filePath,
        name: filename,
      })

      writeFileToDisk(fileDescriptor, cacheFolder)

      return { updated: true, filename }
    })
  )

  // If file is not present, add it. If present, update checksum
  await updateFilesInEnvDesciptor(userSession, {
    files,
    envDescriptor,
    project,
    env,
  })

  return pushedFiles
}

/**
 * Push only modified files present in cache folder.
 */
const pushModifiedFiles = async (
  userSession,
  { project, env, absoluteProjectPath }
) => {
  // create keystone cache folder
  const cacheFolder = getCacheFolder(absoluteProjectPath)

  const changes = getModifiedFilesFromCacheFolder(
    cacheFolder,
    absoluteProjectPath
  )

  const modifiedFiles = changes.filter(c => c.status !== 'ok')

  if (modifiedFiles.length === 0) {
    console.log('No modified files. Nothing to push.')
    return
  }

  // return

  const formatModifiedFiles = modifiedFiles.map(({ path }) => ({
    filename: path.replace(`${absoluteProjectPath}/`, ''),
    fileContent: fs.readFileSync(path).toString(),
  }))

  const pushedFiles = push(userSession, {
    project,
    env,
    files: formatModifiedFiles,
    absoluteProjectPath,
  })
  return pushedFiles
}

module.exports = { push, pushModifiedFiles }
