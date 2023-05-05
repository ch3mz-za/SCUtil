# SCUtil

Utility that does some useful things with your Star Citizen directory

## 1. Features

### 1.1. Clear all data except p4k

> Clears all data in the `Star Citizen` folder except the `Data.p4k` file. This is useful
> to sometimes clear out odd issues that pop up.

### 1.2. Clear `USER` folder (excluding control mappings)

> Clears the `USER` folder in the `Star Citizen` folder, excluding control mappings. This
> is useful to clear issues relating to old user files.

### 1.3. Clear `USER` folder (including control mappings)

> Clears the `USER` folder in the `Star Citizen` folder, including control mappings. This
> is useful to clear issues relating to old user files.

### 1.4. Reads all the filenames in the p4k data file

> This read all the filenames included wihtin the `Data.p4k` file and
> writes it out to a file (`P4k_filenames.txt`). This is for the curious
> individuals.

### 1.5. Search p4k filenames

> This features takes a phrase and searches for filenames within the
> Data.p4k which contain the phrase.

### 1.6. Clear Star Citizen App Data (Windows AppData)

> This clears out error logs that are typically found within Star
> Citizen's App Data which sometimes prevent the game from starting.

### 1.7. Clear RSI Launcher data (Windows AppData)

> This clears out logs and cached items that are typically found within > the RSI Launcher's App Data which sometimes prevent the game from
> starting.

### 1.8. Backup and restore control mapping

> This does a backup of your game's custom control mappings and stores it within the `BACKUPS` folder as well as give you the option to restore previously backed up control mappings.
> The backup location can be found in the parent directory of your Star Citizen game files.

### 1.9. Backup screenshots

> This gives the user the ability to backup their Star Citizen screenshots within the `BACKUPS` folder.

### 1.10. Backup USER directory

> This will store a copy of the game's USER directory within the `BACKUPS` folder.

## 2. Running SCUtil

### 2.1. Executable

Simply download the release and place `SCUtil.exe` within any folder of your choosing. Set your game directory in setting. A `config.toml` file should appear within the root directory of the `SCUtil.exe`. Below is a breakdown of the app's directory.

```txt
SomeDirectory
│   SCUtil.exe
│   config.toml
└───BACKUPS
    └───Screenshots
    │   └───LIVE
    │   └───PTU
    └───ControlMappings
        └───LIVE
        └───PTU

```

### 2.2. Compile & Run SCUtil

On windows, with latest Golang version, simply compile the code using:

```bash
go build -o bin/SCUtil.exe main.go
```

With the executable, follow the instructions in section `2.1`.

> **NOTE** <br>
> Since version 2.0.0 of SCUtil, compilation will require the use of [Fyne](https://github.com/fyne-io/fyne) and a GCC compiler.
>
