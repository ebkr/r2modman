const fs = require("fs-extra");
const path = require('path');

class Configuration {

    gamePath            = ""    // String
    keepZipDownloads    = false // Boolean
    currentItems        = []    // Array of strings

    // If valid then allow game to start.
    GamePathValid() {
        if (this.gamePath.length == 0) {
            return false;
        }
        if (fs.existsSync(this.gamePath)) {
            return true;
        }
        return false;
    }
    
    // Record moved items
    AddInstalledItem(path) {
        this.currentItems.push(path);
    }

    // Remove moved items
    RemoveInstalledMods() {
        for (let i=0; i<this.currentItems.length; i++) {
            fs.removeSync(this.currentItems[i]);
        }
        this.currentItems = [];
    }

    // Update config file
    Update(configPath) {
        fs.writeFileSync(configPath, JSON.stringify(this, null, 4), (e)=>{
            console.log(e);
        })
    }

    // Create or load values from /mods/profiles/config.json
    constructor(configPath) {
        if (fs.existsSync(configPath)) {
            let data = fs.readFileSync(configPath, "utf8");
            let json = JSON.parse(data);
            for (let x in json) {
                this[x] = json[x];
            }
        } else {
            this.Update(configPath);
        }
    }

    // Install mod list to RoR2 folder
    InstallMods(modList) {
        this.RemoveInstalledMods()
        for (let i=0; i<modList.length; i++) {
            if (modList[i].Enabled) {
                let fileList = installMod(modList[i], this.gamePath);
                for (let key in fileList) {
                    if (fileList[key].length > 0) {
                        let copyDir = path.join(this.gamePath, "BepInEx", key, modList[i].FullName)
                        fs.mkdirsSync(copyDir);
                        for (let fileIter=0; fileIter<fileList[key].length; fileIter++) {
                            fs.copy(fileList[key][fileIter], path.join(copyDir, path.basename(fileList[key][fileIter])))
                        }
                        this.AddInstalledItem(copyDir);
                    }
                }
            }
        }
    }
    
}

exports.Configuration = Configuration;

function installBepInEx(mod, root) {
    // Custom install for BepInEx.
    // ! Do not add BepInEx install to currentItems.
    let readDir = fs.readdirSync(path.join(mod.Path, "BepInExPack/"));

    for (let i=0; i<readDir.length; i++) {
        fs.copySync(path.join(mod.Path, "BepInExPack/", readDir[i]), path.join(root, readDir[i]))
    }
    return {};
}

function installMod(mod, root) {
    if (mod.Name === "BepInExPack") {
        // Custom BepInEx install
        return installBepInEx(mod, root)
    } else {
        // Normal mod
        return installCustom(mod, root);
    }
}

function installCustom(mod, root) {
    // List possible folder names, and initialise dictionary
    let folderNames = ["plugins", "monomod", "patchers", "core", "config"];
    let traverseDirectory = (nestedPath) => {
        let toInstall = {}
        for (let index in folderNames) {
            toInstall[folderNames[index]] = [];
        }

        // Traverse directory
        let nestDir = fs.readdirSync(nestedPath);
        for (let i=0; i<nestDir.length; i++) {
            let absolute = path.join(nestedPath, nestDir[i]);
            // If dir
            if (fs.lstatSync(absolute).isDirectory()) {
                let lower = path.basename(absolute).toLowerCase();
                let found = false;
                // Search for folder matching an entry in folderNames
                for (let j=0; j<folderNames.length; j++) {
                    // If found, add absolute path to toInstall[folderName]
                    if (folderNames[j] === lower) {
                        found = true;
                        let toCopy = fs.readdirSync(path.join(nestedPath, nestDir[i]));
                        for (let j=0; j<toCopy.length; j++) {
                            toInstall[lower].push(path.join(absolute, toCopy[j]));
                        }
                    }
                }
                // If not found, continue traversing
                if (!found) {
                    let res = traverseDirectory(path.join(nestedPath, path.basename(absolute)));
                    for (let key in res) {
                        for (let nestLoop=0; nestLoop<res[key].length; nestLoop++) {
                            toInstall[key].push(res[key][nestLoop]);
                        }
                    }
                }
            } else {
                // File is not a folder, therefore find a suitable location.
                if (path.extname(path.basename(absolute, ".dll")) === ".mm") {
                    toInstall["monomod"].push(absolute);
                } else {
                    toInstall["plugins"].push(absolute);
                }
            }
        }
        return toInstall;
    }
    return traverseDirectory(mod.Path);
}