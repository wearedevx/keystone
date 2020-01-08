const debug = require('debug')('keystone:core:file')
const fileCache = require('./cache')
const fs = require('fs')
const pathUtil = require('path')
const walk = require('walkdir')
const hash = require('object-hash')
const Path = require('path')

const fsp = fs.promises

const {
  PUBKEY,
  KEYSTONE_HIDDEN_FOLDER,
  KEYSTONE_CONFIG_PATH,
  KEYSTONE_ENV_CONFIG_PATH,
} = require('../constants')

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

const getCacheFolder = absoluteProjectPath => {
  const cacheFolder = `${getKeystoneFolder(absoluteProjectPath)}/cache/`
  if (!fs.existsSync(cacheFolder)) {
    fs.mkdirSync(cacheFolder)
  }

  return cacheFolder
}

const getKeystoneFolder = absoluteProjectPath => {
  const keystoneFolder = `${absoluteProjectPath}/${KEYSTONE_HIDDEN_FOLDER}`
  if (!fs.existsSync(keystoneFolder)) {
    fs.appendFileSync('.gitignore', `\n${KEYSTONE_HIDDEN_FOLDER}/`)
    fs.mkdirSync(keystoneFolder)
  }

  return keystoneFolder
}

const isFileExist = filePath => {
  return fs.existsSync(filePath)
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

const deleteFolderRecursive = function(path) {
  if (fs.existsSync(path)) {
    fs.readdirSync(path).forEach((file, index) => {
      const curPath = Path.join(path, file)
      if (fs.lstatSync(curPath).isDirectory()) {
        // recurse
        deleteFolderRecursive(curPath)
      } else {
        // delete file
        fs.unlinkSync(curPath)
      }
    })
    fs.rmdirSync(path)
  }
}

async function changeEnvConfig({ env, absoluteProjectPath }) {
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

module.exports = {
  writeFileToDisk,
  readFileFromDisk,
  deleteFileFromDisk,
  getCacheFolder,
  getKeystoneFolder,
  getModifiedFilesFromCacheFolder,
  isFileExist,
  changeEnvConfig,
  deleteFolderRecursive,
}
