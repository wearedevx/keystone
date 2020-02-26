const gitDiff = require('git-diff')
const { getCacheFolder, readFileFromDisk } = require('../../file')
const path = require('path')
const diff = async (userSession, { absoluteProjectPath, filePath, file }) => {
  const cacheFolder = getCacheFolder(absoluteProjectPath)

  let previousContent
  let currentContent
  try {
    previousContent = await readFileFromDisk(path.join(cacheFolder, filePath))
  } catch (err) {
    throw new Error("The file hasn't been modified.")
  }
  try {
    currentContent = await readFileFromDisk(path.join(file))
  } catch (err) {
    throw new Error(`The file ${file} does not exist.`)
  }

  const diffOutput = await gitDiff(previousContent, currentContent, {
    color: true,
  })
  if (!diffOutput) throw new Error("The file hasn't been modified.")

  return diffOutput
}

module.exports = diff
