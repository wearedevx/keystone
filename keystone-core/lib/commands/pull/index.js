const daffy = require('daffy')
const { deepCopy } = require('../../utils')
const { getPath } = require('../../descriptor-path')
const {
  getLatestDescriptorByPath,
  uploadDescriptorForEveryone,
  getLatestProjectDescriptor,
  getLatestEnvDescriptor,
  getLatestMembersDescriptor,
  extractMembersByRole,
  getDescriptor,
  mergeContents,
} = require('../../descriptor')

const KeystoneError = require('../../error')
const { findProjectByUUID, getProjects } = require('../../projects')
const {
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
  writeFileToDisk,
  readFileFromDisk,
} = require('../../file')

const { ROLES } = require('../../constants')

const pull = async (
  userSession,
  { project, env, absoluteProjectPath, force = false, cache = true, origin }
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
  if (cache) {
    const uncommitted = changes.filter(change => change.status !== 'ok')
    console.log('UNCOMMITED', uncommitted)

    if (uncommitted.length > 0 && !force) {
      // throw new KeystoneError(
      //   'PullWhileFilesModified',
      //   'You should push your changes first.',
      //   uncommitted
      // )
    }
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
    origin: origin || projectByUUID.createdBy,
  })

  await getLatestProjectDescriptor(userSession, {
    project,
    origin: origin || projectByUUID.createdBy,
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

  if (
    ownEnvDescriptor &&
    envDescriptor.checksum === ownEnvDescriptor.checksum
  ) {
    return [{ descriptorUpToDate: true }]
  }

  const envMembersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
    type: 'members',
  })

  const { files } = envDescriptor.content

  const allLatestFileDescriptors = await Promise.all(
    files.map(async file => {
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

        const fileDescriptor = await getLatestDescriptorByPath(userSession, {
          descriptorPath: filePath,
          members: extractMembersByRole(envMembersDescriptor, [
            ROLES.ADMINS,
            ROLES.CONTRIBUTORS,
          ]),
        })

        const fileModified = changes.find(c => {
          const filename = c.path.replace(`${absoluteProjectPath}/`, '')
          return file.name === filename
        })

        const fileDescriptorToWriteOnDisk = deepCopy(fileDescriptor)
        if (fileModified && fileModified.status !== 'ok') {
          const currentVersion = await readFileFromDisk(fileModified.path)
          const base = daffy.applyPatch(
            fileDescriptor.content,
            fileDescriptor.history.find(
              h => h.version === fileDescriptor.version - 1
            ).content || ''
          )

          fileDescriptorToWriteOnDisk.content = mergeContents({
            left: currentVersion,
            right: fileDescriptor.content,
            base,
          })
        }

        // Write files on disk
        writeFileToDisk(fileDescriptorToWriteOnDisk, absoluteProjectPath)
        writeFileToDisk(fileDescriptorToWriteOnDisk, cacheFolder)

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
