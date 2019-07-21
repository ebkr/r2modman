let tabs = [];
let searchBar = document.getElementById("searchBar");
let configOpen = false;

let actionableElements = [];

{
    let elements = document.getElementsByTagName("div");
    for (let i=0; i<elements.length; i++) {
        let x = elements[i];
        if (x.getAttribute("data-tabAssoc") != null) {
            tabs.push(x)
            let tab  = document.getElementById(x.getAttribute("data-tabAssoc"));
            console.log(tab)
            tab.addEventListener("click", ()=>{
                for (let j=0; j<tabs.length; j++) {
                    tabs[j].setAttribute("hidden", true);
                    let jTab = document.getElementById(tabs[j].getAttribute("data-tabAssoc"));
                    jTab.removeAttribute("selected")
                }
                x.removeAttribute("hidden");
                tab.setAttribute("selected", true)
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
})

window.ipcRenderer.on("configureClose", (event) => {
    configOpen = false;
})

function DrawInstalled() {
    document.getElementById("installedList").innerHTML = "";
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

        let missing = modManager.installed.HasMissingDependency(modList[i], modManager);
        console.log(missing);

        if (modManager.installed.HasUpdate(modList[i], modManager)) {
            // Alert
            let alert = document.createElement("button");
            alert.innerHTML = "&#8682;";
            alert.className = "buttonSmall";
            row.appendChild(alert);
            row.setAttribute("hasUpdate", "true")
            actionableElements.push(alert);
        } else if (!!missing) {
            // Warn
            let warn = document.createElement("button");
            warn.innerHTML = "&#9888;";
            warn.className = "buttonSmall";
            row.appendChild(warn);
            row.setAttribute("missingDependency", "true")
            actionableElements.push(warn);
        } else {
            row.setAttribute("valid", "");
            if (!modList[listPos].Enabled) {
                row.setAttribute("disabled", true);
            }
        }

        document.getElementById("installedList").appendChild(row);

    }
}

function DrawDownloadable() {
    document.getElementById("downloadList").innerHTML = "";
    let modList = modManager.available.GetModsByFilter(searchBar.value, modManager);
    for (let i=0; i<modList.length; i++) {

        // Store current iteration position
        let listPos = i;
        
        // Create row
        let row = document.createElement("div");
        row.className = "row";
        row.setAttribute("title", modList[i].versions[0].description);

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

        document.getElementById("downloadList").appendChild(row);

    }
}

searchBar.addEventListener("input", ()=>{
    DrawInstalled();
    DrawDownloadable();
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