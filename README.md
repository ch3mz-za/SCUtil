# SCUtil
Utility that does some useful things with your Star Citizen directory

## 1. Features

### 1.1 Clear all data except p4k

> Clears all data in the `Star Citizen` folder except the `Data.p4k` file. This is useful 
> to sometimes clear out odd issues that pop up.

### 1.2 Clear `USER` folder

> Clears the `USER` folder in the `Star Citizen` folder This is useful to sometimes clear 
> out odd issues that pop up.

## 2. Compiling & Run SCUtil

On windows, with Golang, simply compile the code using :

```bash
go build -o bin/SCUtil.exe main.go
```

Place `SCUtil.exe` within a folder located in one of the parent directories of your `Star Citizen` folder:

```
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
