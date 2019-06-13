# r2modman is in the progress of a framework overhaul due to stability purposes. It will be transitioning from GTK3 to Qt. All current updates are paused whilst this is in progress.

![Logo](https://i.imgur.com/rdImc3h.png)

## r2modman : Risk of Rain 2 Mod Manager

A simple, elegant, and easy-to-use mod manager for Risk of Rain 2.

---

### Current Features:
- Thunderstore Integration
- Local Mods
- Mod Updates
- Enable or disable mods

Mods can be downloaded, and updated, using the Thunderstore integration directly within the application.

---

### Screenshots:

![MainScreen](https://i.imgur.com/gpk8zNk.png)

![OnlineScreen](https://i.imgur.com/PQFfCwA.png)

---

### Notes:
- BepInEx is not currently installable through the application.

---

### Credits

r2modman is written in [Go](https://golang.org).

The interface uses [GTK+3](https://gtk.org), using the [gotk3](https://github.com/gotk3/gotk3) binding.

Moving mods to the correct directories was aided using the [copy](https://github.com/otiai10/copy) package provided by otiai10.

Unzipping uses the [unzip](https://github.com/xyproto/unzip) package provided by xyproto.
