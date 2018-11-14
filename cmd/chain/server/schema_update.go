package server

import (
    "fmt"
    "github.com/ssor/sql2graphql/cmd/chain/pipeline"

    "github.com/ssor/sql2graphql/helper"
)

func newSchemaUpdateStep() *SchemaUpdateStep {
    return &SchemaUpdateStep{}
}

type SchemaUpdateStep struct {
    schemas string
}

func (su *SchemaUpdateStep) Exec(request *pipeline.Request) *pipeline.Result {
    serverLogger.Infof("try to alter schema")

    var schemas helper.Schemas
    schemas = schemas.Add((&Ledger{}).Schemes()...)
    schemas = schemas.Add((&Block{}).Schemes()...)
    schemas = schemas.Add((&Transaction{}).Schemes()...)
    //serverLogger.Info("schemas : \n", schemas.String())
    su.schemas = schemas.String()

    err := helper.Alter(schemas, dgClient)
    if err != nil {
        logger.Failedf("alter schema failed: %s", err)
    }
    return &pipeline.Result{
        Error: err,
    }
}

func (su *SchemaUpdateStep) WorkDescription() (s string) {
    s = fmt.Sprintf("schemas update: \n [%s] \n", su.schemas)
    return
}
