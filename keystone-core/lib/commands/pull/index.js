const daffy = require('daffy')
const path = require('path')

const KeystoneError = require('../../error')
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

// const KeystoneError = require('../../error')
const { findProjectByUUID, getProjects } = require('../../projects')
const {
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
  writeFileToDisk,
  readFileFromDisk,
  isFileExist,
} = require('../../file')

const { ROLES } = require('../../constants')

const { updateFilesInEnvDesciptor } = require('../../descriptor')

/**
 * For all files in env descriptor, and not in cache folder,
 * write them in cache folder and work folder.
 * @param {*} userSession
 * @param {*} param1
 */
const checkFilesToWrite = (
  userSession,
  {
    project,
    env,
    envMembersDescriptor,
    envDescriptor,
    cacheFolder,
    absoluteProjectPath,
  }
) => {
  const { username } = userSession.loadUserData()

  return Promise.all(
    envDescriptor.content.files.map(async file => {
      if (!isFileExist(path.join(cacheFolder, file.name))) {
        const filePath = getPath({
          project,
          env,
          type: 'file',
          filename: file.name,
          blockstackId: username,
        })

        const fileDescriptor = await getLatestDescriptorByPath(userSession, {
          descriptorPath: filePath,
          members: extractMembersByRole(envMembersDescriptor, [
            ROLES.ADMINS,
            ROLES.CONTRIBUTORS,
          ]),
        })

        writeFileToDisk(fileDescriptor, cacheFolder)
        writeFileToDisk(fileDescriptor, absoluteProjectPath)
      }
    })
  )
}

const pull = async (
  userSession,
  { project, env, absoluteProjectPath, force = false, cache = true, origin }
) => {
  if (!env) {
    const latestProjectDescriptor = await getLatestProjectDescriptor(
      userSession,
      {
        project,
      }
    )
    throw new KeystoneError(
      'MissingEnv',
      `You need to checkout an env in order to pull files.`,
      { envs: latestProjectDescriptor.content.env }
    )
  }
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

    if (uncommitted.length > 0 && !force) {
      throw new KeystoneError(
        'PullWhileFilesModified',
        'You should push your changes first.',
        uncommitted
      )
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

  const envMembersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
    type: 'members',
  })

  if (
    ownEnvDescriptor &&
    envDescriptor.checksum === ownEnvDescriptor.checksum
  ) {
    // Check files to write
    await checkFilesToWrite(userSession, {
      project,
      env,
      envMembersDescriptor,
      envDescriptor,
      cacheFolder,
      absoluteProjectPath,
    })

    if (!force) return [{ descriptorUpToDate: true }]
  }

  const { files } = envDescriptor.content

  const allLatestFileDescriptors = await Promise.all(
    files.map(async file => {
      const ownFile =
        ownEnvDescriptor &&
        ownEnvDescriptor.content.files.find(f => f.name === file.name)

      if (
        !ownFile ||
        (ownFile && file.checksum !== ownFile.checksum) ||
        force
      ) {
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
        let conflict

        if (fileModified && fileModified.status !== 'ok') {
          const currentVersion = await readFileFromDisk(fileModified.path)
          const base = daffy.applyPatch(
            fileDescriptor.content,
            fileDescriptor.history.find(
              h => h.version === fileDescriptor.version - 1
            ).content || ''
          )

          const mergeResult = mergeContents({
            left: currentVersion,
            right: fileDescriptor.content,
            base,
          })
          fileDescriptorToWriteOnDisk.content = mergeResult.result
          conflict = mergeResult.conflict
        } else {
          writeFileToDisk(fileDescriptorToWriteOnDisk, cacheFolder)
        }

        // Write files on disk
        writeFileToDisk(fileDescriptorToWriteOnDisk, absoluteProjectPath)

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
        return { fileDescriptor, updated: true, conflict }
      }
      return { fileDescriptor: ownFile, updated: false }
    })
  )
  await updateFilesInEnvDesciptor(userSession, {
    files: allLatestFileDescriptors
      .filter(
        file =>
          !(typeof file.conflict === 'boolean') &&
          !file.conflict &&
          file.updated
      )
      .map(file => ({
        filename: file.fileDescriptor.name,
        fileContent: file.fileDescriptor.content,
      })),
    envDescriptor,
    project,
    env,
  })

  return allLatestFileDescriptors
}

module.exports = pull
