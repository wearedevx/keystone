const debug = require('debug')('keystone:commands:env:create')

const { createEnv } = require('../../env')
const { getDescriptor } = require('../../descriptor')

const create = async (userSession, { project, env }) => {
  try {
    await createEnv(userSession, {
      env,
      projectDescriptor: await getDescriptor(userSession, {
        project,
        type: 'project',
      }),
    })
    debug(`▻ Environment name successfully created`)
  } catch (err) {
    console.log(err)
    debug(`▻ Environment creation failed : ${err}`)
  }
}

module.exports = {
  create,
}
