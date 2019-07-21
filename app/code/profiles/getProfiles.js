window.ipcRenderer.send("getProfiles");
window.ipcRenderer.on("getProfiles", (event, args) => {
    document.getElementById("profiles").innerHTML = "";
    for (let i=0; i<args.length; i++) {
        let opt = document.createElement("option");
        opt.innerHTML = args[i];
        document.getElementById("profiles").appendChild(opt);
    }
})

document.getElementById("selectProfile").onclick = function() {
    if (document.getElementById("profiles").value.length > 0) {
        // select
        window.ipcRenderer.send("selectedProfile", document.getElementById("profiles").value);
    }
}

document.getElementById("createProfile").addEventListener("click", function() {
    window.ipcRenderer.send("promptProfileCreation");
})

document.getElementById("deleteProfile").addEventListener("click", function() {
    let value = document.getElementById("profiles").value
    if (value.length > 0) {
        localStorage.setItem("selectedProfile", value)
        window.ipcRenderer.send("deleteProfile", value);
    }
})