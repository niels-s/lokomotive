// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kinvolk/lokomotive/pkg/k8sutil"
	"github.com/kinvolk/lokomotive/pkg/lokomotive"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get the health of a cluster",
	Run:   runHealth,
}

// nolint:gochecknoinits
func init() {
	RootCmd.AddCommand(healthCmd)
	pf := healthCmd.PersistentFlags()
	pf.BoolVarP(&debug, "debug", "", false, "Print debug messages")
}

func runHealth(cmd *cobra.Command, args []string) {
	contextLogger := log.WithFields(log.Fields{
		"command": "lokoctl health",
		"args":    args,
	})

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	lokoConfig, diags := getLokoConfig()
	if diags.HasErrors() {
		for _, diagnostic := range diags {
			contextLogger.Error(diagnostic.Error())
		}

		contextLogger.Fatal("Errors found while loading cluster configuration")
	}

	kubeconfig, err := getKubeconfig(contextLogger, lokoConfig, true)
	if err != nil {
		contextLogger.Debugf("Error in finding kubeconfig file: %s", err)
		contextLogger.Fatal("Suitable kubeconfig file not found. Did you run 'lokoctl cluster apply' ?")
	}

	cs, err := k8sutil.NewClientset(kubeconfig)
	if err != nil {
		contextLogger.Fatalf("Error in creating Kubernetes client: %q", err)
	}

	// We can skip error checking here, as getKubeconfig() already checks it.
	p, _ := getConfiguredPlatform(lokoConfig, true)

	cluster := lokomotive.NewCluster(cs, p.Meta().ExpectedNodes)

	ns, err := cluster.GetNodeStatus()
	if err != nil {
		contextLogger.Fatalf("Error getting node status: %q", err)
	}

	ns.PrettyPrint()

	if !ns.Ready() {
		contextLogger.Fatalf("The cluster is not completely ready.")
	}

	components, err := cluster.Health()
	if err != nil {
		contextLogger.Fatalf("Error in getting Lokomotive cluster health: %q", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	// Print the header.
	fmt.Fprintln(w, "Name\tStatus\tMessage\tError\t")

	// An empty line between header and the body.
	fmt.Fprintln(w, "\t\t\t\t")

	for _, component := range components {

		// The client-go library defines only one `ComponenetConditionType` at the moment,
		// which is `ComponentHealthy`. However, iterating over the list keeps this from
		// breaking in case client-go adds another `ComponentConditionType`.
		for _, condition := range component.Conditions {
			line := fmt.Sprintf(
				"%s\t%s\t%s\t%s\t",
				component.Name, condition.Status, condition.Message, condition.Error,
			)

			fmt.Fprintln(w, line)
		}

		w.Flush()
	}
}
