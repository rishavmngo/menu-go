package menu

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkg/term"
	"github.com/rishavmngo/menu-go/models"
)

type Node struct {
	name      string
	parent    *Node
	childrens []*Node
	action    func()
}

type Menu struct {
	Main    *Node
	running bool
}

func (menu *Menu) IsRunning() bool {
	return menu.running
}

func (menu *Menu) Exit() {

	menu.running = false

}

func (menu *Menu) Display() {
	head := menu.Main
	var buffer bytes.Buffer

	var currentItem = models.Get()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

mainLoop:
	for {

		if !menu.running {
			cancel()
			break mainLoop

		}

		ClearScreenStandalone()
		headingOfList(head, &buffer, currentItem)
		getListItems(head, &buffer, currentItem)

		_, err := buffer.WriteTo(os.Stdout)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing buffer to stdout: %v\n", err)
		}
		inp := getInput()

		switch inp {
		case '\n', 13:
			expand(&head, currentItem)
		case 3:
			cancel()
			break mainLoop
		case 127: //backspace
		case 23:
		case 32: //space
		case 'j':
			currentItem.Increment()
		case 'k':
			currentItem.Decrement()
		case 'u':
			currentItem.Decrement()
		case 'd':
			currentItem.Increment()
		case '<':
			collapse(&head, currentItem)
		case 'l':
			expand(&head, currentItem)
		case 'h':
			collapse(&head, currentItem)
		case '>':
			expand(&head, currentItem)
		default:
		}
	}

}
func getInput() byte {

	t, _ := term.Open("/dev/tty")

	err := term.RawMode(t)
	if err != nil {
		log.Fatal(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)

	t.Restore()
	t.Close()
	if read == 3 {

		switch readBytes[2] {
		case 'A':
			return 'u' // Up arrow
		case 'B':
			return 'd' // Down arrow
		case 'C':
			return '>' // right arrow
		case 'D':
			return '<' // left arrow
		}

	}

	return readBytes[0]

}

func NewMenu(name string) *Menu {

	head := &Node{name: name}

	return &Menu{Main: head, running: true}

}

func (node *Node) Add(name string, action func()) *Node {

	newNode := &Node{name: name}
	newNode.parent = node
	newNode.action = action
	node.childrens = append(node.childrens, newNode)
	return newNode

}
