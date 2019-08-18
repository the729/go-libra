/*
Package types implements core Libra data structs.

Some structs are 'proven', namely `ProvenLedgerInfo`, `ProvenAccountState`, `ProvenAccountBlob`,
`ProvenTransactionList` and `ProvenTransaction`.

These proven structs can only be created by verification of corresponding crypto proofs. And they
do not export any member values. You can only access the member value by getter functions, which
guarantees that no member value can be changed. Thus, you can stay asured that a proven struct is genuine.
Be careful, however, that a genuine proven data struct can be outdated. Always remember to check the
ledger version.

Not proven-prefixed structs have all member values exported. You can create a data struct and its proofs
either by gRPC queries, or by any other means such as loading from file.
*/
package types
