package transaction

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/console"
	"polygon2ejudge/lib/ejudge_session"
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

	tmpPath       string
	ejudgeSession *ejudge_session.TEjudgeSession
	noChange      bool
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
	}

	return transaction
}

func (t *Transaction) SetNoChange() {
	t.noChange = true
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

func (t *Transaction) GetEjudgeSession() (*ejudge_session.TEjudgeSession, error) {
	if t.setupCurrent() != nil {
		return nil, t.err
	}
	if t.ejudgeSession == nil {
		session, err := ejudge_session.NewEjudgeSession(t.current.ServeCfg.ContestID)
		if err != nil {
			t.err = fmt.Errorf("can not connect to ejudge, error: %s", err.Error())
		} else {
			t.ejudgeSession = session
		}
	}
	return t.ejudgeSession, t.err
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
	t.commitMessages = append(t.commitMessages, commitMessage)
	if t.current == nil {
		return nil
	}
	t.applied.ServeCfg = t.current.ServeCfg
	t.applied.RemovePaths = append(t.applied.RemovePaths, t.current.RemovePaths...)
	t.applied.MovePaths = append(t.applied.MovePaths, t.current.MovePaths...)
	t.current = nil
	return nil
}

func (t *Transaction) Finish() {
	if t.err != nil {
		color.Set(color.FgRed)
		fmt.Printf("polygon2ejudge finished with error: %s\n", t.err.Error())
		color.Unset()
		if len(t.commitMessages) > 0 {
			color.Set(color.Bold)
			fmt.Println("Some changes were successful:")
			for _, message := range t.commitMessages {
				fmt.Println(message)
			}
			color.Unset()

			if console.YesOrNo("Apply them?") {
				t.apply()
			}
		}
	} else {
		color.Set(color.Bold)
		fmt.Println("Finished successfully:")
		for _, message := range t.commitMessages {
			fmt.Println(message)
		}
		color.Unset()

		t.apply()
	}

	err := os.RemoveAll(t.tmpPath)
	if err != nil {
		panic(err)
	}
}

func (t *Transaction) apply() {
	if t.noChange {
		return
	}
	color.Set(color.FgHiBlack)
	defer color.Unset()
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
		t.err = os.MkdirAll(t.current.TmpPath, 0774)
		return t.err
	}
	return nil
}
