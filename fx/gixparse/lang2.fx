// Lenguage example for testing gixparse. The iter sentence is not right, it need one more argument.

//macro definition
func line(int x, int y){
				//last number in loop is the step
	iter (i := 0; x){	//declares it, scope is the loop
		circle(2, 3, y, 5);
	}
}

//macro entry
func main(){
	iter (i := 0; 3, 1){
		rect(i, i, 3, 0xff);
	}
	iter (j := 0; 8, 2){	//loops 0 2 4 6 8
		rect(j, j, 8, 0xff);
	}
	circle(4, 5, 2, 0x11000011);
}
