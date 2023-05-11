package main

import(
	"fmt"
	"os"
	"os/exec"
)

func printUsage(){
	panic("Format to deal with users:\n"+
		"Register <user_id>\n"+
		"Enroll <user_id> <password_generated_after_registering>\n"
	)
}

func main(){
	if len(os.Args)<2{
		printUsage()
	}

	operation:=os.Args[0]
	if operation == "Register" && len(os.Args)<2{
		printUsage()
	}
	if operation == "Enroll" && len(os.Args)<3{
		printUsage()
	}

	if operation == "Register"{
		fmt.Printf("Intiating registration of user...\n")
		userId:=os.Args[1]
		registerUser(userId)
	} else if operation == "Enroll"{
		fmt.Printf("Intiating enrollment of user...\n")
		userId:=os.Args[1]
		password:=os.Args[2]
		enrollUser(userId, password)
	}
}

func registerUser(userId string){   // this function should take in userid and print the password
	fmt.Printf("Copy the provided password for future reference\n")
	

}

func enrollUser(userId string, password string){  // this function should take in userid and password, then it should return/print the public+private key msp


}