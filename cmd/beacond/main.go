// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package main

import (
	"log/slog"
	"os"

	clibuilder "github.com/berachain/beacon-kit/cli/builder"
	clicomponents "github.com/berachain/beacon-kit/cli/components"
	"github.com/berachain/beacon-kit/config/spec"
	nodebuilder "github.com/berachain/beacon-kit/node-core/builder"
	"go.uber.org/automaxprocs/maxprocs"
)

// run runs the beacon node.
func run() error {
	// Set the uber max procs
	if _, err := maxprocs.Set(); err != nil {
		return err
	}

	// Build the node using the node-core.
	nb := nodebuilder.New(
		// Set the Runtime Components to the Default.
		nodebuilder.WithComponents(
			DefaultComponents(),
		),
	)

	// Build the root command using the builder
	cb := clibuilder.New(
		// Set the Name to the Default.
		clibuilder.WithName(
			"beacond",
		),
		// Set the Description to the Default.
		clibuilder.WithDescription(
			"A beacon-kit node usable with most Ethereum execution clients",
		),
		// Set the Runtime Components to the Default.
		clibuilder.WithComponents(
			clicomponents.DefaultClientComponents(),
		),
		// Set the NodeBuilderFunc to the NodeBuilder Build.
		clibuilder.WithNodeBuilderFunc(nb.Build),
		clibuilder.WithChainSpecBuilderFunc(spec.Create),
	)

	cmd, err := cb.Build()
	if err != nil {
		return err
	}

	// eventually we want to decouple from cosmos cli, and just pass in a built
	// Node and Cmd to a runner

	// for now, running the cmd will start the node
	return cmd.Run(clicomponents.DefaultNodeHome)
}

// main is the entry point.
func main() {
	if err := run(); err != nil {
		//nolint:sloglint // todo fix.
		slog.Error("startup failure", "error", err)
		os.Exit(1)
	}
}
