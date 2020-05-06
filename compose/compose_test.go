package compose

import (
	"fmt"
	"testing"
)

func TestCompose(t *testing.T) {
	fmt.Println("start")
	project := NewProject("test", "./docker-compose.yml")
	fmt.Println(project.Services)
	for k, v := range project.Services {
		fmt.Printf("start %s", k)
		fmt.Println(v.Up(project.Client))
	}
	// if db, ok := project.Services["db"]; ok {
	// 	t.Logf("start %s", db.Name)
	// 	db.up(project.Client)
	// }
}
