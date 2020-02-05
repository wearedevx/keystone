const fs = require('fs')
const {
  SHARED_MEMBER,
  PUBKEY,
  KEYSTONE_CONFIG_PATH,
  KEYSTONE_ENV_CONFIG_PATH,
  KEYSTONE_HIDDEN_FOLDER,
} = require('@keystone.sh/core/lib/constants')

const pathUtil = require('path')

const gaiaUtil = require('@keystone.sh/core/lib/file/gaia')
const diskUtil = require('@keystone.sh/core/lib/file/disk')

const fsp = fs.promises
gaiaUtil.readFileFromGaia = async (
  userSession,
  { path, json = true, origin, decrypt }
) => {
  const { username } = userSession.loadUserData()
  if (decrypt && userSession.sharedPrivateKey) {
    decrypt = userSession.sharedPrivateKey
  }
  console.log(path)
  path = path.replace(/\//g, '|')
  console.log('READ FILE ', path)

  path = `${username}--${path}`
  path = pathUtil.join(__dirname, '..', pathUtil.sep, 'hub', path)

  try {
    const fetchedFile = await fsp.readFile(path, 'utf-8')
    if (json) {
      const fetchedFileJSON = JSON.parse(fetchedFile)
      return fetchedFileJSON
    }

    return fetchedFile
  } catch (error) {
    // we can't retrieve the file from remote
    // make it like a 404 not found
    console.error(error)
    return null
  }
}

gaiaUtil.writeFileToGaia = async (userSession, { path, content, encrypt }) => {
  const { username } = userSession.loadUserData()
  path = path.replace(/\//g, '|')
  console.log('WRITE FILE ', path)
  path = `${username}--${path}`
  path = pathUtil.join(__dirname, '..', pathUtil.sep, 'hub', path)
  await fsp.writeFile(path, content)
  return content
}

gaiaUtil.getPubkey = async (
  userSession,
  { blockstack_id: blockstackId, publicKey }
) => {
  if (publicKey) return publicKey

  if (new RegExp(SHARED_MEMBER).test(blockstackId))
    return blockstackId.split('-')[1]
  const pubkeyFile = await gaiaUtil.readFileFromGaia(userSession, {
    decrypt: false,
    path: `${PUBKEY}`,
    origin: blockstackId,
    json: false,
    verify: true,
  })
  if (pubkeyFile) {
    return pubkeyFile
  }
  throw new Error(
    `Keystone public application key not found on ${blockstackId}`
  )
}

gaiaUtil.listFilesFromGaia = async userSession => {
  return fs.readdirSync(pathUtil.join(__dirname, '..', pathUtil.sep, 'hub'))
}
module.exports = gaiaUtil

diskUtil.writeFileToDisk = (fileDescriptor, absoluteProjectPath) => {
  try {
    let pathFile

    if (fileDescriptor.name === KEYSTONE_ENV_CONFIG_PATH) {
      pathFile = pathUtil.join(
        __dirname,
        '../local',
        KEYSTONE_HIDDEN_FOLDER,
        fileDescriptor.name
      )
    } else {
      pathFile = pathUtil.join(__dirname, '../local', fileDescriptor.name)
    }
    const lastIndex = pathFile.lastIndexOf(pathUtil.sep)
    const folder = pathFile.substring(0, lastIndex)

    console.log('WRITE FILE DISK ', fileDescriptor.name)
    if (folder) fs.mkdirSync(folder, { recursive: true })

    // if JSON object, stringify
    let { content } = fileDescriptor
    if (typeof content === 'object') content = JSON.stringify(content)

    fs.writeFile(pathFile, content, err => {
      if (err) throw new Error(err)
    })
  } catch (err) {
    throw new Error(err)
  }
  return fileDescriptor
}

diskUtil.readFileFromDisk = async filename => {
  const buffer = await fsp.readFile(pathUtil.join(__dirname, '..', filename))
  return buffer.toString('utf-8')
}

module.exports = gaiaUtil

diskUtil.deleteFileFromDisk = path => {
  debug('deleteFileFromDisk', path)
  return fs.unlinkSync(pathUtil.join(__dirname, '../local', path))
}

diskUtil.getKeystoneFolder = absoluteProjectPath => {
  const keystoneFolder = pathUtil.join(
    __dirname,
    '../local',
    KEYSTONE_HIDDEN_FOLDER
  )
  if (!fs.existsSync(keystoneFolder)) {
    const gitignorePath = pathUtil.join(__dirname, '../local', '.gitignore')
    let gitignoreContent
    try {
      gitignoreContent = fs.readFileSync(gitignorePath).toString()
    } catch (err) {}
    // check in gitignoe if .keystone/ is present
    if (
      !gitignoreContent ||
      (gitignoreContent &&
        !gitignoreContent
          .split('\n')
          .find(line => line === `${KEYSTONE_HIDDEN_FOLDER}/`))
    ) {
      fs.appendFileSync(gitignorePath, `\n${KEYSTONE_HIDDEN_FOLDER}/`)
      fs.mkdirSync(keystoneFolder)
    }
  }

  return keystoneFolder
}
diskUtil.getCacheFolder = absoluteProjectPath => {
  const cacheFolder = pathUtil.join(
    getKeystoneFolder(absoluteProjectPath),
    `/cache/`
  )
  if (!fs.existsSync(cacheFolder)) {
    fs.mkdirSync(cacheFolder, { recursive: true })
  }

  console.log('CACHE FOLDER', cacheFolder)

  return cacheFolder
}

diskUtil.isFileExist = filePath => {
  return fs.existsSync(filePath)
}

diskUtil.getModifiedFilesFromCacheFolder = (
  cacheFolder,
  absoluteProjectPath
) => {
  const paths = walk.sync(cacheFolder)
  const changes = paths.map(currentPath => {
    const relativePath = currentPath.replace(cacheFolder, '')

    const realPath = path.join(absoluteProjectPath, relativePath)
    // does file still exist?
    if (!fs.existsSync(realPath)) {
      return {
        path: realPath,
        status: 'deleted',
      }
    }
    // if path is not a folder, check the content
    if (fs.lstatSync(currentPath).isFile()) {
      const cacheContent = fs.readFileSync(currentPath)
      const content = fs.readFileSync(realPath)
      if (hash(cacheContent) !== hash(content)) {
        return {
          path: realPath,
          status: 'modified',
        }
      }
    }
    return {
      path: realPath,
      status: 'ok',
    }
  })
  return changes
}

diskUtil.deleteFolderRecursive = folderPath => {
  if (fs.existsSync(folderPath)) {
    fs.readdirSync(folderPath).forEach((file, index) => {
      const curPath = path.join(folderPath, file)
      if (fs.lstatSync(curPath).isDirectory()) {
        // recurse
        deleteFolderRecursive(curPath)
      } else {
        // delete file
        fs.unlinkSync(curPath)
      }
    })
    fs.rmdirSync(folderPath)
  }
}

diskUtil.changeEnvConfig = async ({ env, absoluteProjectPath }) => {
  const envConfigDescriptor = {
    name: KEYSTONE_ENV_CONFIG_PATH,
    content: {
      env,
    },
  }

  await writeFileToDisk(envConfigDescriptor, getKeystoneFolder('.'))

  // clean cache
  const cachePath = getCacheFolder(absoluteProjectPath)
  deleteFolderRecursive(cachePath)
  return envConfigDescriptor.content
}

// module.exports = {
//   writeFileToDisk,
//   readFileFromDisk,
//   deleteFileFromDisk,
//   getCacheFolder,
//   getKeystoneFolder,
//   getModifiedFilesFromCacheFolder,
//   isFileExist,
//   changeEnvConfig,
//   deleteFolderRecursive,
// }
