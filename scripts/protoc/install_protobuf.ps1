# Specify the URL of the protobuf compiler (protoc) zip file
$url = "https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip"

# Define the download path
$downloadPath = "$env:USERPROFILE\Downloads\protoc.zip"

# Define the extraction path
$extractPath = "$env:USERPROFILE\protobuf"

# Define the bin directory which contains the protoc binary
$binPath = "$extractPath\bin"

# Check if the extraction path already exists
if (-Not (Test-Path -Path $extractPath)) {
  # Download the zip file
  Invoke-WebRequest -Uri $url -OutFile $downloadPath

  # Extract the zip file
  Expand-Archive -Path $downloadPath -DestinationPath $extractPath

  Write-Host "Protobuf compiler (protoc) has been downloaded and extracted."
} else {
  Write-Host "Protobuf compiler (protoc) already exists at $extractPath. Skipping download and extraction."
}

# Add the bin directory to the system PATH environment variable if it's not already there
if ($env:Path -notcontains $binPath) {
  $env:Path += ";$binPath"

  # To make it permanent, update the user's PATH environment variable
  [Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::User)

  Write-Host "Protobuf compiler (protoc) has been added to PATH."
} else {
  Write-Host "Protobuf compiler (protoc) is already in PATH."
}
