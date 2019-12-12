const debug = require('debug')('keystone:core:file')
const NodeCache = require('node-cache')
const fs = require('fs')
const pathUtil = require('path')
const walk = require('walkdir')
const hash = require('object-hash')

const fsp = fs.promises

const { PUBKEY, KEYSTONE_HIDDEN_FOLDER } = require('../constants')

const fileCache = new NodeCache()

const writeFileToDisk = (fileDescriptor, absoluteProjectPath) => {
  try {
    const pathFile = pathUtil.join(absoluteProjectPath, fileDescriptor.name)
    const lastIndex = pathFile.lastIndexOf('/')
    const folder = pathFile.substring(0, lastIndex)

    debug('Write file to disk', pathFile)

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

const readFileFromDisk = async filename => {
  const buffer = await fsp.readFile(filename)
  return buffer.toString()
}

const deleteFileFromDisk = path => {
  debug('deleteFileFromDisk', path)
  return fs.unlinkSync(path)
}

// TODO: we shouldn't need origin as encrypt already has that information: if set to true or false, it's the logged user else it's the blockstack id of another user
const writeFileToGaia = async (
  userSession,
  { path, origin = 'self', content, encrypt = true }
) => {
  debug('')
  debug(`Write file on gaia to ${path} with encrypt ${encrypt || 'no encrypt'}`)

  const cacheKey = `${origin}/${path}`

  await userSession.putFile(path, content, { encrypt })

  fileCache.set(cacheKey, JSON.parse(content))

  return content
}

const readFileFromGaia = async (
  userSession,
  { path, origin = 'self', decrypt = true, json = true }
) => {
  const { username } = userSession.loadUserData()
  const cacheKey = `${origin}/${path}`
  const file = fileCache.get(cacheKey)

  if (origin === username) {
    origin = 'self'
  }

  debug('')

  if (file) {
    debug(`Get file from cache ${path} from ${origin}, cacheKey=${cacheKey}`)

    return file
  }
  debug(
    `Get file from gaia ${path} from ${origin || 'self'}, cacheKey=${cacheKey}`
  )

  const options = {
    decrypt,
  }

  if (origin !== 'self') {
    options.username = origin
  }

  const fetchedFile = await userSession.getFile(path, options)

  // debug(`${path} => `, fetchedFile)

  if (json) {
    const fetchedFileJSON = JSON.parse(fetchedFile)
    fileCache.set(cacheKey, fetchedFileJSON)
    return fetchedFileJSON
  }

  fileCache.set(cacheKey, fetchedFile)
  return fetchedFile
}

const deleteFilesFromGaia = (
  userSession,
  { path, opts = { wasSigned: true } }
) => {
  debug('deleteFilesFromGaia', path)
  return userSession.deleteFile(path, opts)
}

const getPubkey = async (userSession, { blockstackId }) => {
  const pubkeyFile = await readFileFromGaia(userSession, {
    decrypt: false,
    path: PUBKEY,
    origin: blockstackId,
    json: false,
  })
  if (pubkeyFile) {
    return pubkeyFile
  }
  throw new Error(
    `Keystone public application key not found on ${blockstackId}`
  )
}

const getCacheFolder = absoluteProjectPath => {
  const cacheFolder = `${absoluteProjectPath}/${KEYSTONE_HIDDEN_FOLDER}`
  if (!fs.existsSync(cacheFolder)) {
    fs.mkdirSync(cacheFolder)
  }

  return cacheFolder
}

const getModifiedFilesFromCacheFolder = (cacheFolder, absoluteProjectPath) => {
  const paths = walk.sync(cacheFolder)
  const changes = paths.map(path => {
    const relativePath = path.replace(cacheFolder, '')

    // TODO
    // Use path.join to create path os friendly.
    const realPath = `${absoluteProjectPath}${relativePath}`
    // does file still exist?
    if (!fs.existsSync(realPath)) {
      return {
        path: realPath,
        status: 'deleted',
      }
    }
    // if path is not a folder, check the content
    if (fs.lstatSync(path).isFile()) {
      const cacheContent = fs.readFileSync(path)
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

module.exports = {
  writeFileToDisk,
  readFileFromDisk,
  deleteFileFromDisk,

  writeFileToGaia,
  readFileFromGaia,
  deleteFilesFromGaia,

  getPubkey,
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
}
