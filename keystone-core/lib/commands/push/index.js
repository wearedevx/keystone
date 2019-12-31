const fs = require('fs')
const hash = require('object-hash')
const { getPath } = require('../../descriptor-path')

const {
  updateDescriptor,
  getLatestEnvDescriptor,
  updateFilesInEnvDesciptor,
} = require('../../descriptor')
const {
  writeFileToDisk,
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
} = require('../../file')

const push = async (
  userSession,
  { project, env, files, absoluteProjectPath }
) => {
  const { username } = userSession.loadUserData()

  // create keystone cache folder
  const cacheFolder = getCacheFolder(absoluteProjectPath)

  const envDescriptor = await getLatestEnvDescriptor(userSession, {
    project,
    env,
  })

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
