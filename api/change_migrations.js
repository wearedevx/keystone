/*
 * This small piece of javascript aims at transforming 
 * sequential migratiions into a timestamp ones.
 *
 * NOTE: This doesn't handle the database schema_migrations table,
 * it will have to be done by other means
 */ 
const fs = require('fs');
const fsp = require('fs/promises');
const path = require('path');
const { execSync } = require('child_process');

const MIGRATION_FOLDER = path.join(process.cwd(), 'db', 'migrations');

const SUFFIXES = {
	UP: 'up.sql',
	DOWN: 'down.sql'
};


// Main 
(async () => {
	const files = await fsp.readdir(MIGRATION_FOLDER)
	console.log('[START]')

	console.log('')
	console.log("RENAMING SEQUENTIAL FILES -----------")

	const newFilenames = await changeFilenames(files)

	console.log('')
	console.log("RENAMING BAD DOWNS ------------------")

	await renameBadDowns(newFilenames)

	console.log('')
	console.log('[DONE]')
})() 

/**
	* After the first pass of renaming, the down migrations
	* might have a different timestamp. This function 
	* makes it so up and down migrations have the same one.
	*
	* @param {Array<string>} filenames
	*
	* @return {Promise<>}
	*/
async function renameBadDowns(filenames) {
	const splits = splitAll(filenames) 
	const pairs = byPair(splits)

	const promises = Object.entries(pairs).map(async ([migrationName, prefixes]) => {
		if (!prefixes.up) {
			console.warn(`WARNING: No up for ${migrationName}`)
			return
		}

		if (!prefixes.down) {
			console.warn(`WARNING: No down for ${migrationName}`)
			return
		}

		if (prefixes.up != prefixes.down) {
			const orig = `${prefixes.down}_${migrationName}.${SUFFIXES.DOWN}`
			const dest = `${prefixes.up}_${migrationName}.${SUFFIXES.DOWN}`

			mv(orig, dest)
		}
	})

	await Promise.all(promises)
}

/**
	* @typedef {Object} Migration
	* @property {string} prefix
	* @property {string} migrationName
	* @property {string} suffix
	*/

/**
	* Splits a migration name into a more manageable object
	* @param {string} filename a migration filename
	* @return {Migration}
	*/
function splitName(filename) {
	const indexOfUnderscore = filename.indexOf('_')
	const indexOfExt = filename.indexOf('.')
	const prefix = filename.substr(0, indexOfUnderscore)
	const suffix = filename.substr(indexOfExt + 1)
	const migrationName = filename.substring(indexOfUnderscore + 1, indexOfExt)

	return {
		prefix,
		migrationName,
		suffix
	}
}

/**
	* Splits all given filenames into more manageable objects.
	* @param {Array<string>} filenames
	* @return {Array<Migration>}
	*/
function splitAll(filenames) {
	return filenames.map(f => splitName(f))
}

/**
	* Orders migrations by up and down pairs
	* @param {Array<Migration>}
	* @return {Object<string, {up: string, down: string}>}
	*/
function byPair(filenames) {
	const map = {}

	filenames.forEach(f => {
		if (!map.hasOwnProperty(f.migrationName)) {
			map[f.migrationName] = {}
		}

		if (f.suffix === SUFFIXES.UP) {
			map[f.migrationName].up = f.prefix
		}

		if (f.suffix === SUFFIXES.DOWN) {
			map[f.migrationName].down = f.prefix
		}
	})

	return map
}

/**
	* Change all migration filenames from sequential ones
	* into timestamp ones
	*
	* @param {Array<string>}
	*
	* @return {Promise<Array<string>>}
	*/
async function changeFilenames(filenames) {
	const promises = filenames.map(async filename => changeFilename(filename))

	return await Promise.all(promises)
}

/**
	* Change the given migration filename into 
	* a timestamp one
	*
	* @param {string}
	*
	* @return {Promise<string>}
	*/
async function changeFilename(filename) {
	const newFilename = await getNewName(MIGRATION_FOLDER, filename)

	const orig = path.join(MIGRATION_FOLDER, filename)
	const dest = path.join(MIGRATION_FOLDER, newFilename)

	mv(orig, dest)

	return dest
}

/**
	* Gives the new name of a file using its modification time
	* and the sequence number.
	* It sets the hours at 9 am, and uses the sequence number as minutes,
	* in order to preserve the original order of migrations.
	*
	* @param {string} folder  The location of the migrations
	* @param {string} name    The current name
	*
	* @return {Promise<string>} The timestamp-prefixed name
	*/
async function getNewName(folder, name) {
	const indexOfUnderscore = name.indexOf('_')
	const seqNumber = name.substr(0, indexOfUnderscore)
	const migrationName = name.substr(indexOfUnderscore + 1)

	if (seqNumber.length !== 6) {
		return name
	}

	const filepath = path.join(folder, name)

	const stat = await fsp.stat(filepath)
	const lastModified = stat.mtime

	if (isNaN(lastModified)) {
		return name
	}

	lastModified.setHours(9)
	lastModified.setMinutes(parseInt(seqNumber))

	const prefix = formatDate(lastModified)

	return `${prefix}_${migrationName}`
}

/**
	* Formats a date as migration time stamp
	* @param {Date} date
	* @return {string}
	*/
function formatDate(date) {
	const year = padNumber(date.getFullYear(), 4)
	const month = padNumber(date.getMonth() + 1, 2)
	const day = padNumber(date.getDate(), 2)
	const hours = padNumber(date.getHours(), 2)
	const minutes = padNumber(date.getMinutes(), 2)
	const seconds = padNumber(date.getSeconds(), 2)

	return `${year}${month}${day}${hours}${minutes}${seconds}`
}

/**
	* Pads a number with zero on the left side
	* @param {number} number
	* @param {number} pad How many zero to pad the number with
	*/
function padNumber(number, pad) {
	return number.toString().padStart(pad, '0')
}

/**
	* Moves a file using git for proper versionning
	*
	* @param {string} orig
	* @param {string} dest
	*/
function mv(orig, dest) {
	if (orig != dest) {
		execSync(`git mv ${orig} ${dest}`)
		console.log(`${orig} -> ${dest}`);
	}
}
