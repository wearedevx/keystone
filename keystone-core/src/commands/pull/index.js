const { getPath } = require('../../descriptor-path')
const {
  getLatestDescriptorByPath,
  uploadDescriptorForEveryone,
  getLatestProjectDescriptor,
  getLatestEnvDescriptor,
  getLatestMembersDescriptor,
  extractMembersByRole,
} = require('../../descriptor')
console.log('TCL: KeystoneError', 'puuuulll')
const KeystoneError = require('../../error')
const { findProjectByUUID, getProjects } = require('../../projects')
const {
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
  writeFileToDisk,
} = require('../../file')

const { ROLES } = require('../../constants')

const pull = async (
  userSession,
  { project, env, absoluteProjectPath, force = false }
) => {
  // create keystone cache folder
  const cacheFolder = getCacheFolder(absoluteProjectPath)

  // Can't pull if you have modified files under Keystone watch
  // - get all files from the cache folder
  // - check if they still exists on the current folder
  // - check if the content is the same.
  const changes = getModifiedFilesFromCacheFolder(
    cacheFolder,
    absoluteProjectPath
  )

  const uncommitted = changes.filter(change => change.status !== 'ok')

  if (uncommitted.length > 0 && !force) {
    throw new KeystoneError(
      'PullWhileFilesModified',
      'You should push your changes first.',
      uncommitted
    )
  }

  const { username } = userSession.loadUserData()
  const projects = await getProjects(userSession)
  const projectByUUID = findProjectByUUID(projects, project)
  if (!projectByUUID) {
    throw new Error('The project does not exist in user workspace')
  }

  await getLatestMembersDescriptor(userSession, {
    project,
    origin: projectByUUID.createdBy,
  })

  await getLatestProjectDescriptor(userSession, {
    project,
    origin: projectByUUID.createdBy,
  })

  const ownEnvDescriptor = await getLatestEnvDescriptor(userSession, {
    project,
    env,
    type: 'env',
  })

  const envMembersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
    type: 'members',
  })

  const { files } = ownEnvDescriptor.content

  const allLatestFileDescriptors = await Promise.all(
    files.map(file => {
      const filePath = getPath({
        project,
        env,
        type: 'file',
        blockstackId: username,
        filename: file.name,
      })

      return getLatestDescriptorByPath(userSession, {
        descriptorPath: filePath,
        members: extractMembersByRole(envMembersDescriptor, [
          ROLES.ADMINS,
          ROLES.CONTRIBUTORS,
        ]),
      })
    })
  )

  // For all files, update them
  const writtenFile = await Promise.all(
    allLatestFileDescriptors.map(async fileDescriptor => {
      writeFileToDisk(fileDescriptor, absoluteProjectPath)
      writeFileToDisk(fileDescriptor, cacheFolder)
      try {
        uploadDescriptorForEveryone(userSession, {
          members: extractMembersByRole(envMembersDescriptor, [
            ROLES.ADMINS,
            ROLES.CONTRIBUTORS,
          ]),
          env,
          project,
          descriptor: fileDescriptor,
          type: 'file',
        })
      } catch (error) {
        console.error(error)
        throw new Error(
          `Failed to upload file descriptor on private remote storage`
        )
      }
      return fileDescriptor
    })
  )

  return writtenFile
}

module.exports = pull
