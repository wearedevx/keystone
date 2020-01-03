const debug = require('debug')('keystone:core:env')
const fs = require('fs')
const del = require('del')
const path = require('path')
const {
  getLatestProjectDescriptor,
  getLatestEnvDescriptor,
} = require('../descriptor')
const { KEYSTONE_CONFIG_PATH, KEYSTONE_HIDDEN_FOLDER } = require('../constants')

const checkoutEnv = async (
  userSession,
  { project, env, absoluteProjectPath }
) => {
  const projectDescriptor = await getLatestProjectDescriptor(userSession, {
    project,
    type: 'project',
  })

  // Retrieve updated project descriptor
  const envFound = projectDescriptor.content.env.find(
    envObject => envObject === env
  )

  if (envFound) {
    await getLatestEnvDescriptor(userSession, {
      project,
      env,
      type: 'env',
    })
    const configFile = JSON.parse(
      fs.readFileSync(path.join(absoluteProjectPath, KEYSTONE_CONFIG_PATH))
    )
    configFile.env = env
    fs.writeFileSync(
      path.join(absoluteProjectPath, KEYSTONE_CONFIG_PATH),
      JSON.stringify(configFile)
    )

    // clean cache
    const cachePath = path.join(absoluteProjectPath, KEYSTONE_HIDDEN_FOLDER)
    console.log('TCL: cachePath', cachePath)
    const deletedPaths = await del([cachePath], { force: true })
    console.log('Deleted files and directories:\n', deletedPaths.join('\n'))
    // return configFile
  }
  throw new Error(`The environment ${env} is not defined in this project`)
}

module.exports = {
  checkoutEnv,
}
