const debug = require('debug')('keystone:core:file')
const fileCache = require('./cache')
const fs = require('fs')
const pathUtil = require('path')
const walk = require('walkdir')
const hash = require('object-hash')
const path = require('path')

const fsp = fs.promises
const KeystoneError = require('../error')
const {
  PUBKEY,
  KEYSTONE_HIDDEN_FOLDER,
  KEYSTONE_CONFIG_PATH,
  KEYSTONE_ENV_CONFIG_PATH,
} = require('../constants')

const writeFileToDisk = (fileDescriptor, absoluteProjectPath) => {
  const pathFile = pathUtil.join(absoluteProjectPath, fileDescriptor.name)
  const lastIndex = pathFile.lastIndexOf(pathUtil.sep)
  const folder = pathFile.substring(0, lastIndex)
  debug('Write file to disk', pathFile)

  if (folder) fs.mkdirSync(folder, { recursive: true })

  // if JSON object, stringify
  let { content } = fileDescriptor
  if (typeof content === 'object') content = JSON.stringify(content)

  fs.writeFile(pathFile, content, err => {
    if (err) throw new Error(err)
  })

  return fileDescriptor
}

const readFileFromDisk = async filename => {
  const buffer = await fsp.readFile(filename)
  return buffer.toString('utf-8')
}

const deleteFileFromDisk = path => {
  debug('deleteFileFromDisk', path)
  return fs.unlinkSync(path)
}

const getCacheFolder = absoluteProjectPath => {
  const cacheFolder = path.join(
    getKeystoneFolder(absoluteProjectPath),
    `/cache/`
  )
  if (!fs.existsSync(cacheFolder)) {
    fs.mkdirSync(cacheFolder, { recursive: true })
  }

  return cacheFolder
}

const getKeystoneFolder = absoluteProjectPath => {
  const keystoneFolder = `${absoluteProjectPath}/${KEYSTONE_HIDDEN_FOLDER}`
  if (!fs.existsSync(keystoneFolder)) {
    const gitignorePath = path.join(absoluteProjectPath, '.gitignore')
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

const isFileExist = filePath => {
  return fs.existsSync(filePath)
}

const getModifiedFilesFromCacheFolder = (cacheFolder, absoluteProjectPath) => {
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

const deleteFolderRecursive = function(folderPath) {
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

async function changeEnvConfig({ env, absoluteProjectPath }) {
  const envConfigDescriptor = {
    name: KEYSTONE_ENV_CONFIG_PATH,
    content: {
      env,
    },
  }

  await writeFileToDisk(
    envConfigDescriptor,
    getKeystoneFolder(absoluteProjectPath)
  )

  // clean cache
  const cachePath = getCacheFolder(absoluteProjectPath)
  deleteFolderRecursive(cachePath)
  return envConfigDescriptor.content
}

function resetLocalFiles(absoluteProjectPath, confirm = false) {
  const modifiedFiles = getModifiedFilesFromCacheFolder(
    getCacheFolder(absoluteProjectPath),
    absoluteProjectPath
  ).filter(f => f.status !== 'ok')

  if (modifiedFiles.length === 0)
    throw new KeystoneError('NoPendingModification')

  if (!confirm)
    throw new KeystoneError('PendingModification', '', modifiedFiles)

  modifiedFiles.forEach(f => {
    const filename = f.path.replace(path.join(absoluteProjectPath, '/'), '')
    const previousContent = fs
      .readFileSync(path.join(getCacheFolder(absoluteProjectPath), filename))
      .toString()
    fs.writeFileSync(f.path, previousContent)
  })
  return modifiedFiles
}

module.exports = {
  resetLocalFiles,
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
