const { on, EVENTS } = require('@keystone/core/lib/pubsub')
const editorSpawn = require('child_process')
const path = require('path')
const fs = require('fs')
const daffy = require('daffy')
const { mergeContents } = require('@keystone/core/lib/descriptor')
const { getCacheFolder } = require('@keystone/core/lib/file')
const { getProjectConfigFolderPath } = require('../lib/blockstackLoader')

on(EVENTS.CONFLICT, ({ conflictFiles }) => {
  return new Promise(async resolve => {
    const cacheFolder = getCacheFolder(getProjectConfigFolderPath('.ksconfig'))

    const pathToFile = path.join(cacheFolder, `${conflictFiles[0].name}.merge`)

    // Get previous version of the content if exist
    const history = conflictFiles[0].history.find(
      x => x.version === conflictFiles.version - 1
    )

    // merge the two descriptor based on the previous version
    const { result, conflict } = mergeContents({
      left: conflictFiles[0].content,
      right: conflictFiles[1].content,
      base: history ? daffy.applyPatch(conflictFiles[0].content, history) : '',
    })

    if (conflict) {
      // write merge file in chache folder
      fs.writeFile(pathToFile, result, () => {})

      // open the file  in default editor
      const editorSpawned = editorSpawn.spawn(
        process.env.EDITOR || 'vi',
        [pathToFile],
        {
          stdio: 'inherit',
          detached: true,
        }
      )

      // on editor exit, return the new content (merged by the user)
      editorSpawned.on('close', () => {
        const newContent = fs.readFileSync(pathToFile)

        resolve(newContent)
      })
    }
    resolve(result)
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
