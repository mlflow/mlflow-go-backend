package sql

import (
	"fmt"
	"strings"

	"github.com/mlflow/mlflow-go-backend/pkg/contract"
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
)

// HandleResourceAlreadyExistError handles name conflicts in the Model Registry.
func HandleResourceAlreadyExistError(name string, isExistingEntityPrompt, isNewEntityPrompt bool) *contract.Error {
	// Determine the entity types
	oldEntity := "Registered Model"
	if isExistingEntityPrompt {
		oldEntity = "Prompt"
	}

	newEntity := "Registered Model"
	if isNewEntityPrompt {
		newEntity = "Prompt"
	}

	// Check if there is a conflict between different entity types
	if oldEntity != newEntity {
		return contract.NewError(
			protos.ErrorCode_RESOURCE_ALREADY_EXISTS,
			fmt.Sprintf(
				"Tried to create a %s with name '%s', but the name is already taken by a %s. "+
					"MLflow does not allow creating a model and a prompt with the same name.",
				strings.ToLower(newEntity), name, strings.ToLower(oldEntity)),
		)
	}

	// Raise an error if the entity already exists
	return contract.NewError(
		protos.ErrorCode_RESOURCE_ALREADY_EXISTS,
		fmt.Sprintf("%s (name=%s) already exists.", newEntity, name),
	)
}
