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

	const many = 3
	const name = "my promised transaction"

	for transactions := 1; transactions < many; transactions++ { 
		
		ctx := TnT.PromiseContext_Begin(name)

		fmt.Println("Do something atomic")
		
		TnT.PromiseContext_End(ctx)
	}
}

