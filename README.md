# wat (Automation Toolbox)

A clean, modern, high-performance desktop toolbox built with **Wails (Go + Svelte 5 + TypeScript + Vite)**.

## Live Development

To run the application in live development mode:
```bash
wails dev
```

## Production Building

### Modern Linux (Arch Linux, Fedora, Ubuntu 23.04+)
Modern distributions use **WebKit2GTK 4.1** (as `webkit2gtk-4.0` has been deprecated and removed). To compile for these platforms, build with the `webkit2_41` tag:
```bash
wails build -tags webkit2_41
```
* **Output Binary**: `build/bin/wat`

### Windows Cross-Compilation
To cross-compile a standalone single-file Windows executable (`wat.exe`) from your Linux environment, ensure you have the `mingw-w64-gcc` toolchain installed, and run:
```bash
wails build -platform windows/amd64
```
* **Output Binary**: `build/bin/wat.exe`

### Binary Compression (UPX)
To compress the compiled binaries (often reducing their size by 50% to 70%), install the `upx` utility:
```bash
# On Arch Linux
ipkgs upx
```
Once installed, append the `-upx` flag to your build commands:
```bash
# Compressed Linux Build
wails build -tags webkit2_41 -upx

# Compressed Windows Build
wails build -platform windows/amd64 -upx
```


