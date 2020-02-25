const path = require('path')
const { getLatestEnvDescriptor, updateDescriptor } = require('../../descriptor')
const { writeFileToGaia, listFilesFromGaia } = require('../../file/gaia')
const { deleteFileFromDisk, getCacheFolder } = require('../../file/disk')
const KeystoneError = require('../../error')
const { getProjects } = require('../../projects')
const { deepCopy } = require('../../utils')
const { PROJECTS_STORE } = require('../../constants')

const deleteFiles = async (
  userSession,
  { project, env, files, absoluteProjectPath }
) => {
  const { username } = userSession.loadUserData()

  const latestEnvDescriptor = await getLatestEnvDescriptor(userSession, {
    project,
    env,
    username,
  })
  const envDescriptorCloned = deepCopy(latestEnvDescriptor)

  envDescriptorCloned.content.files = latestEnvDescriptor.content.files.filter(
    envFile => !files.find(file => file === envFile.name)
  )

  if (
    envDescriptorCloned.content.files.length ===
    latestEnvDescriptor.content.files.length
  ) {
    throw new Error('No file to delete')
  }

  await updateDescriptor(userSession, {
    env,
    project,
    content: envDescriptorCloned.content,
    type: 'env',
  })
  const cacheFolder = await getCacheFolder(absoluteProjectPath)
  files.forEach(f => deleteFileFromDisk(path.join(cacheFolder, f)))
  return files
}

const deleteProject = async (userSession, { project }) => {
  const projects = await getProjects(userSession)
  const filteredProjects = projects.filter(p => p.name !== project)

  if (filteredProjects.length === projects.length) {
    throw new KeystoneError(
      'InvalidProjectName',
      `The project ${project} does not exist in your workspace.`
    )
  }
  writeFileToGaia(userSession, {
    content: JSON.stringify(filteredProjects),
    path: PROJECTS_STORE,
  })

  const projectFiles = (await listFilesFromGaia(userSession)).filter(f =>
    f.includes(project)
  )

  projectFiles.map(f => {
    console.log(`Deleted : ${f}`)
    userSession.deleteFile(f)
  })
}

module.exports = { deleteProject, deleteFiles }
