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

	const some_test = true
	const many = 3
	const policy_condition = "(state_of_contention) | state_flag & state_of_uncertainty"

	TnT.InitializeContext()	// Reset context set

	for transactions := 1; transactions <= many; transactions++ { 

		if some_test {

			TnT.ContextActive("state_of_uncertainty")
		}
		
		if TnT.IsDefinedContext(policy_condition) {

			fmt.Println("ACTIVE POLICY",TnT.ContextSet())
		} else {
			fmt.Println("No matching context",TnT.ContextSet())
		}

		TnT.ContextActive("state_flag")
	}
}

