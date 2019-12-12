/**
 * Deeply copies an object
 * @param {object} obj Object to Copy
 * @return {object}
 */
const deepCopy = obj => {
  if (obj === null || typeof obj !== 'object') return obj
  if (obj instanceof Date) {
    const copy = new Date()
    copy.setTime(obj.getTime())
    return copy
  }
  if (obj instanceof Array) {
    const copy = []
    for (let i = 0, len = obj.length; i < len; i += 1) {
      copy[i] = deepCopy(obj[i])
    }
    return copy
  }
  if (obj instanceof Object) {
    const copy = {}

    Object.keys(obj).forEach(attr => {
      if (Object.prototype.hasOwnProperty.call(obj, attr))
        copy[attr] = deepCopy(obj[attr])
    })
    return copy
  }

  throw new Error('Unable to copy obj this object.')
}

module.exports = {
  deepCopy,
}
