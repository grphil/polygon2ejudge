package import_problem

import (
	"fmt"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
	"strings"
)

func CreateImportTaskForUpdate(
	taskCommon common.TaskCommon,
	transaction *transaction2.Transaction,
	ejudgeProblemID int,
) (*ImportTask, error) {
	serveCfg, err := transaction.EditServeCfg()
	if err != nil {
		return nil, err
	}

	prob, ok := serveCfg.Problems[ejudgeProblemID]
	if !ok {
		transaction.SetError(fmt.Errorf("problem %d not found", ejudgeProblemID))
		return nil, transaction.Err()
	}

	extID := prob.GetStr("extid")

	var problemUrl string
	var polygonID int

	if strings.HasPrefix(extID, "https://polygon.codeforces.com/") {
		problemUrl = extID
	} else {
		n, err := fmt.Sscanf(extID, "polygon:%d", &polygonID)
		if err != nil {
			transaction.SetError(fmt.Errorf("can not get problem %d polygon id, error: %s", ejudgeProblemID, err.Error()))
			return nil, err
		}
		if n == 0 {
			transaction.SetError(fmt.Errorf("can not get problem %d polygon id", ejudgeProblemID))
			return nil, transaction.Err()
		}
	}

	shortName := prob.GetStr("short_name")
	internalName := prob.GetStr("internal_name")
	importTask := &ImportTask{
		TaskCommon:   taskCommon,
		ShortName:    &shortName,
		EjudgeId:     &ejudgeProblemID,
		Transaction:  transaction,
		InternalName: internalName,
	}

	if problemUrl != "" {
		importTask.PolygonProbUrl = &problemUrl
	}
	importTask.PolygonProbID = &polygonID
	return importTask, nil
}
