const fs = require('fs')
// use file API with promises - more elegant.
const fsp = fs.promises

const createFolder = async ({ path }) => {
    try {
        await fsp.mkdir(path, {recursive: true})
    } catch (error) {
        throw error
    }
}

const write = async ({ path, filename, content }) => {
    try {
        await fsp.access(path, fs.constants.F_OK)
    } catch (error) {
        try {
            await createFolder({path})
        } catch (error) {
            throw error
        }
        
    } finally {
        const filepath = `${path}/${filename}`
        await fsp.writeFile(filepath, content)
    }
}

const read = async ({ path = '', filename }) => {
    const filepath = `${path}${filename}`
    try {
        const buffer = await fsp.readFile(filepath)
        //all files stored for CLI usage should be in JSON
        return JSON.parse(buffer.toString())
    } catch (error) {
        throw error
    }
}

const del = async ({ path = '', filename }) => {
    const filepath = `${path}${filename}`
    try {
        await fsp.unlink(filepath)
    } catch (error) {
        throw error
    }
}

module.exports = {
    write,  
    read,
    del,
  createFolder
}
