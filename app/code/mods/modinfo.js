
/*
* ModInfo.js
* -------------
* Used to manage mod data and their conversions
*/

// Manage mods
class ModManager {
    installedMods = [] // Mod_V1
    availableMods = function() {
        let modArray = [];
        let modList = JSON.parse(localStorage.getItem("modList"));
        for (let i=0; i<modList.length; i++) {
            let mod = new ThunderstoreMod();
            mod.FromJSON(modList[i]);
            modArray.push(mod);
        }
        return modArray;
    }()

    FilesToInstalled(fileDataList) {
        this.installedMods = [];
        let json = JSON.parse(fileDataList);
        for (let i=0; i<json.length; i++) {
            if (json[i].Manifest == 1) {
                let mod = new Mod_V1();
                mod.FromFile(json[i]);
                this.installedMods.push(mod);
            }
        }
    }



    installed = {
        GetModsByFilter: function(str, self) {
            let myArr = [];
            for (let i=0; i<self.installedMods.length; i++) {
                if (self.installedMods[i].Name.toLowerCase().search(str) >= 0) {   
                    myArr.push(self.installedMods[i]);
                }
            }
            return myArr;
        },
        HasUpdate: function(mod, self) {
            if (!mod.IsHTTP) {
                console.log("Not HTTP");
                return false;
            }
            let comparison = null;
            for (let i=0; i<self.availableMods.length; i++) {
                if (self.availableMods[i].uuid4 === mod.Uuid4) {
                    comparison = self.availableMods[i].versions[0];
                    break;
                }
            }
            if (comparison !== null) {
                let oldMod = mod;
                mod = mod.Version;
                let oldComparison = comparison;
                comparison = new Version();
                comparison.ConvertFromString(oldComparison.version_number);
                if (mod.Major < comparison.Major) {
                    return true;
                } else if (mod.Major == comparison.Major && mod.Minor < comparison.Minor) {
                    return true;
                } else if (mod.Major === comparison.Major && mod.Minor === comparison.Minor && mod.Patch < comparison.Patch) {
                    return true;
                }
            }
            return false;
        },
        HasMissingDependency: function(mod, self) {
            let listToScan = self.installed.GetModsByFilter("", self);
            for (let depIter in mod.Dependencies) {
                let found = false;
                let dep = mod.Dependencies[depIter];
                for (let listIter in listToScan) {
                    let installed = listToScan[listIter];
                    if (installed.FullName === dep.Name) {
                        let ver = new Version();
                        let comVer = new Version();
                        ver.FromDict(installed.Version);
                        comVer.FromDict(dep.Version);
                        if (ver.IsNewerOrEqual(comVer)) {
                            console.log("Version Newer/Equal")
                            found = true;
                            break;
                        }
                    }
                }
                if (!found) {
                    console.log("Not found");
                    return dep.Name;
                }
            }
            return false;
        }
    }
    available = {
        GetModsByFilter: function(str, self) {
            let myArr = [];
            for (let i=0; i<self.availableMods.length; i++) {
                if (self.availableMods[i].name.toLowerCase().search(str) >= 0) {
                    myArr.push(self.availableMods[i]);
                }
            }
            return myArr;
        },
        GetLatestVersion: function(uuid4, self) {
            for (let i=0; i<self.availableMods.length; i++) {
                if (self.availableMods[i].uuid4 === uuid4) {
                    return self.availableMods[i].versions[0];
                }
            }
            return false
        },
        GetRootFromUUID4(uuid4, self) {
            for (let i=0; i<self.availableMods.length; i++) {
                if (self.availableMods[i].uuid4 === uuid4) {
                    return self.availableMods[i];
                }
            }
            return false
        }
    }
}

// Record version information for a mod
class Version {

    Major = 0   // Int
    Minor = 0   // Int
    Patch = 0   // Int

    // Convert the Version to a VersionString
    ConvertToString() {
        return `${this.Major}.${this.Minor}.${this.Patch}`;
    }

    // Convert and use a VersionString
    ConvertFromString(str) {
        let conv = str.split(".")
        this.Major = Number(conv[0]);
        this.Minor = Number(conv[1]);
        this.Patch = Number(conv[2]);
    }

    // Return true if newer or equal version
    IsNewerOrEqual(comp) {
        if (this.Major === comp.Major && this.Minor === comp.Minor && this.Patch === comp.Patch) {
            return true;
        }
        return this.Major > comp.Major ? true : 
                this.Major === comp.Major && this.Minor > comp.Minor ? true : 
                this.Major === comp.Major && this.Minor === comp.Minor && this.Patch > comp.Patch ? true : false;
    }

    FromDict(dict) {
        for (let key in dict) {
            this[key] = dict[key];
        }
    }
}

// Used to reference dependencies easier
class Dependency {
    Name    = "" // String
    Version = new Version() // Version

    FromString(str) {
        let split = str.split("-");
        for (let i=0; i<split.length-1; i++) {
            this.Name += split[i];
            if (i<split.length-2) {
                this.Name += "-";
            }
        }

        this.Version.ConvertFromString(split[split.length-1]);
    }
}

// Used to record the type of manifest used
class ManifestVersion {
    Manifest = 1    // Int
}

// The primary mod object
class Mod_V1 extends ManifestVersion {
    Name            = ""            // String
    URL             = ""            // String
    Path            = ""            // String
    Icon            = ""            // String (URL)
    Uuid4           = ""            // String
    Author          = ""            // String
    FullName        = ""            // String
    Description     = ""            // String
    Dependencies    = []            // Dependency
    Enabled         = true          // Boolean
    IsHTTP          = false         // Boolean
    Version         = new Version() // Version

    // Convert JSON to Mod
    //Returns false if Mod is V2
    FromManifest(raw) {
        let json = JSON.parse(raw.trim(), true);
        for (let prop in json) {
            // Manifest V1
            let isManifestV2 = false;
            switch(prop) {
                case "name":
                    this.Name = json[prop];
                    this.FullName = json[prop];
                    break;
                case "version_number":
                    this.Version.ConvertFromString(json[prop]);
                    break;
                case "website_url":
                    this.URL = json[prop];
                    break;
                case "description":
                    this.Description = json[prop];
                    break;
                case "dependencies":
                    for (let i=0; i<json[prop].length; i++) {
                        let dep = new Dependency();
                        dep.FromString(json[prop][i])
                        this.Dependencies.push(dep);
                    }
                    break
                default:
                    // TODO: Add support for Manifest V2
                    this.Manifest = 2;
                    return false;
            }
            // Manifest is V1.
            this.Manifest = 1;
        }
        return true
    }

    FromFile(json) {
        for (let prop in json) {
            this[prop] = json[prop]
        }
    }
}

// ! : Used in the future to support V1 <-> V2 compatibility
class ModVersionBridge {
    FullName
    Name
}

// ------------------------
// * Thunderstore Related *
// ------------------------

// Used to record data from Thunderstore
class ThunderstoreMod {
    name            = ""    // String
    owner           = ""    // String
    uuid4           = ""    // String
    full_name       = ""    // String
    package_url     = ""    // String
    date_created    = ""    // String
    date_updated    = ""    // String
    versions        = []    // ThunderstoreVersion
    is_pinned       = false // Boolean
    is_deprecated   = false // Boolean

    FromJSON(json) {
        for (let prop in json) {
            switch(prop) {
                case "versions":
                    for (let i=0; i<json[prop].length; i++) {
                        let tsv = new ThunderstoreVersion();
                        tsv.FromJSON(json[prop][i]);
                        this.versions.push(tsv);
                    }
                    break;
                default:
                    this[prop] = json[prop];
                    break;
            }
        }
    }
}

// Information about mod versions via Thunderstore
class ThunderstoreVersion {
    icon            = ""    // String
    name            = ""    // String
    uuid4           = ""    // String
    full_name       = ""    // String
    description     = ""    // String
    website_url     = ""    // String
    download_url    = ""    // String
    date_created    = ""    // String
    version_number  = ""    // String
    downloads       = 0     // Int
    dependencies    = []    // Dependency
    is_active       = true  // Boolean

    FromJSON(json) {
        for (let prop in json) {
            switch(prop) {
                case "dependencies":
                    for (let i=0; i<json[prop].length; i++) {
                        let dep = new Dependency();
                        dep.FromString(json[prop][i]);
                        this.dependencies.push(dep);
                    }
                    break;
                default:
                    this[prop] = json[prop];
                    break;
            }
        }
    }
}

// Will error on web-view, but works on server.
exports.Mod_V1 = Mod_V1
exports.Version = Version