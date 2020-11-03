package main

import "fmt"

type Node struct {
	data int
	next *Node
}

type LinkedList struct {
	head Node
}

func (l LinkedList) countNodes() int {
	i := 1
	currentNode := l.head
	for currentNode.next != nil{
		i++
		currentNode = *currentNode.next
	}

	return i
}

func (l LinkedList) addNode(n Node) {
	currentNode := &l.head

	for currentNode.next != nil {
		currentNode = currentNode.next
	}

	currentNode.next = &n
}

func (l LinkedList) find(nb int) (i int) {
	currentNode := l.head

	for i = 1; nb != currentNode.data; i++ {
		currentNode = *currentNode.next
		if currentNode.next == nil {
			i = 0
			break
		}
	}

	return i
}

func main() {
	a := Node{data: 5}
	b := Node{data: 54}
	c := Node{data: 1}
	d := Node{data: 6543}
	e := Node{data: 543}
	f := Node{data: 4}
	g := Node{data: 12}
	h := Node{data: 63}
	i := Node{data: 54}

	a.next = &b
	b.next = &c
	c.next = &d
	d.next = &e
	e.next = &f
	f.next = &g
	g.next = &h

	list := LinkedList{head: a}

	fmt.Println(list.countNodes())
	list.addNode(i)
	fmt.Println(list.countNodes())
	fmt.Println(list.find(643))
}