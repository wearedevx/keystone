const debug = require('debug')('keystone:core:file')
const fileCache = require('./cache')

// TODO: we shouldn't need origin as encrypt already has that information: if set to true or false, it's the logged user else it's the blockstack id of another user
const writeFileToGaia = async (
  userSession,
  { path, origin = 'self', content, encrypt = true }
) => {
  debug('')
  debug(`Write file on gaia to ${path} with encrypt ${encrypt || 'no encrypt'}`)

  const cacheKey = `${origin}/${path}`

  await userSession.putFile(path, content, { encrypt })

  fileCache.put(cacheKey, JSON.parse(content))

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
    fileCache.put(cacheKey, fetchedFileJSON)
    return fetchedFileJSON
  }

  fileCache.put(cacheKey, fetchedFile)
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

module.exports = {
  writeFileToGaia,
  readFileFromGaia,
  deleteFilesFromGaia,
  getPubkey,
}
