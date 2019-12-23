const { merge } = require('three-way-merge')
const daffy = require('daffy')
const { conflictedDescriptors, conflictLessDescriptors } = require('./dataset')

describe('Manage conflicts', () => {
  it('should pop a conflict up', async () => {
    const left = conflictedDescriptors.left.content
    const right = conflictedDescriptors.right.content
    const base = conflictedDescriptors.base.content
    const merged = merge(left, base, right)

    expect(merged.conflict).toBeTruthy()
  })

  it('should not trow a conflict and merge', async () => {
    const left = conflictLessDescriptors.left.content
    const right = conflictLessDescriptors.right.content
    const base = conflictLessDescriptors.base.content
    const merged = merge(left, base, right)

    expect(merged.conflict).toBeFalsy()
  })
})

describe('Get back in history', () => {
  it('should patch a text', () => {
    const previousContent = 'this is the previous content \n second line'
    const newContent = 'this is the new content \nSecond line'
    const patch = daffy.createPatch(newContent, previousContent)

    const t = daffy.applyPatch(newContent, patch)

    expect(t).toContain('this is the previous content \n second line')
  })
})
