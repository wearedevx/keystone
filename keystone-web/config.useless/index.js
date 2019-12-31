// const getEnv = env => {
//   if (env === 'development') {
//     return import('./env.development')
//   }
//   if (env === 'staging') {
//     return import('./env.staging')
//   }
//   if (env === 'production') {
//     return import('./env.production')
//   }
// }

// export default async env => {
//   const currentEnv = await getEnv(env)

//   Object.keys(currentEnv).forEach(key => {
//     process.env[key] = currentEnv[key]
//   })
// }
