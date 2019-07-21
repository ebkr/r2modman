let tabAssoc = {};
let searchBar = document.getElementById("searchBar");
let configOpen = false;
let actionableElements = [];

{
    let elements = document.getElementById("nav").children;
    for (let i=0; i<elements.length; i++) {
        let x = elements[i];
        let idSub = x.getAttribute("id").substr(1);
        if (document.getElementById("d" + idSub)) {
            x.addEventListener("click", ()=>{
                for (let j=0; j<elements.length; j++) {
                    elements[j].removeAttribute("selected");
                    if (document.getElementById("d" + elements[j].getAttribute("id").substr(1))) {
                        let nestAssoc = document.getElementById("d" + elements[j].getAttribute("id").substr(1));
                        nestAssoc.setAttribute("hidden", "true");
                    }
                }
                x.setAttribute("selected", "true");
                document.getElementById("d" + idSub).removeAttribute("hidden");
                if (idSub === "1" || idSub === "2") {
                    document.getElementById("topBar").removeAttribute("style");
                } else {
                    document.getElementById("topBar").style.display = "none";
                }
                searchBar.value = "";
                onInput();
            });
        }
    }
}

var modManager = new ModManager();

window.ipcRenderer.send("getInstalledMods", "");
window.ipcRenderer.on("getInstalledMods", (event, args) => {
    for (let ae in actionableElements) {
        actionableElements[ae].remove();
    }
    actionableElements = [];
    modManager.FilesToInstalled(args)
    DrawInstalled();
    DrawDownloadable();
});

window.ipcRenderer.on("configureClose", (event) => {
    configOpen = false;
});

function DrawInstalled() {
    // Remove existing children
    let using = document.getElementById("d1");
    // Fully clear
    using.innerHTML = "";
    // Initialise mod rows
    let modList = modManager.installed.GetModsByFilter(searchBar.value, modManager);
    for (let i=0; i<modList.length; i++) {
        // Store current iteration position
        let listPos = i;

        // Create row
        let row = document.createElement("div");
        row.className = "row";

        // Create image
        let image = document.createElement("img");
        image.src = modList[i].Icon;
        image.className = "listImage";

        // Create text positional elements
        let textOuter = document.createElement("tr");
        let textInner = document.createElement("td");
        let text = document.createElement("span");
        text.innerHTML = modList[i].Name;

        // Settings
        let settings = document.createElement("button");
        settings.innerHTML = "Configure";
        settings.className = "buttonSmall";
        settings.addEventListener("click", ()=>{
            if (configOpen) { return }
            configOpen = true;
            localStorage.setItem("configureMod", JSON.stringify(modList[listPos]));
            window.ipcRenderer.send("configureMod", modList[listPos]);
        })
        actionableElements.push(settings);

        // Set element parents
        row.appendChild(image);
        row.appendChild(textOuter);
        textOuter.appendChild(textInner);
        textInner.appendChild(text);
        row.appendChild(settings);

        let missing = modManager.installed.HasMissingDependency(modList[listPos], modManager);

        if (modManager.installed.HasUpdate(modList[listPos], modManager)) {
            // HasUpdate
            let alertObj = document.createElement("button");
            alertObj.innerHTML = "&#8682;";
            alertObj.className = "buttonSmall";
            row.appendChild(alertObj);
            row.setAttribute("problematic", "true")
            actionableElements.push(alertObj);
            alertObj.addEventListener("click", ()=>{
                let ml = modManager.available.GetModsByFilter("", modManager);
                for (let iter in ml) {
                    let tsmod = ml[iter];
                    if (tsmod.uuid4 === modList[listPos].Uuid4) {
                        if (confirm("Do you want to update " + modList[listPos].Name + "?")) {
                            let res = modManager.available.GetLatestVersion(modList[listPos].Uuid4, modManager);
                            if (res) {
                                console.log(modManager.available.GetRootFromUUID4(modList[listPos].Uuid4, modManager));
                                window.ipcRenderer.send("downloadMod", res, modManager.available.GetRootFromUUID4(modList[listPos].Uuid4, modManager));
                            }
                        }
                    }
                }
            })
        } else if (!!missing) {
            // MissingDependency
            let warn = document.createElement("button");
            warn.innerHTML = "&#9888;";
            warn.className = "buttonSmall";
            row.appendChild(warn);
            row.setAttribute("problematic", "true")
            actionableElements.push(warn);
            warn.addEventListener("click", ()=>{
                let ml = modManager.available.GetModsByFilter("", modManager);
                for (let iter in ml) {
                    let tsmod = ml[iter];
                    if (tsmod.full_name === missing) {
                        if (confirm("Install dependency: " + tsmod.name + "?")) {
                            window.ipcRenderer.send("downloadMod", tsmod.versions[0], tsmod);
                        }
                    }
                }
            })
        } else {
            row.setAttribute("valid", "");
            if (!modList[listPos].Enabled) {
                row.setAttribute("disabled", true);
            }
        }

        using.appendChild(row);
    }
    using.appendChild(document.createElement("br"));
}

function DrawDownloadable() {
    // Remove existing children
    let using = document.getElementById("d2");
    // Fully clear
    using.innerHTML = "";
    let modList = modManager.available.GetModsByFilter(searchBar.value, modManager);
    for (let i=0; i<modList.length; i++) {

        // Store current iteration position
        let listPos = i;
        
        // Create row
        let row = document.createElement("div");
        row.className = "row";
        row.setAttribute("title", modList[i].versions[0].description);
        row.setAttribute("online", "true");

        // Create image
        let image = document.createElement("img");
        image.src = modList[i].versions[0].icon;
        image.className = "listImage";

        // Create text positional elements
        let textOuter = document.createElement("tr");
        let textInner = document.createElement("td");
        let text = document.createElement("span");
        text.innerHTML = modList[i].name;

        // Version Select
        let storedVersions = {};
        let selection = document.createElement("select");
        for (let version=0; version<modList[i].versions.length; version++) {
            let v = modList[i].versions[version];
            let opt = document.createElement("option");
            opt.innerText = v.version_number;
            selection.appendChild(opt);
            storedVersions[v.version_number] = v;
        }

        // Download
        let download = document.createElement("button");
        download.innerHTML = "&#9660;";
        download.className = "buttonSmall";
        download.addEventListener("click", ()=>{
            window.ipcRenderer.send("downloadMod", storedVersions[selection.value], modList[listPos]);
        })
        actionableElements.push(download);

        // Set element parents
        row.appendChild(image);
        row.appendChild(textOuter);
        textOuter.appendChild(textInner);
        textInner.appendChild(text);
        row.appendChild(selection)
        row.appendChild(download);

        using.appendChild(row);

    }
    using.appendChild(document.createElement("br"));
}

function onInput() {
    DrawInstalled();
    DrawDownloadable();
}
searchBar.addEventListener("input", ()=>{
    onInput();
})

// Force refresh from server
window.ipcRenderer.on("refresh", ()=>{
    window.location.reload();
});

// User refresh
document.getElementById("refresh").addEventListener("click", ()=>{
    if (configOpen) { 
        alert("Can't refresh whilst mod config is open!")
        return 
    }
    window.location.reload();
})

document.getElementById("addLocal").addEventListener("click", ()=>{
    window.ipcRenderer.send("addLocalMod", "");
})

document.getElementById("gameStarter").addEventListener("click", ()=>{
    window.ipcRenderer.send("playRoR2", "");
})

document.getElementById("exportProfile").addEventListener("click", ()=>{
    window.ipcRenderer.send("exportProfile");
});

document.getElementById("importProfile").addEventListener("click", ()=>{
    window.ipcRenderer.send("importProfile");
});

window.ipcRenderer.on("isAppAssociated", (event, associated) => {
    if (associated) {
        document.getElementById("associate").style.display = "none";
    }
});

window.ipcRenderer.on("isAssociatedSuccess", (event, associated) => {
    console.log("isassoc:", associated);
    if (associated) {
        alert("You can now use the 'Install with Mod Manager' button.");
        document.getElementById("associate").style.display = "none";
    } else {
        alert("This feature requires admin privileges.");
    }
});

document.getElementById("associate").addEventListener("click", ()=>{
    window.ipcRenderer.send("associateHandler");
});

document.getElementById("disableAll").addEventListener("click", ()=>{
    window.ipcRenderer.send("disableAll");
});

document.getElementById("enableAll").addEventListener("click", ()=>{
    window.ipcRenderer.send("enableAll");
});

document.getElementById("deleteAll").addEventListener("click", ()=>{
    window.ipcRenderer.send("deleteAll");
});

window.ipcRenderer.send("isAppAssociated");