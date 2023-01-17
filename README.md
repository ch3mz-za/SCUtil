# SCUtil

Utility that does some useful things with your Star Citizen directory

## 1. Features

### 1.1 Clear all data except p4k

> Clears all data in the `Star Citizen` folder except the `Data.p4k` file. This is useful
> to sometimes clear out odd issues that pop up.

### 1.2 Clear `USER` folder (excluding control mappings)

> Clears the `USER` folder in the `Star Citizen` folder, excluding control mappings. This
> is useful to clear issues relating to old user files.

### 1.3 Clear `USER` folder (including control mappings)

> Clears the `USER` folder in the `Star Citizen` folder, including control mappings. This
> is useful to clear issues relating to old user files.

### 1.4 Reads all the filenames in the p4k data file

> This read all the filenames included wihtin the `Data.p4k` file and
> writes it out to a file (`P4k_filenames.txt`). This is for the curious
> individuals.

### 1.5 Search p4k filenames

> This features takes a phrase and searches for filenames within the
> Data.p4k which contain the phrase.

### 1.6 Clear Star Citizen App Data (Windows AppData)

> This clears out error logs that are typically found within Star
> Citizen's App Data which sometimes prevent the game from starting.

### 1.7 Clear RSI Launcher data (Windows AppData)

> This clears out logs and cached items that are typically found within > the RSI Launcher's App Data which sometimes prevent the game from
> starting.

## 2. Compiling & Run SCUtil

On windows, with Golang, simply compile the code using :

```bash
go build -o bin/SCUtil.exe main.go
```

Place `SCUtil.exe` within a folder located in one of the parent directories of your `Star Citizen` folder:

```txt
Your_Game_Dir
│
└───SCUtil
│   │   SCUtil.exe
│   
└───Star Citizen
│   └───LIVE
│   └───PTU
│
└───RSI Launcher
```
