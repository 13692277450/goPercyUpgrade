# GoPercyGrade

GoPercyGrade is a sample Go library that only need two lines code to upgrade your application for 3 platform: Windows, Linux, MacOS.

## Features

- Simple and clean API, only two lines code to upgrade your application
- Auto detect platform and choose correct download url
- Auto move current application to a new name to replace application from new one.
- Auto remove old files if needed.

## Installation

```bash
go get github.com/13692277450/gopercyupgrde
```

## Usage

```go
import "github.com/13692277450/gopercyupgrade"

func main() {
    // Use the library functions here
    var currentVersion string = "1.0.1"
    GoPercyUpgradeConfig(currentVersion, "http://example.com/version.json")
    /*
    put version.json url in your upgrade server.
    Version.json file structure:
    {
    "versionwindows": "1.0.1",
    "versionlinux": "1.0.1",
    "versionmac": "1.0.1",
    "downloadUrlwindows": "http://example.com/windows",
    "downloadUrllinux": "http://example.com/linux",
    "downloadUrlmac": "http://example.com/mac",
    "noteswindows": "Bug fixes and improvements",
    "noteslinux": "Bug fixes and improvements",
    "notesmac": "Bug fixes and improvements",
    "pub_date": "2025-11-09T00:00:00Z",
    "shownotesmessages": true,
    "removeoldfiles": true
    }
    */
    //done, so easy.
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change, or you can email m13692277450@outlook.com or call +86-13692277450 if needed.

## License

[MIT](LICENSE)
