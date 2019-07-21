// All of the Node.js APIs are available in the preload process.
// It has the same sandbox as a Chrome extension.
window.addEventListener('DOMContentLoaded', () => {
  
})

const { ipcRenderer } = require('electron');
function init() {
    // add global variables to your web page
    window.isElectron = true
    window.ipcRenderer = ipcRenderer
}

init();
