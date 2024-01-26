package transaction

import (
	"fmt"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/console"
	"polygon2ejudge/lib/serve_cfg"
	"strconv"
)

type movePathAction struct {
	src string
	dst string
}

type transactionActions struct {
	RemovePaths []string
	MovePaths   []movePathAction

	ServeCfg *serve_cfg.ServeCFG
	TmpPath  string
}

type Transaction struct {
	applied transactionActions
	current *transactionActions

	err            error
	commitMessages []string

	tmpPath string
}

func NewTransaction(contestID int) *Transaction {
	serveCfg, err := serve_cfg.NewServeCFG(contestID)
	if err != nil {
		err = fmt.Errorf("can not parse serve.cfg, error: %s", err.Error())
	}

	transaction := &Transaction{
		applied: transactionActions{
			RemovePaths: nil,
			MovePaths:   nil,
			ServeCfg:    serveCfg,
		},
		err:            err,
		commitMessages: nil,
	}

	if transaction.err == nil {
		transaction.tmpPath = filepath.Join(transaction.applied.ServeCfg.Path(), "tmp")
		err = os.RemoveAll(transaction.tmpPath)
		if err != nil {
			transaction.err = fmt.Errorf("can not clear tmp path %s, err: %s", transaction.tmpPath, err.Error())
		}
		err = os.Mkdir(transaction.tmpPath, 0774)
		if err != nil {
			transaction.err = fmt.Errorf("can not create tmp path %s, err: %s", transaction.tmpPath, err.Error())
		}
	}

	return transaction
}

func (t *Transaction) GetTmp() (string, error) {
	if t.setupCurrent() != nil {
		return "", t.err
	}

	return t.current.TmpPath, nil
}

func (t *Transaction) EditServeCfg() (*serve_cfg.ServeCFG, error) {
	if t.setupCurrent() != nil {
		return nil, t.err
	}
	return t.current.ServeCfg, nil
}

func (t *Transaction) RemovePath(path string) error {
	if t.setupCurrent() != nil {
		return t.err
	}

	t.current.RemovePaths = append(t.current.RemovePaths, path)
	return nil
}

func (t *Transaction) MovePath(src string, dst string) error {
	if t.setupCurrent() != nil {
		return t.err
	}

	err := os.MkdirAll(filepath.Dir(dst), 0774)
	if err != nil {
		t.err = fmt.Errorf("can not move path %s to %s, error: %s", src, dst, err)
		return t.err
	}

	t.current.MovePaths = append(t.current.MovePaths, movePathAction{src: src, dst: dst})
	return nil
}

func (t *Transaction) SetError(err error) {
	if t.err == nil {
		t.err = err
	}
}

func (t *Transaction) Err() error {
	return t.err
}

func (t *Transaction) Commit(commitMessage string) error {
	if t.err != nil {
		return t.err
	}
	if t.current == nil {
		return nil
	}
	t.applied.ServeCfg = t.current.ServeCfg
	t.applied.RemovePaths = append(t.applied.RemovePaths, t.current.RemovePaths...)
	t.applied.MovePaths = append(t.applied.MovePaths, t.current.MovePaths...)
	t.commitMessages = append(t.commitMessages, commitMessage)
	return nil
}

func (t *Transaction) Finish() {
	if t.err != nil {
		fmt.Printf("polygon2ejudge finished with error: %s\n", t.err.Error())
		if len(t.commitMessages) > 0 {
			fmt.Println("Some changes were successful:")
			for _, message := range t.commitMessages {
				fmt.Println(message)
			}

			if console.YesOrNo("Apply them?") {
				t.apply()
			}
		}
	} else {
		fmt.Println("Finished successfully:")
		for _, message := range t.commitMessages {
			fmt.Println(message)
		}

		t.apply()
	}

	err := os.RemoveAll(t.tmpPath)
	if err != nil {
		panic(err)
	}
}

func (t *Transaction) apply() {
	fmt.Println("Applying changes")

	err := t.applied.ServeCfg.Write()
	if err != nil {
		fmt.Println("Can not write serve.cfg")
		panic(err)
	}
	fmt.Println("Written serve.cfg")

	for _, path := range t.applied.RemovePaths {
		err = os.RemoveAll(path)
		if err != nil {
			fmt.Printf("Can not remove path %s err: %s\n", path, err.Error())
			panic(err)
		}
		fmt.Printf("Removed path %s\n", path)
	}

	for _, moveAction := range t.applied.MovePaths {
		err = os.RemoveAll(moveAction.dst)
		if err != nil {
			panic(err)
		}

		err = os.Rename(moveAction.src, moveAction.dst)
		if err != nil {
			fmt.Printf("Can not move path %s to %s error: %s\n", moveAction.src, moveAction.dst, err.Error())
			panic(err)
		}
		fmt.Printf("Moved path %s to %s\n", moveAction.src, moveAction.dst)
	}
}

func (t *Transaction) setupCurrent() error {
	if t.err != nil {
		return t.err
	}
	if t.current == nil {
		t.current = &transactionActions{
			RemovePaths: nil,
			MovePaths:   nil,
			ServeCfg:    t.applied.ServeCfg.Clone(),
			TmpPath:     filepath.Join(t.tmpPath, strconv.Itoa(len(t.commitMessages))),
		}
		return os.Mkdir(t.current.TmpPath, 0774)
	}
	return nil
}
