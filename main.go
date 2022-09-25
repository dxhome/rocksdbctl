package main

import (
	"fmt"
	"strings"

	"github.com/linxGnu/grocksdb"
	"github.com/spf13/cobra"
)

var (
	isDebug    = false
	dbPath     = ""
	keyPrefix  = ""
	dumpLimit  = 50
	showInByte = false
)

type rdbConnect struct {
	db *grocksdb.DB
	ro *grocksdb.ReadOptions
	wo *grocksdb.WriteOptions
}

func connectRocksDB() *rdbConnect {
	var err error

	conn := new(rdbConnect)

	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(false)

	conn.db, err = grocksdb.OpenDb(opts, dbPath)
	if err != nil {
		return nil
	}

	conn.wo = grocksdb.NewDefaultWriteOptions()
	conn.ro = grocksdb.NewDefaultReadOptions()

	return conn
}

func newDumpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "dump all key/value",
		Run:   runDumpCmd,
		Args:  cobra.ExactArgs(0),
	}

	cmd.Flags().StringVar(&keyPrefix, "prefix", "", "key prefix")
	cmd.Flags().IntVarP(&dumpLimit, "limit", "l", 50, "dump limited keys")

	return cmd
}

func runDumpCmd(cmd *cobra.Command, args []string) {
	var count int

	conn := connectRocksDB()
	if conn == nil {
		fmt.Printf("cannot connect to rocksdb at %s\n", dbPath)
		return
	}

	it := conn.db.NewIterator(conn.ro)
	defer it.Close()

	it.Seek([]byte(keyPrefix))
	for it = it; it.Valid(); it.Next() {
		key := it.Key()

		if !strings.HasPrefix(string(key.Data()), keyPrefix) {
			continue
		}

		value := it.Value()
		fmt.Printf("(%s) -> %d bytes\n", key.Data(), len(value.Data()))
		key.Free()
		value.Free()

		count += 1
		if count >= dumpLimit {
			break
		}
	}
	if err := it.Err(); err != nil {
		fmt.Printf("db Iterator failed, %s", err)
		return
	}
	fmt.Println("-----------------")
	fmt.Printf("%d keys dumped\n", count)
}

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "get a key/value",
		Run:   runGetCmd,
		Args:  cobra.ExactArgs(1),
	}

	cmd.Flags().BoolVarP(&showInByte, "byte", "b", false, "disply mode in byte, default display mode is string")

	return cmd
}

func runGetCmd(cmd *cobra.Command, args []string) {
	conn := connectRocksDB()
	if conn == nil {
		fmt.Printf("cannot connect to rocksdb at %s\n", dbPath)
		return
	}

	key := args[0]

	value, err := conn.db.Get(conn.ro, []byte(key))
	if err != nil {
		fmt.Printf("db Get failed, %s, %s", key, err)
		return
	}

	fmt.Printf("(%s) %d bytes\n", key, len(value.Data()))
	if showInByte {
		fmt.Printf("%v\n", value.Data())
	} else {
		fmt.Printf("%s\n", value.Data())
	}
}

func runCommand() error {

	rootCmd := &cobra.Command{
		Use:   "rocksdbctl",
		Short: "rocksdb management tool",
	}

	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "path", "p", "", "rocksdb path")
	rootCmd.MarkPersistentFlagRequired("path")

	rootCmd.AddCommand(
		newDumpCmd(),
		newGetCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func main() {
	runCommand()
}
