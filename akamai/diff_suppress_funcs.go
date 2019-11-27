package akamai

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"

)

// suppressEquivalentTypeStringBoolean provides custom difference suppression for TypeString booleans
// Some arguments require three values: true, false, and "" (unspecified), but
// confusing behavior exists when converting bare true/false values with state.
func suppressEquivalentTypeStringBoolean(k, old, new string, d *schema.ResourceData) bool {
	if old == "false" && new == "0" {
		return true
	}
	if old == "true" && new == "1" {
		return true
	}
	return false
}

func suppressEquivalentJsonDiffs(k, old, new string, d *schema.ResourceData) bool {

	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func suppressEquivalentJsonRules(k, old, new string, d *schema.ResourceData) bool {

	// Deserialize and serialize through edgegrid-golang to ensure that the serialized strings are equivalent
	// This handles the case where edgegrid-golang has a different "omitEmpty" scheme to other api implementations
	//
	// When marshaling, we only consider the "Rules" part and not the header
	//
	// Note: if this function determines that the two rule sets are different, Terraform will show ALL
	// differences in the plan, even those that considered trivial
	//

        nrules := papi.NewRules()
        orules := papi.NewRules()

	if err := json.Unmarshal([]byte(old), orules); err != nil {
		return false
	}
	nold, err := json.Marshal(orules.Rule);
	if (err != nil) {
		return false
	}

	if err := json.Unmarshal([]byte(new), nrules); err != nil {
		return false
	}
	nnew, err := json.Marshal(nrules.Rule);
	if (err != nil) {
		return false
	}

	return suppressEquivalentJsonDiffs(k, string(nold), string(nnew), d)
}

