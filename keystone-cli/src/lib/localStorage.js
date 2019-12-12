const cache = {};

function getItem(key) {
    return cache[key];
}

function setItem(key, item) {
    cache[key] = item;
}

function removeItem(key) {
    cache[key] = null;
}

function getAll(){
    return cache;
}

function createLocalStorage() {
    const storage = {
        getItem,
        setItem,
        removeItem,
        getAll
    };
   
    return storage;
}

function updateLocalStorage(storage, data = {}){
    Object.keys(data).forEach(key => {
        storage.setItem(key, data[key]);
    });
    return storage
}

module.exports = {
    createLocalStorage,
    updateLocalStorage
}