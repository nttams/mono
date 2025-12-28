use aiken/collection/list
use aiken/crypto.{VerificationKeyHash}
use cardano/transaction.{OutputReference, Transaction}

pub type Datum {
  owner: VerificationKeyHash,
}

pub type Redeemer {
  msg: ByteArray,
}

validator gift {
  spend(
    datum: Option<Datum>,
    _redeemer: Redeemer,
    _own_ref: OutputReference,
    self: Transaction,
  ) {
    expect Some(Datum { owner }) = datum
    list.has(self.extra_signatories, owner)?
  }

  else(_) {
    fail
  }
}

test gift_test() {
  let datum = Datum { owner: #"0000" }

  let redeemer = Redeemer { msg: "hello" }

  let placeholder_utxo = OutputReference { transaction_id: "", output_index: 0 }

  gift.spend(
    Some(datum),
    redeemer,
    placeholder_utxo,
    // transaction.placeholder,
    Transaction { ..transaction.placeholder, extra_signatories: [datum.owner] },
  )
}
