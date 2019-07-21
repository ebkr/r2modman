let toggle = true;

let json = JSON.parse(localStorage.getItem("configureMod"));
let children = document.getElementsByTagName("span");
for (let iter in children) {
    if (children[iter].innerHTML === undefined) { break; }
    let matches = children[iter].innerHTML.match(/{mod.[^}]*}/g)
    let str = "";
    let editingString = children[iter].innerHTML;
    for (let match in matches) {
        let foundSub = editingString.search(matches[match]);
        str += editingString.substr(0, foundSub)
        let exp = matches[match];
        let newExp = exp.substr(1, exp.length-2);
        let splitStr = newExp.split(".");
        let res = json;
        for (let i=1; i<splitStr.length; i++) {
            res = res[splitStr[i]]
        }
        str += res;
        editingString = editingString.substr(foundSub + exp.length)
    }
    children[iter].innerHTML = str;
}

function toggleText() {
    toggle = !toggle;
    if (toggle) {
        document.getElementById("toggle").innerHTML = "Disable Mod"
    } else {
        document.getElementById("toggle").innerHTML = "Enable Mod"
    }
}

toggle = json.Enabled;
if (!toggle) {
    document.getElementById("toggle").innerHTML = "Enable Mod"
}
document.getElementById("toggle").addEventListener("click", ()=>{
    window.ipcRenderer.send("toggleModUsability")
    toggleText()
})

if (json.URL.length === 0) {
    document.getElementById("view").remove();
} else {
    document.getElementById("view").addEventListener("click", ()=>{
        window.ipcRenderer.send("goToSite", json.URL);
    })
}

document.getElementById("remove").addEventListener("click", ()=>{
    if (confirm("Do you want to uninstall " + json.Name + "?")) {
        window.ipcRenderer.send("removeMod")
        toggleText()
    }
})