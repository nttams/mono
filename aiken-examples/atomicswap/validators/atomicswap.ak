use aiken/collection/list
use cardano/address.{Address}
use cardano/assets.{lovelace_of}
use cardano/transaction.{OutputReference, Transaction}

pub type Datum {
  receiver: Address,
  amount_lovelace: Int,
}

// NOT DONE YET
// - Write test
// - Is using list.any a good idea?
// - the receiver: Address in type Datum, will it be serialized/deserialized correctly?
validator atomicswap {
  spend(
    datum: Option<Datum>,
    _redeemer: Data,
    _utxo: OutputReference,
    self: Transaction,
  ) {
    expect Some(Datum { receiver, amount_lovelace }) = datum

    list.any(
      self.outputs,
      fn(o) { o.address == receiver && lovelace_of(o.value) > amount_lovelace },
    )?
  }

  else(_) {
    fail
  }
}
