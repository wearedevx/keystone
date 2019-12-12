module.exports = require('@oclif/command')

process
  .on('uncaughtException', err => {
    console.log('UNCAUGHT§', err)
  })
  .on('unhandledRejection', err => {
    console.log('UNCAUGHT§', err)
  })
