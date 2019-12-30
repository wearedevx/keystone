const { on, EVENTS } = require('@keystone/core/lib/pubsub')
const editorSpawn = require('child_process')
const path = require('path')
const fs = require('fs')
const daffy = require('daffy')
const inquirer = require('inquirer')
const chalk = require('chalk')
const treeify = require('treeify')
const { mergeContents } = require('@keystone/core/lib/descriptor')
const { getCacheFolder } = require('@keystone/core/lib/file')
const {
  updateFilesInEnvDesciptor,
} = require('@keystone/core/lib/commands/push')
const { getProjectConfigFolderPath } = require('../lib/blockstackLoader')

on(EVENTS.CONFLICT, ({ conflictFiles }) => {
  return new Promise(async resolve => {
    if (typeof conflictFiles[0].content === 'string') {
      const cacheFolder = getCacheFolder(
        getProjectConfigFolderPath('.ksconfig')
      )

      const pathToFile = path.join(
        cacheFolder,
        `${conflictFiles[0].name}.merge`
      )

      // Get previous version of the content if exist
      const history = conflictFiles[0].history.find(
        x => x.version === conflictFiles.version - 1
      )

      // merge the two descriptor based on the previous version
      const { result, conflict } = mergeContents({
        left: conflictFiles[0].content,
        right: conflictFiles[1].content,
        base: history
          ? daffy.applyPatch(conflictFiles[0].content, history)
          : '',
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
          const stringContent = newContent.toString()
          resolve(stringContent)
        })
      } else {
        resolve(result)
      }
    } else {
      const choices = [conflictFiles[0], conflictFiles[1]]
        .reduce((files, descriptor) => {
          files.push(
            ...(descriptor.content.files || descriptor.content.members)
          )
          return files
        }, [])
        .reduce((unqFiles, file) => {
          if (!unqFiles.find(f => f.name === file.name)) {
            unqFiles.push({ name: file.name, value: file, checked: true })
          }
          return unqFiles
        }, [])

      const user1Items = (
        conflictFiles.right.content.files || conflictFiles.right.content.members
      ).map(i => i.blockstack_id || i.name)

      const user2Items = (
        conflictFiles[0].content.files || conflictFiles[0].content.members
      ).map(i => i.blockstack_id || i.name)

      const itemsByOwner = {
        [chalk.green(conflictFiles.right.author)]: user1Items,
        [chalk.green(conflictFiles[0].author)]: user2Items,
      }
      console.log('\x1Bc')
      console.log(treeify.asTree(itemsByOwner, true))

      const { items } = await inquirer.prompt([
        {
          type: 'checkbox',
          name: 'items',
          message: `Which files you want to keep from the env ?`,
          choices,
        },
      ])
      resolve(items)
    }
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
//        [0]: descriptorsWithMaxVersion[0].content,
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
