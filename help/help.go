package help

import (
	"fmt"
)

func Help() {
	fmt.Println("Usage:\n" +
		"\tcobble new [name]\t- Creates a new project interactively.\n" +
		"\tcobble help\t- Show this help menu")
}
