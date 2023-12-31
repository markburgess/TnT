//
// Copyright © Mark Burgess, ChiTek-i (2023)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// ****************************************************************************

package main

import (
	"fmt"
	"TnT"
)

// ***********************************************************************

func main() {

	TnT.InitializeContext()

	a := 0.1
	b := 0.2
	c := 0.3
	d := 0.4
	e := 0.5
	f := 0.6
	g := 0.7
	test := 0.0

	TnT.SetContext("a",a)
	TnT.SetContext("b",b)
	TnT.SetContext("c",c)
	TnT.SetContext("d",d)
	TnT.SetContext("e",e)
	TnT.SetContext("f",f)
	TnT.SetContext("g",g)


	// *************

	str1 := "test & ( a | b)"
	cmp1 := test * (a+b)
	expr1,res1 := TnT.ContextEval(str1)
	fmt.Println("1.",str1,"---->",expr1,res1,"CMP",cmp1,"\n")

	str2 := "(test2 & ( a | b))|(e.f.g)"
	cmp2 := (test * (a+b)) + (e*f*g)
	expr2,res2 := TnT.ContextEval(str2)
	fmt.Println("2.",str2,"---->",expr2,res2,"CMP",cmp2,"\n")

	str3 := "(test3 & ( c | d))"
	cmp3 := (test * (c+d))
	expr3,res3 := TnT.ContextEval(str3)
	fmt.Println("3.",str3,"---->",expr3,res3,"CMP",cmp3,"\n")

	str3a := "(test3a) (& ( c | d))"
	cmp3a := (test * (c+d))
	expr3a,res3a := TnT.ContextEval(str3a)
	fmt.Println("4.",str3a,"---->",expr3a,res3a,"CMP",cmp3a,"\n")

	str3b := "(test3b) & (( c | d))"
	cmp3b := (test * (c+d))
	expr3b,res3b := TnT.ContextEval(str3b)
	fmt.Println("5.",str3b,"---->",expr3b,res3b,"CMP",cmp3b,"\n")

	str4 := "!d"
	cmp4 := 0
	expr4,res4 := TnT.ContextEval(str4)
	fmt.Println("6.",str4,"---->",expr4,res4,"CMP",cmp4,"\n")

	str4a := "nosuchsymbol"
	cmp4a := 0
	expr4a,res4a := TnT.ContextEval(str4a)
	fmt.Println("7.",str4a,"---->",expr4a,res4a,"CMP",cmp4a,"\n")

	str4b := "!nosuchsymbol"
	cmp4b := 1
	expr4b,res4b := TnT.ContextEval(str4b)
	fmt.Println("7.",str4b,"---->",expr4b,res4b,"CMP",cmp4b,"\n")

	str4c := "!(nosuchsymbol)"
	cmp4c := 1
	expr4c,res4c := TnT.ContextEval(str4c)
	fmt.Println("8.",str4c,"---->",expr4c,res4c,"CMP",cmp4c,"\n")

	str4d := "!(nosuchsymbol|d)"
	cmp4d := 0
	expr4d,res4d := TnT.ContextEval(str4d)
	fmt.Println("9.",str4d,"---->",expr4d,res4d,"CMP",cmp4d,"\n")

	str4e := "!(nosuchsymbol|!d)"
	cmp4e := 1
	expr4e,res4e := TnT.ContextEval(str4e)
	fmt.Println("10.",str4e,"---->",expr4e,res4e,"CMP",cmp4e,"\n")

	str5 := " ( a ) | b&c"
	cmp5 := 1
	expr5,res5 := TnT.ContextEval(str5)
	fmt.Println("11.",str5,"---->",expr5,res5,"CMP",cmp5,"\n")


}

