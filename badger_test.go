package idbenchmark_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/dgraph-io/badger"
)

const (
	syncWrites   = true
	noSyncWrites = false
)

var badgerKey = []byte(idbenchmarkKey)

// Logger is implemented by any logging system that is used for standard logs.
type Logger interface {
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
	Warningf(string, ...interface{})
}

var badgerLogger Logger

func SetLogger(l Logger) { badgerLogger = l }

func Errorf(format string, v ...interface{}) {
	badgerLogger.Errorf(format, v...)
}

func Infof(format string, v ...interface{}) {
	badgerLogger.Infof(format, v...)
}

func Warningf(format string, v ...interface{}) {
	badgerLogger.Warningf(format, v...)
}

type defaultLog struct {
	*log.Logger
}

var defaultLogger = &defaultLog{Logger: log.New(ioutil.Discard, "badger", 0)}

func UseDefaultLogger() { SetLogger(defaultLogger) }

func (l *defaultLog) Errorf(f string, v ...interface{}) {
	l.Printf("ERROR: "+f, v...)
}

func (l *defaultLog) Infof(f string, v ...interface{}) {
	l.Printf("INFO: "+f, v...)
}

func (l *defaultLog) Warningf(f string, v ...interface{}) {
	l.Printf("WARNING: "+f, v...)
}

func badgerConnect(syncWrites bool) (db *badger.DB, err error) {
	badger.SetLogger(defaultLogger)
	opts := badger.DefaultOptions
	opts.Dir = "."
	opts.ValueDir = "."
	opts.SyncWrites = syncWrites
	db, err = badger.Open(opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

func runBadger(b *testing.B, db *badger.DB, bandwidth uint64) {
	seq, err := db.GetSequence(badgerKey, bandwidth)
	if err != nil {
		log.Printf("Sequence badger GetSequence error: %v", err)
		return
	}
	defer seq.Release()
	seq.Next()
	runBadgerSeq(b, db, seq)
}

func runBadgerSeq(b *testing.B, db *badger.DB, seq *badger.Sequence) {
	for n := 0; n < b.N; n++ {
		id, err := seq.Next()
		if err != nil {
			log.Printf("Sequence badger Next error: %v", err)
			break
		}
		if id == 0 {
			log.Printf("id=0")
			break
		}
	}
}

func benchmarkBadger(b *testing.B, bandwidth uint64, syncWrites bool) {
	db, err := badgerConnect(syncWrites)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	b.ResetTimer()
	runBadger(b, db, bandwidth)
	b.StopTimer()
}

func BenchmarkBadger1(b *testing.B) {
	benchmarkBadger(b, 1, syncWrites)
}

func BenchmarkBadger64(b *testing.B) {
	benchmarkBadger(b, 64, syncWrites) // 2^4
}

func BenchmarkBadger256(b *testing.B) {
	benchmarkBadger(b, 256, syncWrites) // 2^8
}

func BenchmarkBadger4096(b *testing.B) {
	benchmarkBadger(b, 4096, syncWrites) // 2^12
}

func BenchmarkBadger65536(b *testing.B) {
	benchmarkBadger(b, 65536, syncWrites) // 2^16
}

func BenchmarkBadgerNoSync1(b *testing.B) {
	benchmarkBadger(b, 1, noSyncWrites)
}

func BenchmarkBadgerNoSync64(b *testing.B) {
	benchmarkBadger(b, 64, noSyncWrites)
}

func BenchmarkBadgerNoSync256(b *testing.B) {
	benchmarkBadger(b, 256, noSyncWrites)
}

func BenchmarkBadgerNoSync4096(b *testing.B) {
	benchmarkBadger(b, 4096, noSyncWrites)
}

func BenchmarkBadgerNoSync65536(b *testing.B) {
	benchmarkBadger(b, 65536, noSyncWrites)
}

func runBadgerParallel(b *testing.B, bandwidth uint64, syncWrites bool) {
	db, err := badgerConnect(syncWrites)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	key := []byte(idbenchmarkKey)

	seq, err := db.GetSequence(key, bandwidth)
	if err != nil {
		log.Printf("Sequence badger GetSequence error: %v", err)
		return
	}
	defer seq.Release()
	seq.Next()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runBadgerSeq(b, db, seq)
		}
	})
	b.StopTimer()
}

func BenchmarkBadgerParallel1(b *testing.B) {
	runBadgerParallel(b, 1, syncWrites)
}

func BenchmarkBadgerParallel10(b *testing.B) {
	runBadgerParallel(b, 10, syncWrites)
}

func BenchmarkBadgerParallel100(b *testing.B) {
	runBadgerParallel(b, 100, syncWrites)
}

func BenchmarkBadgerParallel1000(b *testing.B) {
	runBadgerParallel(b, 1000, syncWrites)
}

func BenchmarkBadgerParallel10000(b *testing.B) {
	runBadgerParallel(b, 10000, syncWrites)
}

func BenchmarkBadgerParallel100000(b *testing.B) {
	runBadgerParallel(b, 100000, syncWrites)
}

func BenchmarkBadgerParallelNoSync1(b *testing.B) {
	runBadgerParallel(b, 1, noSyncWrites)
}

func BenchmarkBadgerParallelNoSync10(b *testing.B) {
	runBadgerParallel(b, 10, noSyncWrites)
}

func BenchmarkBadgerParallelNoSync100(b *testing.B) {
	runBadgerParallel(b, 100, noSyncWrites)
}

func BenchmarkBadgerParallelNoSync1000(b *testing.B) {
	runBadgerParallel(b, 1000, noSyncWrites)
}

func BenchmarkBadgerParallelNoSync10000(b *testing.B) {
	runBadgerParallel(b, 10000, noSyncWrites)
}

func BenchmarkBadgerParallelNoSync100000(b *testing.B) {
	runBadgerParallel(b, 100000, noSyncWrites)
}
