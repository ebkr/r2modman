window.ipcRenderer.send("enableBorderless")
window.ipcRenderer.send("resize", document.body.clientWidth, document.body.clientHeight)

fetch("https://thunderstore.io/api/v1/package/", {
    method: 'get'
}).then((response) => {
    return response.text()
}).then((response) => {
    localStorage.setItem("modList", response)
    window.ipcRenderer.send("splashFinished", response)
})