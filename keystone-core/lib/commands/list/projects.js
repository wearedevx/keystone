const { readFileFromGaia } = require('../../file')
const { PROJECTS_STORE } = require('../../constants')
const listProjects = async userSession => {
  const projects = await readFileFromGaia(userSession, { path: PROJECTS_STORE })
  if (!projects) {
    console.log("You don't have any project in your workspace")
    return null
  }
  console.log(`Projects: ${projects.length} found`)

  const printableProjects = projects.map(
    ({ name, createdBy, pendingInvite }) => ({
      name,
      createdBy,
      invitation: pendingInvite ? 'waiting for access' : 'ok',
    })
  )

  printableProjects.forEach(project => {
    console.log(
      ` > ${project.name.split('/')[0]}/${(
        project.name.split('/')[1]
      )}`
    )
  })
}

module.exports = {
  listProjects,
}
