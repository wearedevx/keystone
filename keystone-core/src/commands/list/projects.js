const { readFileFromGaia } = require('../../file')
const { PROJECTS_STORE } = require('../../constants')

const listProjects = async userSession => {
  const projects = await readFileFromGaia(userSession, { path: PROJECTS_STORE })

  console.log(`Projects: ${projects.length} found`)

  const printableProjects = projects.map(
    ({ name, createdBy, pendingInvite }) => ({
      name,
      createdBy,
      invitation: pendingInvite ? 'waiting for access' : 'ok',
    })
  )

  console.table(printableProjects)
}

module.exports = {
  listProjects,
}
