package testing

import (
	"bufio"
	"fmt"
	"os"
  "testing"
)

func TestXxx(*testing.T){

}

func main()  {
  var runa rune
  var s string


  f, err := os.Open(file)
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)
	lex, err := NewLexer(r, file)
	fmt.Println(lex)

  for true{
    runa, err = lex.get()
  	s = string(runa)
  	fmt.Println(s)
  }

}
