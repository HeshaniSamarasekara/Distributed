package model

// FileTable -struct fro file tabel
type FileTable struct {
	Files []FileTableEntry
}

// FileTableEntry - Struct of type file details
type FileTableEntry struct {
	IP          string
	Port        string
	FileStrings []string
}

// NodeFiles - List of files in the node
type NodeFiles struct {
	FileNames []string
}
