package help

import (
	"fmt"
)

func Help() {
	fmt.Println("Cobble Help:\n" +
		"\tnew\t- Create a new project\n" +
		"\tupdate\t- Change a project's Script API version\n" +
		"\thelp\t- Show this help menu")
}
