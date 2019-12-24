const debug = require('debug')('keystone:core:file')
const fileCache = require('./cache')
const fs = require('fs')
const pathUtil = require('path')
const walk = require('walkdir')
const hash = require('object-hash')

const fsp = fs.promises

const { PUBKEY, KEYSTONE_HIDDEN_FOLDER } = require('../constants')

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
  const cacheFolder = `${absoluteProjectPath}/${KEYSTONE_HIDDEN_FOLDER}`
  if (!fs.existsSync(cacheFolder)) {
    fs.mkdirSync(cacheFolder)
  }

  return cacheFolder
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

module.exports = {
  writeFileToDisk,
  readFileFromDisk,
  deleteFileFromDisk,
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
  isFileExist,
}
