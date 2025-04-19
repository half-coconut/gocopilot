package day02

import (
	"fmt"
	"io"
	"testing"
)

type IGreeting interface {
	sayHello() string
}

//type iGreeting struct {
//	g Go
//	p PHP
//}
//
//func sayHello(i IGreeting) {
//	i.sayHello()
//}

type Go struct{}

func (g *Go) sayHello() string {
	return "hi, go"
}

type PHP struct{}

func (p *PHP) sayHello() string {
	return "hi, php"
}

func TestInterfaces(t *testing.T) {
	g := Go{}
	p := PHP{}
	t.Log(g.sayHello())
	t.Log(p.sayHello())
}

type Person interface {
	whatJob()
	growUp()
}

type student struct {
	age int
}

func (s student) whatJob() {
	fmt.Println("I'm a student")
}

// 注意：这里使用指针后，会改变 s.age 的值
func (s *student) growUp() {
	s.age++
}

type programmer struct {
	age int
}

func (p programmer) whatJob() {
	fmt.Println("I'm a student")
}

// 这里不使用指针，只是对p.age 的值传递，即为值的副本
func (p programmer) growUp() {
	p.age += 10
}

func TestPerson(t *testing.T) {
	s := &student{
		age: 18,
	}
	s.whatJob()
	s.growUp()
	t.Log(s.age)

	p := programmer{
		age: 100,
	}
	p.whatJob()
	p.growUp()
	t.Log(p.age)
}

func TestNil(t *testing.T) {
	type apple struct{}
	a := &apple{}
	t.Log(a == nil)
}

func TestFloatToInt(t *testing.T) {
	var i = 9
	var f float64
	f = float64(i)
	t.Log(fmt.Sprintf("%T,%v\n", f, f))

	f = 10.8
	a := int(f)
	t.Log(fmt.Sprintf("%T,%v\n", a, a))
}

type MyWriter struct {
}

func (w MyWriter) Write(p []byte) (n int, err error) {
	return
}

// 编译器自动检测类型是否实现了接口
func TestValidType(t *testing.T) {
	var _ io.Writer = (*MyWriter)(nil)

	var _ io.Writer = MyWriter{}
}
