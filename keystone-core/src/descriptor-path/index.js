const getEnvPath = ({ env, project, blockstackId }) => {
  if (!project) {
    throw new Error('Need project!')
  }
  if (!env) {
    throw new Error('Need env!')
  }

  return `${project}/${env}/${blockstackId}.json`
}

const getFilePath = ({ filename, blockstackId, env, project }) => {
  if (!project) {
    throw new Error('Need project!')
  }
  if (!env) {
    throw new Error('Need env!')
  }
  if (!filename) {
    throw new Error('Need filename!')
  }

  return `${project}/${env}/${filename}/${blockstackId}.json`
}

const getMembersPath = ({ blockstackId, env, project }) => {
  if (!project) {
    throw new Error('Need project!')
  }

  const directories = [project]

  if (env) {
    directories.push(env)
  }

  directories[directories.length - 1] += '-members'
  directories.push(`${blockstackId}.json`)

  return directories.join('/')
}

const getProjectPath = ({ blockstackId, project }) => {
  if (!project) {
    throw new Error('Need project!')
  }

  return `${project}/${blockstackId}.json`
}

const getPath = ({ env, project, blockstackId, type, filename }) => {
  let path

  switch (type) {
    case 'project':
      path = getProjectPath({
        blockstackId,
        project,
      })
      break

    case 'env':
      path = getEnvPath({
        blockstackId,
        env,
        project,
      })
      break

    case 'file':
      path = getFilePath({
        blockstackId,
        env,
        project,
        filename,
      })
      break

    case 'members':
      path = getMembersPath({
        blockstackId,
        env,
        project,
      })
      break

    default:
      throw new Error(`Upload file with unknown type: ${type}`)
  }

  return path
}

const changeBlockstackId = (path, blockstackId) => {
  const indexOfSlash = path.lastIndexOf('/')
  return `${path.substring(0, indexOfSlash)}/${blockstackId}.json`
}

module.exports = {
  getPath,
  changeBlockstackId,
}
