const { getPath } = require('../../descriptor-path')
const {
  getLatestDescriptorByPath,
  uploadDescriptorForEveryone,
  getLatestProjectDescriptor,
  getLatestEnvDescriptor,
  getLatestMembersDescriptor,
  extractMembersByRole,
  getDescriptor,
} = require('../../descriptor')

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

  // TODO: we should use createdBy only when it's the first time
  // we pull and we don't have any member list yet.
  await getLatestMembersDescriptor(userSession, {
    project,
    origin: projectByUUID.createdBy,
  })

  await getLatestProjectDescriptor(userSession, {
    project,
    origin: projectByUUID.createdBy,
  })

  const ownEnvDescriptor = await getDescriptor(userSession, {
    env,
    project,
    type: 'env',
  })

  const envDescriptor = await getLatestEnvDescriptor(userSession, {
    project,
    env,
    type: 'env',
  })

  if (envDescriptor.checksum === ownEnvDescriptor.checksum) {
    return [{ descriptorUpToDate: true }]
  }

  const envMembersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
    type: 'members',
  })

  const { files } = envDescriptor.content

  const allLatestFileDescriptors = await Promise.all(
    files.map(file => {
      const ownFile = ownEnvDescriptor.content.files.find(
        f => f.name === file.name
      )
      if (!ownFile || (ownFile && file.checksum !== ownFile.checksum)) {
        const filePath = getPath({
          project,
          env,
          type: 'file',
          blockstackId: username,
          filename: file.name,
        })

        const fileDescriptor = getLatestDescriptorByPath(userSession, {
          descriptorPath: filePath,
          members: extractMembersByRole(envMembersDescriptor, [
            ROLES.ADMINS,
            ROLES.CONTRIBUTORS,
          ]),
        })

        // Write files on disk
        writeFileToDisk(fileDescriptor, absoluteProjectPath)
        writeFileToDisk(fileDescriptor, cacheFolder)

        // Upload in own hub for everyone
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
        return { fileDescriptor, updated: true }
      }
      return { fileDescriptor: ownFile, updated: false }
    })
  )

  return allLatestFileDescriptors
}

module.exports = pull
