const { getLatestEnvDescriptor } = require('../../descriptor')

const listAllFiles = userSession => {
  userSession.listFiles(file => {
    console.log(file)
    return true
  })
}

const listEnvFiles = async (userSession, { project, env }) => {
  console.log("TCL: listEnvFiles ->  project, env",  project, env)
  const envDescriptor = await getLatestEnvDescriptor(userSession, {
    project,
    env,
  })

  if (envDescriptor.content.files.length === 0) {
    console.log(`No files in ${project} for env ${env}`)
  } else {
    envDescriptor.content.files.forEach(file => {
      console.log(`▻ ${file.name}`)
    })
  }
}

module.exports = {
  listAllFiles,
  listEnvFiles,
}
