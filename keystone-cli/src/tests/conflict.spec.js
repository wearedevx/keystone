const { merge } = require('three-way-merge')
const { conflictedDescriptors, conflictLessDescriptors } = require('./dataset')
describe('Manage conflicts', () => {
  it('should pop a conflict up', async () => {
    const left = conflictedDescriptors.left.content
    const right = conflictedDescriptors.right.content
    const base = conflictedDescriptors.base.content
    const merged = merge(left, base, right)

    console.log(merged.conflict)

    console.log(merged.joinedResults())

    expect(merged.conflict).toBeTruthy()
  })

  it('should not trow a conflict and merge', async () => {
    const left = conflictLessDescriptors.left.content
    const right = conflictLessDescriptors.right.content
    const base = conflictLessDescriptors.base.content
    const merged = merge(left, base, right)

    console.log(merged.conflict)

    console.log(merged.joinedResults())

    expect(merged.conflict).toBeFalsy()
  })
})
