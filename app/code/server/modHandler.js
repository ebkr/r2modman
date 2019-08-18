const fs = require("fs-extra");
const path = require('path');
const fetch = require("node-fetch");
const extract = require('extract-zip');
const unzipper = require("unzipper");
const rsync = require("rsyncwrapper");

class ModHandler {
    profile = ""

    Download(mod, version, callback) {
        fetch(version.download_url, {
            method: 'get'
        }).then((response) => {
            return response.arrayBuffer();
        }).then((response) => {
            // Write to file
            fs.writeFileSync(path.join(this.profile, mod.full_name + ".zip"), Buffer.from(response), (e) => {
                callback(false, e)
                return false
            })
            this.Extract(path.join(this.profile, mod.full_name + ".zip"), path.join(this.profile, mod.full_name), (result) => {
                if (!result) {
                    callback(false, e)
                    return false;
                } else {
                    fs.unlinkSync(path.join(this.profile, mod.full_name + ".zip"), (e) => {
                        callback(false, e)
                        return false;
                    })
                    callback(true, {
                        extractLocation: path.join(this.profile, mod.full_name),
                    })
                    return true;
                }
            });
        }).catch((e) => {
            callback(false, e)
        })
    }

    UpdateR2MM(url, callback) {
        let name = "app-update";
        fetch(url, {
            method: 'get'
        }).then((response) => {
            console.log(response);
            return response.arrayBuffer();
        }).then((response) => {
            let proceed = true;
            fs.mkdirSync(path.join(process.cwd(), "tmp"), (err)=>{
                console.log("Error creating temporary directory:", err);
            });
            if (proceed) {
                fs.writeFileSync(path.join(process.cwd(), "tmp", name + ".zip"), Buffer.from(response), (e) => {
                    proceed = false;
                    return false;
                })
            }
            if (proceed) {
                fs.createReadStream(path.join(process.cwd(), "tmp", name + ".zip"))
                    .pipe(unzipper.Extract({ path: path.join(process.cwd(), "tmp", name) }))
                    .on("close", ()=>{
                        let locationPath = "";
                        let zipResources = "";
                        for (let i=0; i<process.argv.length; i++) {
                            if (process.argv[i].toLowerCase().split("ror2mm://").length > 1) {
                                locationPath = process.cwd();
                                zipResources = path.join(process.cwd(), "tmp", name, "resources", "app");
                            }
                        }
                        if (locationPath === "") {
                            locationPath = path.join(process.cwd());
                            zipResources = path.join(process.cwd(), "tmp", name);
                        }
                        console.log("LocationPath:", locationPath);
                        recursiveDirectoryOverwrite(zipResources, locationPath, ()=>{
                            console.log("Done recursive")
                            callback(true);
                        })
                    })
            }
            if (!proceed) {
                callback(false, new Error("Something failed during the update process"));
            }
        }).catch((err)=>{
            console.log("Oof:", err);
            callback(false, err);
        })
    }

    Update(modListString) {
        fs.writeFileSync(path.join(this.profile, "mods.json"), modListString, (e) => {
            console.log(e);
        });
    }

    Extract(zipLocation, location, callback) {
        extract(zipLocation, {dir: location}, (e) => {
            if (e) {
                callback(false, e);
            }
            callback(true);
        });
    }

    Play() {
        // TODO
    }

    GetMods() {
        let modFile = path.join(this.profile, "mods.json");
        if (fs.existsSync(modFile)) {
            // Mod file located:
            let buf = fs.readFileSync(modFile, "utf8");
            return buf
        } else {
            fs.writeFileSync(modFile, "[]", (e) => {
                console.log(e);
            })
            return "[]";
        }
    }

}

async function recursiveDirectoryRunnable(root, location) {
    console.log("rdo: root:", root);
    console.log("rdo: loc:", location);
    let contents = fs.readdirSync(root)
    for (let i=0; i<contents.length; i++) {
        if (fs.statSync(path.join(root, contents[i])).isDirectory()) {
            await recursiveDirectoryRunnable(path.join(root, contents[i]), path.join(location, contents[i]));
        } else {
            fs.removeSync(path.join(location, contents[i]))
            fs.moveSync(path.join(root, contents[i]), path.join(location, contents[i]));
        }
    }
}

async function recursiveDirectoryOverwrite(root, location, callback) {
    await recursiveDirectoryRunnable(root, location)
    callback(true);
}

exports.ModHandler = ModHandler;