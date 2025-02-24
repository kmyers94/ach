// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ach

import (
	"time"
)

// Reversal will transform a File into a Nacha compliant reversal which can be transmitted to undo fund movement.
func (f *File) Reversal(effectiveEntryDate time.Time) error {
	for i := range f.Batches {
		bh := f.Batches[i].GetHeader()

		// Must submit a Reversing Entry within a batch that includes the word "REVERSAL" in the
		// Company Entry Description field of the Company/Batch Header Record.
		//
		// The description “REVERSAL” must replace the original content of the Company Entry Description
		// field transmitted in the original batch, including content that may otherwise have been required by The Rules.
		bh.CompanyEntryDescription = "REVERSAL"

		// For each Reversing Entry, the content on the following fields must remain unchanged from the original,
		// erroneous Entry to which the Reversal relates:
		//  - Standard Entry Class Code
		//  - Company Identification/Originator Identification
		//  - Amount
		//  - Update the following records according to the fund flow

		// Adjust Effective Entry Date for same-day vs standard
		bh.EffectiveEntryDate = effectiveEntryDate.Format("060102")

		f.Batches[i].SetHeader(bh)

		// In EntryDetail records we need to update the TransactionCode fields to undo fund movement.
		entries := f.Batches[i].GetEntries()
		for j := range entries {
			switch entries[j].TransactionCode {
			case
				CheckingCredit, CheckingReturnNOCCredit, CheckingPrenoteCredit, CheckingZeroDollarRemittanceCredit,
				GLCredit, GLPrenoteCredit, GLReturnNOCCredit, GLZeroDollarRemittanceCredit,
				LoanCredit, LoanPrenoteCredit, LoanReturnNOCCredit, LoanZeroDollarRemittanceCredit,
				SavingsCredit, SavingsPrenoteCredit, SavingsReturnNOCCredit, SavingsZeroDollarRemittanceCredit:
				entries[j].TransactionCode += 5

			case
				CheckingDebit, CheckingPrenoteDebit, CheckingReturnNOCDebit, CheckingZeroDollarRemittanceDebit,
				GLDebit, GLPrenoteDebit, GLReturnNOCDebit, GLZeroDollarRemittanceDebit,
				LoanDebit, LoanReturnNOCDebit,
				SavingsDebit, SavingsPrenoteDebit, SavingsReturnNOCDebit, SavingsZeroDollarRemittanceDebit:
				entries[j].TransactionCode -= 5
			}
		}
	}
	return nil
}
