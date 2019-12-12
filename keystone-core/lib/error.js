const constants = require('./constants')

const errorCodes = constants.ERROR_CODES

class KeystoneError extends Error {
  constructor(code, message, data) {
    const errorCode = errorCodes[code]
    if (errorCode) {
      super(message)
      this.code = errorCode
      this.name = 'KeystoneError'
      this.data = data
    } else {
      throw new Error(`Unknown Keystone error code: ${errorCode} - ${message}`)
    }
  }
}

module.exports = KeystoneError
