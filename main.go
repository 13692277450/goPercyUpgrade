package gopercyupgrade

import "fmt"

func main() {
	var currentversion string = "1.0.0"
	fmt.Println("Welcome to use GoPercyUpgrade!")
	GoPercyUpgradeConfig(currentversion, "http://example.com/version.json")

}
