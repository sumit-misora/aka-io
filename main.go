package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wmnsk/milenage"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: aka-io K OPc RAND AMF SQN\n")
		fmt.Fprintf(os.Stderr, "All inputs must be in hexadecimal format\n")
		fmt.Fprintf(os.Stderr, "K: 32 hex chars (16 bytes)\n")
		fmt.Fprintf(os.Stderr, "OPc: 32 hex chars (16 bytes)\n")
		fmt.Fprintf(os.Stderr, "RAND: 32 hex chars (16 bytes)\n")
		fmt.Fprintf(os.Stderr, "AMF: 4 hex chars (2 bytes)\n")
		fmt.Fprintf(os.Stderr, "SQN: 12 hex chars (6 bytes)\n")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 5 {
		flag.Usage()
		os.Exit(1)
	}

	// Parse and validate inputs
	k, err := parseAndValidate(args[0], 32, "K")
	if err != nil {
		log.Fatalf("Invalid K: %v", err)
	}

	opc, err := parseAndValidate(args[1], 32, "OPc")
	if err != nil {
		log.Fatalf("Invalid OPc: %v", err)
	}

	randBytes, err := parseAndValidate(args[2], 32, "RAND")
	if err != nil {
		log.Fatalf("Invalid RAND: %v", err)
	}

	amfBytes, err := parseAndValidate(args[3], 4, "AMF")
	if err != nil {
		log.Fatalf("Invalid AMF: %v", err)
	}

	sqnBytes, err := parseAndValidate(args[4], 12, "SQN")
	if err != nil {
		log.Fatalf("Invalid SQN: %v", err)
	}

	// Convert to proper types
	// AMF: 2 bytes to uint16 (big-endian)
	amf := binary.BigEndian.Uint16(amfBytes)

	// SQN: 6 bytes to uint64 (big-endian, padded with 2 leading zero bytes)
	sqnPadded := make([]byte, 8)
	copy(sqnPadded[2:], sqnBytes)
	sqn := binary.BigEndian.Uint64(sqnPadded)

	// Display input
	fmt.Println("[INPUT]")
	fmt.Printf("K     : %s\n", strings.ToUpper(args[0]))
	fmt.Printf("OPc   : %s\n", strings.ToUpper(args[1]))
	fmt.Printf("RAND  : %s\n", strings.ToUpper(args[2]))
	fmt.Printf("AMF   : %s\n", strings.ToUpper(args[3]))
	fmt.Printf("SQN   : %s\n", strings.ToUpper(args[4]))

	// Create Milenage instance
	m := milenage.NewWithOPc(k, opc, randBytes, sqn, amf)

	// Compute all values (F1, F2, F3, F4, F5, F1*, F5*)
	err = m.ComputeAll()
	if err != nil {
		log.Fatalf("ComputeAll calculation failed: %v", err)
	}

	// Read computed values from Milenage struct
	macA := m.MACA
	macS := m.MACS
	res := m.RES
	ck := m.CK
	ik := m.IK
	ak := m.AK
	aks := m.AKS

	// Manually construct AUTN according to 3GPP TS 33.102 specification
	// AUTN = (SQN XOR AK) || AMF || MAC-A
	autnManual := make([]byte, 16)

	// First 6 bytes: SQN XOR AK
	for i := 0; i < 6; i++ {
		autnManual[i] = sqnBytes[i] ^ ak[i]
	}

	// Next 2 bytes: AMF
	copy(autnManual[6:8], amfBytes)

	// Last 8 bytes: MAC-A
	copy(autnManual[8:16], macA)

	// Generate AUTS using library method (which recalculates MAC-S and AK-S with AMF=0x0000)
	auts, err := m.GenerateAUTS()
	if err != nil {
		log.Fatalf("AUTS generation failed: %v", err)
	}

	// Display output
	fmt.Println("[OUTPUT]")
	fmt.Printf("MAC-A : %s\n", strings.ToUpper(hex.EncodeToString(macA)))
	fmt.Printf("MAC-S : %s\n", strings.ToUpper(hex.EncodeToString(macS)))
	fmt.Printf("RES   : %s\n", strings.ToUpper(hex.EncodeToString(res)))
	fmt.Printf("CK    : %s\n", strings.ToUpper(hex.EncodeToString(ck)))
	fmt.Printf("IK    : %s\n", strings.ToUpper(hex.EncodeToString(ik)))
	fmt.Printf("AK    : %s\n", strings.ToUpper(hex.EncodeToString(ak)))
	fmt.Printf("AKS   : %s\n", strings.ToUpper(hex.EncodeToString(aks)))
	fmt.Printf("AUTN  : %s\n", strings.ToUpper(hex.EncodeToString(autnManual)))
	fmt.Printf("AUTS  : %s\n", strings.ToUpper(hex.EncodeToString(auts)))
}

// parseAndValidate converts a hex string to bytes and validates its length
func parseAndValidate(hexStr string, expectedLen int, fieldName string) ([]byte, error) {
	hexStr = strings.TrimSpace(hexStr)
	if len(hexStr) != expectedLen {
		return nil, fmt.Errorf("%s must be %d hex characters, got %d", fieldName, expectedLen, len(hexStr))
	}

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex format: %v", err)
	}

	return decoded, nil
}
