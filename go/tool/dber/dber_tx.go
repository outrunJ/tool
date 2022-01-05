package dber

import (
	"fmt"
	"context"
)

type TxMark int

const (
	TX_MARK_NONE   TxMark = iota
	TX_MARK_DURING
	TX_MARK_DONE
)

func (d *dber) TxContext(ctx context.Context, opts *TxOptions, fn func(DBer) error, rb func() error) (err error, dbErr error, rbErr error) {

	do := func(d DBer) (err error) {
		defer func() {
			if erro := recover(); erro != nil {
				err = fmt.Errorf("%v", erro)
			}
		}()

		err = fn(d)
		return
	}

	rollback := func() (rbErr error) {
		defer func() {
			if erro := recover(); erro != nil {
				rbErr = fmt.Errorf("%v %v", rbErr, erro)
			}
		}()

		if rb != nil {
			rbErr = rb()
		}
		return
	}

	switch d.txMark {
	case TX_MARK_DONE:
		err = fmt.Errorf("tx has done")
		return
	case TX_MARK_DURING:
		err = do(d)
		if err != nil {
			rbErr = rollback()
		}
		d.txMark = TX_MARK_DONE
		return
	case TX_MARK_NONE:
		var tx Tx
		tx, err = d.db.BeginTx(ctx, opts)
		if err != nil {
			return
		}
		_d := d.clone()
		_d.SQLer = tx
		_d.txMark = TX_MARK_DURING

		err = do((_d).Context(ctx))

		d.txMark = TX_MARK_NONE
		if err != nil {
			dbErr = tx.Rollback()
			rbErr = rollback()
		} else {
			dbErr = tx.Commit()
		}

		return
	default:
		err = fmt.Errorf("tx wrong status")
		return
	}
}

func (d *dber) Tx(fn func(DBer) error, rb func() error) (err error, dbErr error, rbErr error) {
	return d.TxContext(context.Background(), nil, fn, rb)
}
