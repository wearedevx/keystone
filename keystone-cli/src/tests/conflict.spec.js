// const { merge } = require('three-way-merge-lines')
// const daffy = require('daffy')
// const inquirer = require('inquirer')
// const chalk = require('chalk')
// const treeify = require('treeify')

// const {
//   conflictedDescriptors,
//   conflictLessDescriptors,
//   conflictedEnvDescriptors,
// } = require('./dataset')

// describe('Manage conflicts', () => {
//   it('should pop a conflict up', async () => {
//     const left = conflictedDescriptors.left.content
//     const right = conflictedDescriptors.right.content
//     const base = conflictedDescriptors.base.content
//     const merged = merge(left, base, right)

//     expect(merged.conflict).toBeTruthy()
//   })

//   it('should not trow a conflict and merge', async () => {
//     const left = conflictLessDescriptors.left.content
//     const right = conflictLessDescriptors.right.content
//     const base = conflictLessDescriptors.base.content
//     const merged = merge(left, base, right)

//     expect(merged.conflict).toBeFalsy()
//   })
// })

// describe('Get back in history', () => {
//   it('should patch a text', () => {
//     const previousContent = 'this is the previous content \n second line'
//     const newContent = 'this is the new content \nSecond line'
//     const patch = daffy.createPatch(newContent, previousContent)

//     const t = daffy.applyPatch(newContent, patch)

//     expect(t).toContain('this is the previous content \n second line')
//   })
// })

// describe('Manage merge between descriptor with array content', () => {
//   fit('should prompt the user to choose the files he/she wants to keep', async () => {
//     const choices = [
//       conflictedEnvDescriptors.left,
//       conflictedEnvDescriptors.right,
//     ]
//       .reduce((files, descriptor) => {
//         files.push(...(descriptor.content.files || descriptor.content.members))
//         return files
//       }, [])
//       .reduce((unqFiles, file) => {
//         if (!unqFiles.find(f => f.name === file.name)) {
//           unqFiles.push({ name: file.name, value: file, checked: true })
//         }
//         return unqFiles
//       }, [])

//     const user1Items = (
//       conflictedEnvDescriptors.right.content.files ||
//       conflictedEnvDescriptors.right.content.members
//     ).map(i => i.blockstack_id || i.name)

//     const user2Items = (
//       conflictedEnvDescriptors.left.content.files ||
//       conflictedEnvDescriptors.left.content.members
//     ).map(i => i.blockstack_id || i.name)

//     const itemsByOwner = {
//       [chalk.green(conflictedEnvDescriptors.right.author)]: user1Items,
//       [chalk.green(conflictedEnvDescriptors.left.author)]: user2Items,
//     }

//     console.log(treeify.asTree(itemsByOwner, true))

//     const { items } = await inquirer.prompt([
//       {
//         type: 'checkbox',
//         name: 'items',
//         message: `Which files you want to keep from the env ?`,
//         choices,
//       },
//     ])
//     console.log(items)
//   }, 10000000)
// })
