const { getLatestProjectDescriptor } = require('../../descriptor')

const listEnvironments = async (userSession, { project }) => {
  // Get latestProjectDescriptor
  const latestProjectDescriptor = await getLatestProjectDescriptor(
    userSession,
    {
      project,
    }
  )

  latestProjectDescriptor.content.env.forEach(env => {
    console.log(`▻ ${env}`)
  })
}

module.exports = { listEnvironments }
