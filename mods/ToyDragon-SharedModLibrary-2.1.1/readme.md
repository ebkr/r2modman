# Frogtown Mod Manager
Contains UI for enabling, disabling, and updating mods.

![In game popup](https://github.com/ToyDragon/ROR2ModShared/blob/master/Images/ingame.png?raw=true)

## Usage
Toggle mods on and off using the checkbox in the far right column. Most mods don't support being toggled on or off without restarting, so you may see a red "R" in the status column indicating you need to restart for the change to take effect. If the mod has an associated thunderstore page it will automatically be checked for updates, and if one is available you can click the new version text to jump to the page and download it. Close the popup with escape and open it with ctrl+F10.

![Close up](https://github.com/ToyDragon/ROR2ModShared/blob/master/Images/closeup.png?raw=true)

 Use the checkbox in the left column to collapse mods from the same author, so that you can remove clutter and enable/disable all of them at once. Mouse over a mod to remind yourself what it does.

![description](https://github.com/ToyDragon/ROR2ModShared/blob/master/Images/tooltip.png?raw=true)

This mod is a prerequisite for:
- [Frogtown Cheats](https://thunderstore.io/package/ToyDragon/CheatingChatCommands/)
- [Healing Helper Mod](https://thunderstore.io/package/ToyDragon/HealingHelpers/)
- [Character Randomizer](https://thunderstore.io/package/ToyDragon/CharacterRandomizer/)
- [Engineer Fixes](https://thunderstore.io/package/ToyDragon/EngineerLunarCoinsFix/)

## Installation
1. Install [BepInEx Mod Pack](https://thunderstore.io/package/bbepis/BepInExPack/)
2. Download the latest ToyDragon-SharedModLibrary.zip
3. Unzip it and move FrogtownShared.dll to your \BepInEx\plugins folder

## Frogtown Mod Manager Versions
- 2.1.1
  - Better error handling when mods have invalid or duplicate GUIDs.
- 2.1.0
  - Check thunderstore for mod updates instead of github.
  - Better prevention of cheats being after game starts.

## Developers
This library can help you:
- Toggle your mods on/off
- Maintain the isModded flag
- Distribute updates

```C#
public ModDetails modDetails;
public void Awake()
{
    //Initializing a ModDetails object will allow you to assign
    //a short description and github repository to your mod, and allow
    //it to be enabled or disabled without needing to restart the game.
    //Otherwise when disabled the DLL containing your plugin will be 
    //moved to a DisabledMods folder that BepInEx doesn't scan. If your
    //mod relies on any other external files being in the same place
    //this may cause issues.
    modDetails = new ModDetails("com.frogtown.chatcheats")
    {
        //This description shows up as a tooltip when hovering over
        //your mod in the mod list, it should be very short. A link to
        //your thunderstore page will be included, feel free to put any
        //additional documentation there.
        description = "Adds the /change_char and /give_item chat commands.",
        
        //author is a string used to group your mods together in the
        //mod list.
        author = "ToyDragon",

        //thunderstoreFullName is used to search for your mod in the
        //thunderstore API. It will be your dependency string without
        //the version suffix.
        thunderstoreFullName = "ToyDragon-CheatingChatCommands",
    };
    FrogtownShared.RegisterMod(modDetails);
    
    //When all mods are disabled the isModded flag will be updated
    //to false, and when any are enabled it will set it back to true.
}
```

## Releases
The manager will check the thunderstore listing at most twice a day, and display bright blue text if there is a newer version the user doesn't have installed. The user can click the link to be brought to your thunderstore page, so include installation instructions and version notes there.
