package ualert

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCondTransExtended(t *testing.T) {
	// Here we are testing that we got a query that is referenced with multiple different offsets, the migration
	// generated correctly all subqueries for each offset. RefID A exists twice with a different offset (cond1, cond4).

	ordID := int64(1)
	lookup := dsUIDLookup{}

	settings := dashAlertSettings{}

	cond1 := dashAlertCondition{}
	cond1.Evaluator.Params = []float64{-500000}
	cond1.Evaluator.Type = "lt"
	cond1.Operator.Type = "and"
	cond1.Query.DatasourceID = 4
	cond1.Query.Model = []byte(`{"datasource":{"type":"graphite","uid":"1"},"hide":false,"refCount":0,"refId":"A","target":"my_metric_1","textEditor":true}`)
	cond1.Query.Params = []string{
		"A",
		"1h",
		"now",
	}
	cond1.Reducer.Type = "diff"

	cond2 := dashAlertCondition{}
	cond2.Evaluator.Params = []float64{
		-0.01,
		0.01,
	}
	cond2.Evaluator.Type = "within_range"
	cond2.Operator.Type = "or"
	cond2.Query.DatasourceID = 4
	cond2.Query.Model = []byte(`{"datasource":{"type":"graphite","uid":"1"},"hide":true,"refCount":0,"refId":"B","target":"my_metric_2","textEditor":false}`)
	cond2.Query.Params = []string{
		"B",
		"6h",
		"now",
	}
	cond2.Reducer.Type = "diff"

	cond3 := dashAlertCondition{}
	cond3.Evaluator.Params = []float64{
		-500000,
	}
	cond3.Evaluator.Type = "lt"
	cond3.Operator.Type = "or"
	cond3.Query.DatasourceID = 4
	cond3.Query.Model = []byte(`{"datasource":{"type":"graphite","uid":"1"},"hide":false,"refCount":0,"refId":"C","target":"my_metric_3","textEditor":false}`)
	cond3.Query.Params = []string{
		"C",
		"1m",
		"now",
	}
	cond3.Reducer.Type = "diff"

	cond4 := dashAlertCondition{}
	cond4.Evaluator.Params = []float64{
		1000000,
	}
	cond4.Evaluator.Type = "gt"
	cond4.Operator.Type = "and"
	cond4.Query.DatasourceID = 4
	cond4.Query.Model = []byte(`{"datasource":{"type":"graphite","uid":"1"},"hide":false,"refCount":0,"refId":"A","target":"my_metric_1","textEditor":true}`)
	cond4.Query.Params = []string{
		"A",
		"5m",
		"now",
	}
	cond4.Reducer.Type = "last"

	settings.Conditions = []dashAlertCondition{cond1, cond2, cond3, cond4}

	alertQuery1 := alertQuery{
		RefID: "A",
		RelativeTimeRange: relativeTimeRange{
			From: 3600000000000,
		},
		Model: cond1.Query.Model,
	}
	alertQuery2 := alertQuery{
		RefID: "B",
		RelativeTimeRange: relativeTimeRange{
			From: 300000000000,
		},
		Model: []byte(strings.ReplaceAll(string(cond1.Query.Model), "refId\":\"A", "refId\":\"B")),
	}
	alertQuery3 := alertQuery{
		RefID: "C",
		RelativeTimeRange: relativeTimeRange{
			From: 21600000000000,
		},
		Model: []byte(strings.ReplaceAll(string(cond2.Query.Model), "refId\":\"B", "refId\":\"C")),
	}
	alertQuery4 := alertQuery{
		RefID: "D",
		RelativeTimeRange: relativeTimeRange{
			From: 60000000000,
		},
		Model: []byte(strings.ReplaceAll(string(cond3.Query.Model), "refId\":\"C", "refId\":\"D")),
	}
	alertQuery5 := alertQuery{
		RefID:         "E",
		DatasourceUID: "__expr__",
		Model:         []byte(`{"type":"classic_conditions","refId":"E","conditions":[{"evaluator":{"params":[-500000],"type":"lt"},"operator":{"type":"and"},"query":{"params":["A"]},"reducer":{"type":"diff"}},{"evaluator":{"params":[-0.01,0.01],"type":"within_range"},"operator":{"type":"or"},"query":{"params":["C"]},"reducer":{"type":"diff"}},{"evaluator":{"params":[-500000],"type":"lt"},"operator":{"type":"or"},"query":{"params":["D"]},"reducer":{"type":"diff"}},{"evaluator":{"params":[1000000],"type":"gt"},"operator":{"type":"and"},"query":{"params":["B"]},"reducer":{"type":"last"}}]}`),
	}

	expected := &condition{
		Condition: "E",
		OrgID:     ordID,
		Data:      []alertQuery{alertQuery1, alertQuery2, alertQuery3, alertQuery4, alertQuery5},
	}

	c, err := transConditions(settings, ordID, lookup)

	require.NoError(t, err)
	require.Equal(t, expected, c)
}
