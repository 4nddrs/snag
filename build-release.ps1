# Define output paths
$releasePath = "C:\Proyects\releases\v1.0.1"
$distPath = "$releasePath\dist"

# Create directory structure
if (!(Test-Path $releasePath)) { New-Item -ItemType Directory -Path $releasePath }
mkdir -p "$distPath\mac", "$distPath\amd64", "$distPath\arm64", "$distPath\windows"

Write-Host "🚀 Starting multi-platform build..." -ForegroundColor Cyan

# 1. Compile Binaries
Write-Host "📦 Building for macOS..."
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o "$distPath\mac\snag" .

Write-Host "📦 Building for Linux AMD64..."
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o "$distPath\amd64\snag" .

Write-Host "📦 Building for Linux ARM64..."
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o "$distPath\arm64\snag" .

Write-Host "📦 Building for Windows..."
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o "$releasePath\snag.exe" .

# 2. Create .tar.gz archives INSIDE each folder
Write-Host "🗜️  Compressing archives in their respective directories..."

# Mac
tar -czvf "$distPath\mac\snag-mac.tar.gz" -C "$distPath\mac" snag
# Linux AMD64
tar -czvf "$distPath\amd64\snag-linux-amd64.tar.gz" -C "$distPath\amd64" snag
# Linux ARM64
tar -czvf "$distPath\arm64\snag-linux-arm64.tar.gz" -C "$distPath\arm64" snag

# 3. Generate SHA256 Hashes for Homebrew Tap
Write-Host "`n🔑 SHA256 HASHES FOR YOUR FORMULA:" -ForegroundColor Yellow
$archFolders = "mac", "amd64", "arm64"
foreach ($folder in $archFolders) {
    $file = Get-ChildItem "$distPath\$folder\*.tar.gz"
    if ($file) {
        $hash = (Get-FileHash $file.FullName -Algorithm SHA256).Hash.ToLower()
        Write-Host "$($file.Name): $hash"
    }
}

Write-Host "`n✅ Build complete! Files located at: $releasePath" -ForegroundColor Green

# Optional: Open the folder when finished
# explorer $releasePath
