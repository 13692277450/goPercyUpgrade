<img width="554" height="201" alt="image" src="https://github.com/user-attachments/assets/a1dad592-6641-4003-a96e-2a6c14ede96c" />


# GoPercyUpgrade

GoPercyUpgrade is a sample Go library that only need two lines code to upgrade your application for 3 platforms: Windows, Linux, MacOS.

## Features

- Simple and clean API, only two lines code to upgrade your application
- Auto detect platform and choose correct download url
- Auto move current application to a new name to replace application from new one.
- Auto remove old files if needed.

## Installation

```bash
go get github.com/13692277450/gopercyupgrade@latest
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

## Screenshots


Linux version test results:

<img width="2124" height="311" alt="image" src="https://github.com/user-attachments/assets/659b9118-dcbb-41f8-a7d2-23509f1f3629" />


MacOS version test results:

 <img width="2123" height="508" alt="image" src="https://github.com/user-attachments/assets/92452709-44c4-4e69-9769-3100a601040a" />


Windows version test results:

 <img width="1700" height="326" alt="image" src="https://github.com/user-attachments/assets/e0d47467-2566-4cd0-a0d0-10f09cf58277" />





## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change, or you can email m13692277450@outlook.com or call +86-13692277450 if needed.

## License

[MIT](LICENSE)
