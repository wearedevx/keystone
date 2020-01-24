"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const hub_1 = require("./hub");
exports.connectToGaiaHub = hub_1.connectToGaiaHub;
exports.uploadToGaiaHub = hub_1.uploadToGaiaHub;
exports.BLOCKSTACK_GAIA_HUB_LABEL = hub_1.BLOCKSTACK_GAIA_HUB_LABEL;
// export { type GaiaHubConfig } from './hub'
const ec_1 = require("../encryption/ec");
const keys_1 = require("../keys");
const profileLookup_1 = require("../profiles/profileLookup");
const errors_1 = require("../errors");
const userSession_1 = require("../auth/userSession");
const authConstants_1 = require("../auth/authConstants");
const utils_1 = require("../utils");
const fetchUtil_1 = require("../fetchUtil");
const SIGNATURE_FILE_SUFFIX = '.sig';
/**
 * Fetch the public read URL of a user file for the specified app.
 * @param {String} path - the path to the file to read
 * @param {String} username - The Blockstack ID of the user to look up
 * @param {String} appOrigin - The app origin
 * @param {String} [zoneFileLookupURL=null] - The URL
 * to use for zonefile lookup. If falsey, this will use the
 * blockstack.js's [[getNameInfo]] function instead.
 * @return {Promise<string>} that resolves to the public read URL of the file
 * or rejects with an error
 */
async function getUserAppFileUrl(path, username, appOrigin, zoneFileLookupURL) {
    const profile = await profileLookup_1.lookupProfile(username, zoneFileLookupURL);
    let bucketUrl = null;
    if (profile.hasOwnProperty('apps')) {
        if (profile.apps.hasOwnProperty(appOrigin)) {
            const url = profile.apps[appOrigin];
            const bucket = url.replace(/\/?(\?|#|$)/, '/$1');
            bucketUrl = `${bucket}${path}`;
        }
    }
    return bucketUrl;
}
exports.getUserAppFileUrl = getUserAppFileUrl;
/**
 *
 *
 * @deprecated
 * #### v19 Use [[UserSession.encryptContent]].
 *
 * Encrypts the data provided with the app public key.
 * @param {String|Buffer} content - data to encrypt
 * @param {Object} [options=null] - options object
 * @param {String} options.publicKey - the hex string of the ECDSA public
 * key to use for encryption. If not provided, will use user's appPublicKey.
 * @return {String} Stringified ciphertext object
 */
async function encryptContent(content, options, caller) {
    const opts = Object.assign({}, options);
    if (!opts.publicKey) {
        const privateKey = (caller || new userSession_1.UserSession()).loadUserData().appPrivateKey;
        opts.publicKey = keys_1.getPublicKeyFromPrivate(privateKey);
    }
    const cipherObject = await ec_1.encryptECIES(opts.publicKey, content);
    return JSON.stringify(cipherObject);
}
exports.encryptContent = encryptContent;
/**
 *
 * @deprecated
 * #### v19 Use [[UserSession.decryptContent]].
 *
 * Decrypts data encrypted with `encryptContent` with the
 * transit private key.
 * @param {String|Buffer} content - encrypted content.
 * @param {Object} [options=null] - options object
 * @param {String} options.privateKey - the hex string of the ECDSA private
 * key to use for decryption. If not provided, will use user's appPrivateKey.
 * @return {String|Buffer} decrypted content.
 */
function decryptContent(content, options, caller) {
    const opts = Object.assign({}, options);
    if (!opts.privateKey) {
        opts.privateKey = (caller || new userSession_1.UserSession()).loadUserData().appPrivateKey;
    }
    try {
        const cipherObject = JSON.parse(content);
        return ec_1.decryptECIES(opts.privateKey, cipherObject);
    }
    catch (err) {
        if (err instanceof SyntaxError) {
            throw new Error('Failed to parse encrypted content JSON. The content may not '
                + 'be encrypted. If using getFile, try passing { decrypt: false }.');
        }
        else {
            throw err;
        }
    }
}
exports.decryptContent = decryptContent;
/* Get the gaia address used for servicing multiplayer reads for the given
 * (username, app) pair.
 * @private
 * @ignore
 */
async function getGaiaAddress(app, username, zoneFileLookupURL, caller) {
    const opts = normalizeOptions({ app, username, zoneFileLookupURL }, caller);
    let fileUrl;
    if (username) {
        fileUrl = await getUserAppFileUrl('/', opts.username, opts.app, opts.zoneFileLookupURL);
    }
    else {
        if (!caller) {
            caller = new userSession_1.UserSession();
        }
        const gaiaHubConfig = await caller.getOrSetLocalGaiaHubConnection();
        fileUrl = await hub_1.getFullReadUrl('/', gaiaHubConfig);
    }
    const matches = fileUrl.match(/([13][a-km-zA-HJ-NP-Z0-9]{26,35})/);
    if (!matches) {
        throw new Error('Failed to parse gaia address');
    }
    return matches[matches.length - 1];
}
/**
 * @param {Object} [options=null] - options object
 * @param {String} options.username - the Blockstack ID to lookup for multi-player storage
 * @param {String} options.app - the app to lookup for multi-player storage -
 * defaults to current origin
 *
 * @ignore
 */
function normalizeOptions(options, caller) {
    const opts = Object.assign({}, options);
    if (opts.username) {
        if (!opts.app) {
            caller = caller || new userSession_1.UserSession();
            if (!caller.appConfig) {
                throw new errors_1.InvalidStateError('Missing AppConfig');
            }
            opts.app = caller.appConfig.appDomain;
        }
        if (!opts.zoneFileLookupURL) {
            caller = caller || new userSession_1.UserSession();
            if (!caller.appConfig) {
                throw new errors_1.InvalidStateError('Missing AppConfig');
            }
            if (!caller.store) {
                throw new errors_1.InvalidStateError('Missing store UserSession');
            }
            const sessionData = caller.store.getSessionData();
            // Use the user specified coreNode if available, otherwise use the app specified coreNode. 
            const configuredCoreNode = sessionData.userData.coreNode || caller.appConfig.coreNode;
            if (configuredCoreNode) {
                opts.zoneFileLookupURL = `${configuredCoreNode}${authConstants_1.NAME_LOOKUP_PATH}`;
            }
        }
    }
    return opts;
}
/**
 * @deprecated
 * #### v19 Use [[UserSession.getFileUrl]] instead.
 *
 * @param {String} path - the path to the file to read
 * @returns {Promise<string>} that resolves to the URL or rejects with an error
 */
async function getFileUrl(path, options, caller) {
    const opts = normalizeOptions(options, caller);
    let readUrl;
    if (opts.username) {
        readUrl = await getUserAppFileUrl(path, opts.username, opts.app, opts.zoneFileLookupURL);
    }
    else {
        const gaiaHubConfig = await (caller || new userSession_1.UserSession()).getOrSetLocalGaiaHubConnection();
        readUrl = await hub_1.getFullReadUrl(path, gaiaHubConfig);
    }
    if (!readUrl) {
        throw new Error('Missing readURL');
    }
    else {
        return readUrl;
    }
}
exports.getFileUrl = getFileUrl;
/* Handle fetching the contents from a given path. Handles both
 *  multi-player reads and reads from own storage.
 * @private
 * @ignore
 */
async function getFileContents(path, app, username, zoneFileLookupURL, forceText, caller) {
    const opts = { app, username, zoneFileLookupURL };
    const readUrl = await getFileUrl(path, opts, caller);
    const response = await fetchUtil_1.fetchPrivate(readUrl);
    if (!response.ok) {
        throw await utils_1.getBlockstackErrorFromResponse(response, `getFile ${path} failed.`);
    }
    const contentType = response.headers.get('Content-Type');
    if (forceText || contentType === null
        || contentType.startsWith('text')
        || contentType === 'application/json') {
        return response.text();
    }
    else {
        return response.arrayBuffer();
    }
}
/* Handle fetching an unencrypted file, its associated signature
 *  and then validate it. Handles both multi-player reads and reads
 *  from own storage.
 * @private
 * @ignore
 */
async function getFileSignedUnencrypted(path, opt, caller) {
    // future optimization note:
    //    in the case of _multi-player_ reads, this does a lot of excess
    //    profile lookups to figure out where to read files
    //    do browsers cache all these requests if Content-Cache is set?
    const sigPath = `${path}${SIGNATURE_FILE_SUFFIX}`;
    try {
        const [fileContents, signatureContents, gaiaAddress] = await Promise.all([
            getFileContents(path, opt.app, opt.username, opt.zoneFileLookupURL, false, caller),
            getFileContents(sigPath, opt.app, opt.username, opt.zoneFileLookupURL, true, caller),
            getGaiaAddress(opt.app, opt.username, opt.zoneFileLookupURL, caller)
        ]);
        if (!fileContents) {
            return fileContents;
        }
        if (!gaiaAddress) {
            throw new errors_1.SignatureVerificationError('Failed to get gaia address for verification of: '
                + `${path}`);
        }
        if (!signatureContents || typeof signatureContents !== 'string') {
            throw new errors_1.SignatureVerificationError('Failed to obtain signature for file: '
                + `${path} -- looked in ${path}${SIGNATURE_FILE_SUFFIX}`);
        }
        let signature;
        let publicKey;
        try {
            const sigObject = JSON.parse(signatureContents);
            signature = sigObject.signature;
            publicKey = sigObject.publicKey;
        }
        catch (err) {
            if (err instanceof SyntaxError) {
                throw new Error('Failed to parse signature content JSON '
                    + `(path: ${path}${SIGNATURE_FILE_SUFFIX})`
                    + ' The content may be corrupted.');
            }
            else {
                throw err;
            }
        }
        const signerAddress = keys_1.publicKeyToAddress(publicKey);
        if (gaiaAddress !== signerAddress) {
            throw new errors_1.SignatureVerificationError(`Signer pubkey address (${signerAddress}) doesn't`
                + ` match gaia address (${gaiaAddress})`);
        }
        else if (!ec_1.verifyECDSA(fileContents, publicKey, signature)) {
            throw new errors_1.SignatureVerificationError('Contents do not match ECDSA signature: '
                + `path: ${path}, signature: ${path}${SIGNATURE_FILE_SUFFIX}`);
        }
        else {
            return fileContents;
        }
    }
    catch (err) {
        // For missing .sig files, throw `SignatureVerificationError` instead of `DoesNotExist` error.
        if (err instanceof errors_1.DoesNotExist && err.message.indexOf(sigPath) >= 0) {
            throw new errors_1.SignatureVerificationError('Failed to obtain signature for file: '
                + `${path} -- looked in ${path}${SIGNATURE_FILE_SUFFIX}`);
        }
        else {
            throw err;
        }
    }
}
/* Handle signature verification and decryption for contents which are
 *  expected to be signed and encrypted. This works for single and
 *  multiplayer reads. In the case of multiplayer reads, it uses the
 *  gaia address for verification of the claimed public key.
 * @private
 * @ignore
 */
async function handleSignedEncryptedContents(caller, path, storedContents, app, privateKey, username, zoneFileLookupURL) {
    const appPrivateKey = privateKey || caller.loadUserData().appPrivateKey;
    const appPublicKey = keys_1.getPublicKeyFromPrivate(appPrivateKey);
    let address;
    if (username) {
        address = await getGaiaAddress(app, username, zoneFileLookupURL, caller);
    }
    else {
        address = keys_1.publicKeyToAddress(appPublicKey);
    }
    if (!address) {
        throw new errors_1.SignatureVerificationError('Failed to get gaia address for verification of: '
            + `${path}`);
    }
    let sigObject;
    try {
        sigObject = JSON.parse(storedContents);
    }
    catch (err) {
        if (err instanceof SyntaxError) {
            throw new Error('Failed to parse encrypted, signed content JSON. The content may not '
                + 'be encrypted. If using getFile, try passing'
                + ' { verify: false, decrypt: false }.');
        }
        else {
            throw err;
        }
    }
    const signature = sigObject.signature;
    const signerPublicKey = sigObject.publicKey;
    const cipherText = sigObject.cipherText;
    const signerAddress = keys_1.publicKeyToAddress(signerPublicKey);
    if (!signerPublicKey || !cipherText || !signature) {
        throw new errors_1.SignatureVerificationError('Failed to get signature verification data from file:'
            + ` ${path}`);
    }
    else if (signerAddress !== address) {
        throw new errors_1.SignatureVerificationError(`Signer pubkey address (${signerAddress}) doesn't`
            + ` match gaia address (${address})`);
    }
    else if (!ec_1.verifyECDSA(cipherText, signerPublicKey, signature)) {
        throw new errors_1.SignatureVerificationError('Contents do not match ECDSA signature in file:'
            + ` ${path}`);
    }
    else if (typeof (privateKey) === 'string') {
        const decryptOpt = { privateKey };
        return caller.decryptContent(cipherText, decryptOpt);
    }
    else {
        return caller.decryptContent(cipherText);
    }
}
/**
 * Retrieves the specified file from the app's data store.
 * @param {String} path - the path to the file to read
 * @returns {Promise} that resolves to the raw data in the file
 * or rejects with an error
 */
async function getFile(path, options, caller) {
    const defaults = {
        decrypt: true,
        verify: false,
        username: null,
        app: utils_1.getGlobalObject('location', { returnEmptyObject: true }).origin,
        zoneFileLookupURL: null
    };
    const opt = Object.assign({}, defaults, options);
    if (!caller) {
        caller = new userSession_1.UserSession();
    }
    // in the case of signature verification, but no
    //  encryption expected, need to fetch _two_ files.
    if (opt.verify && !opt.decrypt) {
        return getFileSignedUnencrypted(path, opt, caller);
    }
    const storedContents = await getFileContents(path, opt.app, opt.username, opt.zoneFileLookupURL, !!opt.decrypt, caller);
    if (storedContents === null) {
        return storedContents;
    }
    else if (opt.decrypt && !opt.verify) {
        if (typeof storedContents !== 'string') {
            throw new Error('Expected to get back a string for the cipherText');
        }
        if (typeof (opt.decrypt) === 'string') {
            const decryptOpt = { privateKey: opt.decrypt };
            return caller.decryptContent(storedContents, decryptOpt);
        }
        else {
            return caller.decryptContent(storedContents);
        }
    }
    else if (opt.decrypt && opt.verify) {
        if (typeof storedContents !== 'string') {
            throw new Error('Expected to get back a string for the cipherText');
        }
        let decryptionKey;
        if (typeof (opt.decrypt) === 'string') {
            decryptionKey = opt.decrypt;
        }
        return handleSignedEncryptedContents(caller, path, storedContents, opt.app, decryptionKey, opt.username, opt.zoneFileLookupURL);
    }
    else if (!opt.verify && !opt.decrypt) {
        return storedContents;
    }
    else {
        throw new Error('Should be unreachable.');
    }
}
exports.getFile = getFile;
/** @ignore */
class FileContentLoader {
    constructor(content) {
        this.content = content;
    }
    getContentType() {
        if (typeof this.content === 'string') {
            return 'text/plain; charset=utf-8';
        }
        else if (typeof Blob !== 'undefined' && this.content instanceof Blob && this.content.type) {
            return this.content.type;
        }
        else {
            return 'application/octet-stream';
        }
    }
    async loadContent() {
        if (typeof this.content === 'string') {
            return this.content;
        }
        else if (ArrayBuffer.isView(this.content)) {
            return Buffer.from(this.content.buffer);
        }
        else if (typeof Blob !== 'undefined' && this.content instanceof Blob) {
            const reader = new FileReader();
            const readPromise = new Promise((resolve, reject) => {
                reader.onerror = (err) => {
                    reject(err);
                };
                reader.onload = () => {
                    const arrayBuffer = reader.result;
                    resolve(Buffer.from(arrayBuffer));
                };
                reader.readAsArrayBuffer(this.content);
            });
            const result = await readPromise;
            return result;
        }
        const typeName = Object.prototype.toString.call(this.content);
        throw new Error(`Unsupported content object type: ${typeName}`);
    }
    load() {
        if (this.loadedData === undefined) {
            this.loadedData = this.loadContent();
        }
        return this.loadedData;
    }
}
/**
 * Stores the data provided in the app's data store to to the file specified.
 * @param {String} path - the path to store the data in
 * @param {String|Buffer} content - the data to store in the file
 * @return {Promise} that resolves if the operation succeed and rejects
 * if it failed
 */
async function putFile(path, content, options, caller) {
    const contentLoader = new FileContentLoader(content);
    const defaults = {
        encrypt: true,
        sign: false,
        contentType: ''
    };
    const opt = Object.assign({}, defaults, options);
    let { contentType } = opt;
    if (!contentType) {
        contentType = contentLoader.getContentType();
    }
    if (!caller) {
        caller = new userSession_1.UserSession();
    }
    // First, let's figure out if we need to get public/private keys,
    //  or if they were passed in
    let privateKey = '';
    let publicKey = '';
    if (opt.sign) {
        if (typeof (opt.sign) === 'string') {
            privateKey = opt.sign;
        }
        else {
            privateKey = caller.loadUserData().appPrivateKey;
        }
    }
    if (opt.encrypt) {
        if (typeof (opt.encrypt) === 'string') {
            publicKey = opt.encrypt;
        }
        else {
            if (!privateKey) {
                privateKey = caller.loadUserData().appPrivateKey;
            }
            publicKey = keys_1.getPublicKeyFromPrivate(privateKey);
        }
    }
    // In the case of signing, but *not* encrypting,
    //   we perform two uploads. So the control-flow
    //   here will return there.
    if (!opt.encrypt && opt.sign) {
        const contentData = await contentLoader.load();
        const signatureObject = ec_1.signECDSA(privateKey, contentData);
        const signatureContent = JSON.stringify(signatureObject);
        const gaiaHubConfig = await caller.getOrSetLocalGaiaHubConnection();
        try {
            const fileUrls = await Promise.all([
                hub_1.uploadToGaiaHub(path, contentData, gaiaHubConfig, contentType),
                hub_1.uploadToGaiaHub(`${path}${SIGNATURE_FILE_SUFFIX}`, signatureContent, gaiaHubConfig, 'application/json')
            ]);
            return fileUrls[0];
        }
        catch (error) {
            const freshHubConfig = await caller.setLocalGaiaHubConnection();
            const fileUrls = await Promise.all([
                hub_1.uploadToGaiaHub(path, contentData, freshHubConfig, contentType),
                hub_1.uploadToGaiaHub(`${path}${SIGNATURE_FILE_SUFFIX}`, signatureContent, freshHubConfig, 'application/json')
            ]);
            return fileUrls[0];
        }
    }
    // In all other cases, we only need one upload.
    let contentForUpload;
    if (opt.encrypt && !opt.sign) {
        const contentData = await contentLoader.load();
        contentForUpload = await encryptContent(contentData, { publicKey });
        contentType = 'application/json';
    }
    else if (opt.encrypt && opt.sign) {
        const contentData = await contentLoader.load();
        const cipherText = await encryptContent(contentData, { publicKey });
        const signatureObject = ec_1.signECDSA(privateKey, cipherText);
        const signedCipherObject = {
            signature: signatureObject.signature,
            publicKey: signatureObject.publicKey,
            cipherText
        };
        contentForUpload = JSON.stringify(signedCipherObject);
        contentType = 'application/json';
    }
    else {
        contentForUpload = content;
    }
    const gaiaHubConfig = await caller.getOrSetLocalGaiaHubConnection();
    try {
        return await hub_1.uploadToGaiaHub(path, contentForUpload, gaiaHubConfig, contentType);
    }
    catch (error) {
        const freshHubConfig = await caller.setLocalGaiaHubConnection();
        const file = await hub_1.uploadToGaiaHub(path, contentForUpload, freshHubConfig, contentType);
        return file;
    }
}
exports.putFile = putFile;
/**
 * Deletes the specified file from the app's data store.
 * @param path - The path to the file to delete.
 * @param options - Optional options object.
 * @param options.wasSigned - Set to true if the file was originally signed
 * in order for the corresponding signature file to also be deleted.
 * @returns Resolves when the file has been removed or rejects with an error.
 */
async function deleteFile(path, options, caller) {
    if (!caller) {
        caller = new userSession_1.UserSession();
    }
    const gaiaHubConfig = await caller.getOrSetLocalGaiaHubConnection();
    const opts = Object.assign({}, options);
    if (opts.wasSigned) {
        // If signed, delete both the content file and the .sig file
        try {
            await hub_1.deleteFromGaiaHub(path, gaiaHubConfig);
            await hub_1.deleteFromGaiaHub(`${path}${SIGNATURE_FILE_SUFFIX}`, gaiaHubConfig);
        }
        catch (error) {
            const freshHubConfig = await caller.setLocalGaiaHubConnection();
            await hub_1.deleteFromGaiaHub(path, freshHubConfig);
            await hub_1.deleteFromGaiaHub(`${path}${SIGNATURE_FILE_SUFFIX}`, gaiaHubConfig);
        }
    }
    else {
        try {
            await hub_1.deleteFromGaiaHub(path, gaiaHubConfig);
        }
        catch (error) {
            const freshHubConfig = await caller.setLocalGaiaHubConnection();
            await hub_1.deleteFromGaiaHub(path, freshHubConfig);
        }
    }
}
exports.deleteFile = deleteFile;
/**
 * Get the app storage bucket URL
 * @param {String} gaiaHubUrl - the gaia hub URL
 * @param {String} appPrivateKey - the app private key used to generate the app address
 * @returns {Promise} That resolves to the URL of the app index file
 * or rejects if it fails
 */
function getAppBucketUrl(gaiaHubUrl, appPrivateKey) {
    return hub_1.getBucketUrl(gaiaHubUrl, appPrivateKey);
}
exports.getAppBucketUrl = getAppBucketUrl;
/**
 * Loop over the list of files in a Gaia hub, and run a callback on each entry.
 * Not meant to be called by external clients.
 * @param {GaiaHubConfig} hubConfig - the Gaia hub config
 * @param {String | null} page - the page ID
 * @param {number} callCount - the loop count
 * @param {number} fileCount - the number of files listed so far
 * @param {function} callback - the callback to invoke on each file.  If it returns a falsey
 *  value, then the loop stops.  If it returns a truthy value, the loop continues.
 * @returns {Promise} that resolves to the number of files listed.
 * @private
 * @ignore
 */
async function listFilesLoop(caller, hubConfig, page, callCount, fileCount, callback) {
    if (callCount > 65536) {
        // this is ridiculously huge, and probably indicates
        // a faulty Gaia hub anyway (e.g. on that serves endless data)
        throw new Error('Too many entries to list');
    }
    hubConfig = hubConfig || await caller.getOrSetLocalGaiaHubConnection();
    let response;
    try {
        const pageRequest = JSON.stringify({ page });
        const fetchOptions = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': `${pageRequest.length}`,
                Authorization: `bearer ${hubConfig.token}`
            },
            body: pageRequest
        };
        response = await fetchUtil_1.fetchPrivate(`${hubConfig.server}/list-files/${hubConfig.address}`, fetchOptions);
        if (!response.ok) {
            throw await utils_1.getBlockstackErrorFromResponse(response, 'ListFiles failed.');
        }
    }
    catch (error) {
        // If error occurs on the first call, perform a gaia re-connection and retry.
        // Same logic as other gaia requests (putFile, getFile, etc).
        if (callCount === 0) {
            const freshHubConfig = await caller.setLocalGaiaHubConnection();
            return listFilesLoop(caller, freshHubConfig, page, callCount + 1, 0, callback);
        }
        throw error;
    }
    const responseText = await response.text();
    const responseJSON = JSON.parse(responseText);
    const entries = responseJSON.entries;
    const nextPage = responseJSON.page;
    if (entries === null || entries === undefined) {
        // indicates a misbehaving Gaia hub or a misbehaving driver
        // (i.e. the data is malformed)
        throw new Error('Bad listFiles response: no entries');
    }
    let entriesLength = 0;
    for (let i = 0; i < entries.length; i++) {
        // An entry array can have null entries, signifying a filtered entry and that there may be
        // additional pages
        if (entries[i] !== null) {
            entriesLength++;
            const rc = callback(entries[i]);
            if (!rc) {
                // callback indicates that we're done
                return fileCount + i;
            }
        }
    }
    if (nextPage && entries.length > 0) {
        // keep going -- have more entries
        return listFilesLoop(caller, hubConfig, nextPage, callCount + 1, fileCount + entriesLength, callback);
    }
    else {
        // no more entries -- end of data
        return fileCount + entriesLength;
    }
}
/**
 * List the set of files in this application's Gaia storage bucket.
 * @param {function} callback - a callback to invoke on each named file that
 * returns `true` to continue the listing operation or `false` to end it
 * @return {Promise} that resolves to the number of files listed, or rejects with an error.
 */
function listFiles(callback, caller) {
    caller = caller || new userSession_1.UserSession();
    return listFilesLoop(caller, null, null, 0, 0, callback);
}
exports.listFiles = listFiles;
//# sourceMappingURL=index.js.map