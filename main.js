// Modules to control application life and create native browser window
const {app,BrowserWindow,dialog,shell,protocol} = require('electron');
const path = require('path');
const {ipcMain} = require('electron');
const ipcServer = require('node-ipc');
const fs = require("fs-extra");
const prompt = require('electron-prompt');
const modHandler = require("./app/code/server/modHandler")
const modinfo = require("./app/code/mods/modinfo")
const conf = require("./app/code/other/configclass")
const { execSync } = require('child_process');
var regedit = require('regedit')

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
var mainWindow;
var stage = 0;
var selectedProfile = "";
var dir = path.join(process.execPath, "../", "mods", "profiles");
var mHandler = new modHandler.ModHandler();
var tsData = null;
let downloadsInProgress = 0;

// Download using the ror2mm:// protocol
function downloadModFromProtocol(protocol) {
	let strSplit = protocol.replace("//", "/").split("/");
	let location = strSplit[3];
	let author = strSplit[4];
	let mod = strSplit[5];
	let version = strSplit[6];
	if (location && author && mod && version) {
		for (let iter in tsData) {
			let tsMod = tsData[iter];
			if (mod.toLowerCase() === tsMod.name.toLowerCase() && author.toLowerCase() === tsMod.owner.toLowerCase()) {
				// Should be the correct mod
				console.log(tsMod.full_name, true);
				for (let verIter in tsMod.versions) {
					let ver = tsMod.versions[verIter];
					if (ver.version_number === version) {
						let interval;
						interval = setInterval(()=>{
							if (selectedProfile !== "") {
								downloadMod(null, ver, tsMod);
								clearInterval(interval);
							}
						}, 100);
						break;
					}
				}
			}
		}
	} else {
		alert("Invalid install link");
	}
} 

function createWindow() {
	// Create the browser window.
	mainWindow = new BrowserWindow({
		frame: false,
		height: 209,
		width: 600,
		webPreferences: {
			preload: path.join(__dirname, 'preload.js')
		}
	})

	// and load the index.html of the app.
	mainWindow.loadFile('app/splash.html');

	// Open the DevTools.
	// mainWindow.webContents.openDevTools()

	// Emitted when the window is closed.
	mainWindow.on('closed', function () {
		// Dereference the window object, usually you would store windows
		// in an array if your app supports multi windows, this is the time
		// when you should delete the corresponding element.
		if (stage == 0) {
			mainWindow = null;
		}
	})
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', function() {
	// Only allow a single instance of the executable at a time
	let reqLockSuccess = app.requestSingleInstanceLock();
	if (!reqLockSuccess) {
		// If this isn't the single instance,
		// connect to r2mm IPC and send the install parameter.
		ipcServer.connectTo("r2mm", ()=>{
			if (process.argv.length >= 2) {
				ipcServer.of.r2mm.emit("install", process.argv[1]);
			}
			app.quit();
		});
	} else {
		ipcServer.config.id = "r2mm";
		ipcServer.serve(
			"/tmp/app.r2mm",
			() => {
				ipcServer.server.on("install", (res)=>{
					downloadModFromProtocol(res);
				})
			}
		);
		ipcServer.server.start();
	}
	createWindow();
});

// Quit when all windows are closed.
app.on('window-all-closed', function () {
	// On macOS it is common for applications and their menu bar
	// to stay active until the user quits explicitly with Cmd + Q
	if (process.platform !== 'darwin') {
		app.releaseSingleInstanceLock();
		ipcServer.server.stop();
		app.quit();
	}
})

app.on('activate', function () {
	// On macOS it's common to re-create a window in the app when the
	// dock icon is clicked and there are no other windows open.
	if (mainWindow === null) createWindow();
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
ipcMain.on("splashFinished", (e, fetched) => {
	// Profile Window
	tsData = JSON.parse(fetched);
	let newWindow = CreateWindow({
		height: 440,
		width: 410,
		frame: true,
		webPreferences: {
			preload: path.join(__dirname, 'preload.js')
		}
	}, 1)
	mainWindow.loadFile('app/profiles.html');

	// Bind getProfiles to window call
	ipcMain.on('getProfiles', (e) => {
		fs.mkdirsSync(dir);
		mainWindow.webContents.send("getProfiles", fs.readdirSync(dir).filter((file) => {
			return fs.statSync(dir + '/' + file).isDirectory();
		}));
	});

	// Bind promptProfileCreation to window call
	ipcMain.on("promptProfileCreation", (e) => {
		prompt({
			height: 160,
			title: 'Enter a new profile name',
			label: 'Profile Name:',
			value: "",
			inputAttrs: {
				type: 'string'
			}
		}, mainWindow).then((r) => {
			if (!fs.existsSync(path.join(dir, r))) {
				fs.mkdirsSync(path.join(dir, r), {
					recursive: true
				});
				mainWindow.webContents.send("getProfiles", fs.readdirSync(dir).filter((file) => {
					return fs.statSync(dir + '/' + file).isDirectory();
				}));
			}
		}).catch((e) => {
			console.log(e);
		})
	})

	ipcMain.on("deleteProfile", (e, profile) => {
		fs.emptyDirSync(path.join(dir, profile));
		fs.rmdirSync(path.join(dir, profile));
		mainWindow.webContents.send("getProfiles", fs.readdirSync(dir).filter((file) => {
			return fs.statSync(dir + '/' + file).isDirectory();
		}));
	});

	ipcMain.on("selectedProfile", (e, profile) => {
		selectedProfile = profile;
		startManagerWindow();
		if (process.argv.length >= 2) {
			downloadModFromProtocol(process.argv[1]);
		}
	});
});

function startManagerWindow() {
	let newWindow = CreateWindow({
		height: 350,
		width: 800,
		frame: true,
		webPreferences: {
			preload: path.join(__dirname, 'preload.js')
		}
	}, 2)
	mHandler.profile = path.join(dir, selectedProfile);
	ipcMain.on("getInstalledMods", (e, profile) => {
		mainWindow.webContents.send("getInstalledMods", mHandler.GetMods(mainWindow));
	});
	mainWindow.loadFile('app/manager_2.html');
}

ipcMain.on('resize', function (e, x, y) {
	mainWindow.setSize(x, y);
});

function downloadMod(e, version, tsmod) {
	downloadsInProgress += 1;
	console.log(version);
	let progressWindow = new BrowserWindow({
		width: 300,
		height: 130,
		frame: true,
		webPreferences: {
			preload: path.join(__dirname, 'preload.js'),
			devTools: false,
		}
	});
	progressWindow.setMenu(null);
	progressWindow.setAlwaysOnTop(true, "floating");
	progressWindow.loadURL(path.join(__dirname, "app", "progress.html?modName=" + tsmod.full_name));
	mHandler.Download(tsmod, version, (res, message) => {
		progressWindow.close()
		if (!res) {
			downloadsInProgress -= 1;
			mainWindow.webContents.send("downloadError", message);
		} else {
			// Create mod
			let mod = generateModFromManifest(message.extractLocation)
			mod.Uuid4 = tsmod.uuid4;
			mod.FullName = tsmod.full_name;
			mod.IsHTTP = true;
			mod.Path = message.extractLocation;
			mod.URL = tsmod.package_url;
			mod.Description = version.description;
			mod.Version.ConvertFromString(version.version_number);
			mod.Author = tsmod.owner;
			mod.Icon = path.join(message.extractLocation, "icon.png");
			let modArray = [];
			modArray = JSON.parse(mHandler.GetMods(), true)
			insertOrReplaceMod(modArray, mod)
			let listString = JSON.stringify(modArray, null, 4)
			mHandler.Update(listString);
			mainWindow.webContents.send("getInstalledMods", listString);
			downloadsInProgress -= 1;
		}
	});
};
ipcMain.on('downloadMod', function (e, version, tsmod) {
	downloadMod(e, version, tsmod);
})

function CreateWindow(opts, closeStage) {
	//opts.devTools = false
	let window = new BrowserWindow(opts);
	//window.setMenu(null);
	opts.icon = path.join(__dirname, "app", "assets", "r2.ico")
	window.on("close", function () {
		if (stage == closeStage) {
			mainWindow = null;
		}
	})
	stage += 1;
	old = mainWindow;
	mainWindow = window;
	old.close();
}

ipcMain.on("addLocalMod", function (e, x, y) {
	dialog.showOpenDialog(mainWindow, {
		title: "Select Mod to Install",
		filters: [{name: "zip", extensions: ["zip"]}],
		properties: ["openFile", "multiSelections"]
	}, (filePaths)=> {
		if (filePaths == null) {
			return;
		}
		let modArray = [];
		modArray = JSON.parse(mHandler.GetMods(), true)
		for (let i=0; i<filePaths.length; i++) {
			mHandler.Extract(filePaths[i], path.join(dir, selectedProfile, path.basename(filePaths[i], ".zip")), (valid)=>{
				if (valid) {
					let mod = generateModFromManifest(path.join(dir, selectedProfile, path.basename(filePaths[i], ".zip")))
					mod.Path = path.join(dir, selectedProfile, path.basename(filePaths[i], ".zip"))
					mod.Icon = path.join(mod.Path, "icon.png");
					insertOrReplaceMod(modArray, mod)
					let listString = JSON.stringify(modArray, null, 4)
					mHandler.Update(listString);
					mainWindow.webContents.send("getInstalledMods", listString);
				}
			})
		}
	});
});

function generateModFromManifest(directory) {
	let mod = new modinfo.Mod_V1()
	let buf = fs.readFileSync(path.join(directory, "manifest.json"), "utf8");
	buf.trim();
	if (!mod.FromManifest(buf)) {
		// is V2
		// TODO
	}
	return mod;
}

function insertOrReplaceMod(modArray, mod) {
	let foundInArray = false;
	modArray.forEach((value, key) => {
		if (value.FullName === mod.FullName) {
			foundInArray = true;
			modArray[key] = mod;
		}
	})
	if (!foundInArray) {
		if (mod.Name == "BepInExPack") {
			modArray.unshift(mod);
		} else {
			modArray.push(mod);
		}
	}
	return modArray;
}

ipcMain.on("playRoR2", function (e, x, y) {
	let configFile = path.join(dir, "config.json");
	let config = new conf.Configuration(configFile);
	if (!config.GamePathValid()) {
		dialog.showOpenDialog(mainWindow, {
			title: "Select Risk of Rain 2 Install Directory",
			properties: ["openDirectory"]
		}, (filePaths)=> {
			if (filePaths !== null && filePaths.length == 1) {
				let path = filePaths[0];
				config.gamePath = path;
				config.Update(configFile);
			}
		});
	}
	if (config.GamePathValid()) {
		// Install Mods
		config.InstallMods(JSON.parse(mHandler.GetMods()));
		config.Update(configFile);
		execSync('start steam://run/632360');
	}
});

ipcMain.on("configureMod", function(e, mod) {
	tempWindow = new BrowserWindow({
		width: 500,
		height: 350,
		frame: true,
		webPreferences: {
			preload: path.join(__dirname, 'preload.js'),
			devTools: true,
		}
	});
	tempWindow.on("close", function () {
		mainWindow.webContents.send("getInstalledMods", mHandler.GetMods());
		mainWindow.webContents.send("configureClose");
		ipcMain.removeListener("toggleModUsability", toggleModUsability);
		ipcMain.removeListener("removeMod", removeMod);
	})
	let toggleModUsability = ()=>{
		mod.Enabled = !mod.Enabled;
		let modList = JSON.parse(mHandler.GetMods().trim())
		for (let modIndex in modList) {
			if (modList[modIndex].FullName === mod.FullName) {
				modList[modIndex].Enabled = mod.Enabled;
			}
		}
		mHandler.Update(JSON.stringify(modList, null, 4));
		mainWindow.webContents.send("getInstalledMods", JSON.stringify(modList));
	}
	let removeMod = ()=>{
		let modList = JSON.parse(mHandler.GetMods().trim())
		for (let modIndex in modList) {
			if (modList[modIndex].FullName === mod.FullName) {
				modList.splice(modIndex, 1);
				break;
			}
		}
		fs.removeSync(mod.Path);
		mHandler.Update(JSON.stringify(modList, null, 4));
		tempWindow.close();
	}

	ipcMain.on("toggleModUsability", toggleModUsability);
	ipcMain.on("removeMod", removeMod);
	tempWindow.loadURL(path.join(__dirname, "app", "configure.html"));
})

ipcMain.on("goToSite", (e, site)=>{
	execSync('start ' + site);
});

ipcMain.on("exportProfile", (e, site)=>{
	let parsed = JSON.parse(mHandler.GetMods());
	let modArr = [];
	for (let modi in parsed) {
		let mod = parsed[modi];
		if (mod.IsHTTP) {
			modArr.push({
				uuid4: mod.Uuid4,
				version: mod.Version
			});
		}
	}
	fs.writeFileSync(path.join(mHandler.profile, "export.json"), JSON.stringify(modArr));
	shell.showItemInFolder(path.join(mHandler.profile, "export.json"));
});

ipcMain.on("importProfile", (e)=>{
	dialog.showOpenDialog(mainWindow, {
		title: "Select profile to import",
		filters: [{name: "json", extensions: ["json"]}],
		properties: ["openFile"]
	}, (filePaths)=> {
		if (filePaths == null) {
			return;
		}
		if (fs.existsSync(filePaths[0])) {
			let data = fs.readFileSync(filePaths[0], "utf8");
			let json = JSON.parse(data);
			let installList = [];
			for (let iter in json) {
				let imported = json[iter];
				for (let ts in tsData) {
					let tsmod = tsData[ts];
					if (tsmod.uuid4 === imported.uuid4) {
						let vsNum = new modinfo.Version();
						vsNum.FromDict(imported.version);
						let strVer = vsNum.ConvertToString();
						for (let vIter in tsmod.versions) {
							let version = tsmod.versions[vIter];
							if (version.version_number === strVer) {
								installList.push([version, tsmod]);
							}
						}
					}
				}
			}
			let interval;
			interval = setInterval(()=>{
				if (installList.length == 0) {
					console.log("Disconnect Importer");
					clearInterval(interval);
				} else {
					if (downloadsInProgress == 0) {
						// InstallList [{version} {tsmod}]...;
						downloadMod(null, installList[0][0], installList[0][1]);
						installList.shift();
					}
				}
			}, 100);
		}
	});
});

ipcMain.on("associateHandler", (e) => {
	app.setAsDefaultProtocolClient("ror2mm", process.execPath);
	mainWindow.webContents.send("isAssociatedSuccess", true);
});

ipcMain.on("isAppAssociated", (e)=>{
	mainWindow.webContents.send("isAppAssociated", app.isDefaultProtocolClient('ror2mm'));
})

ipcMain.on("disableAll", (e) => {
	let mods = JSON.parse(mHandler.GetMods());
	for (let modIter in mods) {
		let mod = mods[modIter];
		mod.Enabled = false;
	}
	mHandler.Update(JSON.stringify(mods, null, 4));
	mainWindow.webContents.send("refresh");
});

ipcMain.on("enableAll", (e) => {
	let mods = JSON.parse(mHandler.GetMods());
	for (let modIter in mods) {
		let mod = mods[modIter];
		mod.Enabled = true;
	}
	mHandler.Update(JSON.stringify(mods, null, 4));
	mainWindow.webContents.send("refresh");
});

ipcMain.on("deleteAll", (e) => {
	let modList = JSON.parse(mHandler.GetMods().trim())
	for (let modIndex in modList) {
		fs.removeSync(modList[modIndex].Path);
	}
	modList = [];
	mHandler.Update(JSON.stringify(modList, null, 4));
	mainWindow.webContents.send("refresh");
});