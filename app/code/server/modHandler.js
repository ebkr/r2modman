const fs = require("fs");
const path = require('path');
const fetch = require("node-fetch");
const extract = require('extract-zip')

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

    Update(modListString) {
        fs.writeFileSync(path.join(this.profile, "mods.json"), modListString, (e) => {
            console.log(e);
        });
    }

    Extract(zipLocation, location, callback) {
        extract(zipLocation, {dir: location}, (e) => {
            if (e) {
                callback(false);
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

exports.ModHandler = ModHandler;