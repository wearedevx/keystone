const { on, EVENTS } = require('@keystone/core/lib/pubsub')
const editorSpawn = require('child_process')
const path = require('path')
const fs = require('fs')
const { getCacheFolder } = require('@keystone/core/lib/file')
const { getProjectConfigFolderPath } = require('../lib/blockstackLoader')

on(EVENTS.CONFLICT, ({ conflictFiles }) => {
  return new Promise(async resolve => {
    console.log(conflictFiles)
    // inquirer avec conflictFiles
    // const conflictResolved = result of inquirer
    // resolve(conflictResolved)
    const cacheFolder = getCacheFolder(getProjectConfigFolderPath('.ksconfig'))

    const pathToFile = path.join(cacheFolder, `${conflictFiles[0].name}.merge`)
    console.log('TCL: pathToFile', pathToFile)

    fs.writeFile(pathToFile, conflictFiles[0].content, () => {})

    const editorSpawned = editorSpawn.spawn('vim', [pathToFile], {
      stdio: 'inherit',
      detached: true,
    })

    editorSpawned.on('close', () => {
      const newContent = fs.readFileSync(pathToFile)

      resolve(newContent)

      console.log('FINISH')
    })
  })
})

// if (typeof descriptorsWithMaxVersion[0].content === 'string') {
//   const previousVersion = daffy.applyPatch(
//     descriptorsWithMaxVersion[0].content,
//     descriptorsWithMaxVersion[0].history.find(x => x === maxVersion - 1)
//   )
//   let finalContent = previousVersion

//   for (let i = 0; i < descriptorsWithMaxVersion.length; i += 1) {
//     if (descriptorsWithMaxVersion[1]) {
//       const { result } = mergeContents({
//         left: descriptorsWithMaxVersion[0].content,
//         right: descriptorsWithMaxVersion[1].content,
//         base: previousVersion,
//       })
//       finalContent = result
//       descriptorsWithMaxVersion.splice(0, 2)
//       descriptorsWithMaxVersion.splice(0, 0, finalContent)
//     } else {
//       descriptorsWithMaxVersion[0].content = finalContent
//     }
//   }
//   console.log(`\n File ${newDescriptors[0].name} in conflict ! \n`)
//   return descriptorsWithMaxVersion[0]
// }
// return descriptorsWithMaxVersion[0]
// }
// const choices = descriptorsWithMaxVersion
// .map(descriptor => descriptor.content.files || descriptor.content.members)
// .reduce((acc, curr) => {
//   if (!acc.find(f => f.name === curr.name)) {
//     acc.push(curr)
//   }
//   return acc
// }, [])

// const { items } = await inquirer.prompt([
// {
//   type: 'checkbox',
//   name: 'items',
//   message: `Which files you want to keep from the env ?`,
//   choices,
// },
// ])
// }
