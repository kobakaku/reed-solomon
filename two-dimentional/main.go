package main

import (
	"bytes"
	"log"

	"github.com/celestiaorg/rsmt2d"
)

func main() {
	log.Println("🚀 Starting the process...")

	// shareSize is the size of each share (in bytes).
	shareSize := 512
	// Init new codec
	codec := rsmt2d.NewLeoRSCodec()

	ones := bytes.Repeat([]byte{1}, shareSize)
	twos := bytes.Repeat([]byte{2}, shareSize)
	threes := bytes.Repeat([]byte{3}, shareSize)
	fours := bytes.Repeat([]byte{4}, shareSize)

	// Compute parity shares
	eds, err := rsmt2d.ComputeExtendedDataSquare(
		[][]byte{
			ones, twos,
			threes, fours,
		},
		codec,
		rsmt2d.NewDefaultTree,
	)
	if err != nil {
		log.Fatalf("❌ ComputeExtendedDataSquare failed: %v", err)
	}

	rowRoots, err := eds.RowRoots()
	if err != nil {
		log.Fatalf("❌ Return the Merkle roots of all the rows in the square failed: %v", err)
	}
	colRoots, err := eds.ColRoots()
	if err != nil {
		log.Fatalf("❌ returns the Merkle roots of all the columns in the square failed: %v", err)
	}

	flattened := eds.Flattened()

	// Delete some shares, just enough so that repairing is possible.
	flattened[0], flattened[2], flattened[3] = nil, nil, nil
	flattened[4], flattened[5], flattened[6], flattened[7] = nil, nil, nil, nil
	flattened[8], flattened[9], flattened[10] = nil, nil, nil
	flattened[12], flattened[13] = nil, nil

	// Re-import the data square.
	eds, err = rsmt2d.ImportExtendedDataSquare(flattened, codec, rsmt2d.NewDefaultTree)
	if err != nil {
		log.Fatalf("❌ Reimport the data square failed: %v", err)
	}

	// Repair square.
	err = eds.Repair(
		rowRoots,
		colRoots,
	)
	if err != nil {
		log.Fatalf("❌ Repair an incomplete extended data square (EDS) failed: %v", err)
	}

	log.Println("✅ The process completed successfully.")
}
