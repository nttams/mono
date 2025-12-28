use aiken/interval.{Interval}
use cardano/transaction.{OutputReference, Transaction}

pub type Datum {
  valid_after: Int,
}

validator vet {
  spend(
    datum: Option<Datum>,
    _redeemer: Void,
    _utxo: OutputReference,
    tx: Transaction,
  ) {
    expect Some(Datum { valid_after }) = datum

    let t: Interval = interval.after(valid_after)
    tx.validity_range |> interval.includes(t)
  }

  else(_) {
    fail
  }
}

test vet_test() {
  // will fail
  // let datum = Datum { valid_after: 1766760002000 }

  // will pass
  let datum = Datum { valid_after: 1766760002002 }

  let placeholder_utxo = OutputReference { transaction_id: "", output_index: 0 }

  vet.spend(
    Some(datum),
    Void,
    placeholder_utxo,
    Transaction {
      ..transaction.placeholder,
      validity_range: interval.after(1766760002001),
    },
  )
}
