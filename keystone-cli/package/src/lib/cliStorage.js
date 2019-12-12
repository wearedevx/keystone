const fs = require('fs')
// use file API with promises - more elegant.
const fsp = fs.promises

const createFolder = async ({ path }) => {
    try {
        await fsp.mkdir(path, {recursive: true})
    } catch (error) {
        console.log("unable to create folder", error)
    }
}

const write = async ({ path, filename, content }) => {
    try {
        await fsp.access(path, fs.constants.F_OK)
    } catch (error) {
        console.log("error", error)
        await createFolder({path})
    } finally {
        const filepath = `${path}/${filename}`
        await fsp.writeFile(filepath, content)
    }
}

const read = async ({ path, filename }) => {
    const filepath = `${path}/${filename}`
    try {
        const buffer = await fsp.readFile(filepath)
        //all files stored for CLI usage should be in JSON
        return JSON.parse(buffer.toString())
    } catch (error) {
        console.log("Can't read file",error)
    }
}

module.exports = {
    write,  
    read
}