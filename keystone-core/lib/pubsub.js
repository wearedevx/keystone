// @ts-check

const CONSTANTS = require('./constants')

const listeners = []
/**
 * Execute sequantially functions.
 * @param {*} arrPromises
 * @param {*} payload
 */
const executeSequentially = (
  arrPromises,
  payload,
  timer = 1000,
  results = [],
  indc = 0
) => {
  return new Promise(async resolve => {
    if (arrPromises[indc]) {
      results.push(await arrPromises[indc](payload))

      if (arrPromises.length <= indc) {
        resolve(results)
      } else {
        setTimeout(() => {
          resolve(
            executeSequentially(arrPromises, payload, timer, results, indc + 1)
          )
        }, timer)
      }
    } else {
      resolve(results)
    }
  })
}

/**
 * Emit an event.
 * Cal lall listeners of the event name.
 * @param {*} userId User id
 * @param {*} eventName Event name
 * @param {*} payload Event paylaod
 * @param {*} asynchronous Default true. Indicate if need to execute all listeners asynchronously or not.
 * @param {*} timer Default 0. Timer to wait before each listeners calls if synchronous execution.
 */
async function emit(eventName, payload = null, asynchronous = true, timer = 0) {
  if (process.env.NODE_ENV !== 'test') console.log('Event fired:', eventName)
  let res

  // try {
  //   track(userId, eventName, payload)
  // } catch (err) {
  //   console.error(err)
  // }

  if (!listeners[eventName]) return Promise.resolve([])

  if (asynchronous) {
    res = Promise.all(listeners[eventName].map(f => f(payload)))
  } else {
    res = await executeSequentially(listeners[eventName], payload, timer)
  }

  return res
}

/**
 * Declare a function to execute on a event name.
 * @param {*} eventName
 * @param {*} func
 */
function on(eventName, func) {
  if (!listeners[eventName]) listeners[eventName] = []
  listeners[eventName].push(func)
}

const singleton = {
  EVENTS: CONSTANTS.EVENTS,
  on,
  emit,
}

module.exports = singleton
