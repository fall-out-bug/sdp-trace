package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

// isBinaryFile reads the first 512 bytes of path and checks for known binary
// magic numbers (ELF, Mach-O, PE/COFF). It returns false for empty or
// unreadable files so that missing paths do not create false positives.
func isBinaryFile(path string) bool {
	buf, ok := readFileHeader(path)
	if !ok {
		return false
	}
	return isKnownBinaryFormat(buf)
}

// readFileHeader returns up to the first 512 bytes of path. The bool
// indicates whether any bytes were successfully read.
func readFileHeader(path string) ([]byte, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil || n == 0 {
		return nil, false
	}
	return buf[:n], true
}

// isKnownBinaryFormat returns true when buf starts with a binary executable
// signature. The set of recognised formats is intentionally small: ELF for
// Linux, Mach-O for macOS, and PE/COFF for Windows.
func isKnownBinaryFormat(buf []byte) bool {
	return isELF(buf) || isMachO(buf) || isPE(buf)
}

// isELF reports whether buf begins with the ELF magic bytes.
func isELF(buf []byte) bool {
	return len(buf) >= 4 && bytes.Equal(buf[:4], []byte{0x7f, 'E', 'L', 'F'})
}

// isMachO reports whether buf begins with a Mach-O header signature. Both
// 32-bit and 64-bit headers are accepted, as well as universal binaries.
// The check uses big-endian decoding because the magic constants are stored
// in network byte order regardless of the host architecture.
func isMachO(buf []byte) bool {
	if len(buf) < 4 {
		return false
	}
	switch binary.BigEndian.Uint32(buf[:4]) {
	case 0xfeedface, 0xfeedfacf, 0xcafebabe, 0xbebafeca:
		return true
	}
	return false
}

// isPE reports whether buf begins with the MZ DOS header.
func isPE(buf []byte) bool {
	return len(buf) >= 2 && bytes.Equal(buf[:2], []byte("MZ"))
}
