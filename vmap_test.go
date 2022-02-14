package vmap

import (
	"testing"
//	"fmt"
	)

func Test1(t *testing.T){
	
	kvM := Init()
	entries := 100
    kvlist, err := ReadKVRand("NKVList.txt", entries)
    if err != nil {
        t.Errorf("error reading names! %v\n", err)
		return
    }

	kvM.Kvmap = kvlist
	err = kvM.PrintVmap("KVListin.out")
    if err != nil {
        t.Errorf("error printing vmap! %v\n", err)
	}

    err = kvM.FillVmap(kvlist)
    if err != nil {
        t.Errorf("error filling vmap! %v\n", err)
	}

	err = kvM.PrintVmap("KVListnew.out")
    if err != nil {
        t.Errorf("error printing vmap! %v\n", err)
	}

	lenkv := len(*kvlist)
	if lenkv != entries + 1 {
		t.Errorf("kv list has %d entries, should have % d!", lenkv, entries+1)
	}
	err = kvM.PrintAltVmap("KVListAlt.out")
    if err != nil {
        t.Errorf("error printing vmap alt! %v\n", err)
	}

	kvlist = kvM.Kvmap
	for i:=1; i<lenkv; i++ {
		key:= (*kvlist)[i].key
		val, succ := kvM.GetVal(key)
		if !succ { t.Errorf("key %s idx: %d not found!", key, i)}
		if val != (*kvlist)[i].value {
			t.Errorf("vmap line %d key %s value %s is not rec value %s", i, key, val, (*kvlist)[i].value)
		}
	}

	kvmap, err := FillMap2(kvlist)
    if err != nil {
        t.Errorf("error filling map! %v\n", err)
	}

	for i:=1; i<lenkv; i++ {
		key:= (*kvlist)[i].key
		val, ok := kvmap[key]
		if !ok { t.Errorf("key %s not found!", key)}
		if val != (*kvlist)[i].value {
			t.Errorf("map entry %d key %s value %s is not rec value %s", i, key, val, (*kvlist)[i].value)
		}
	}

}

func BenchmarkVmap (b *testing.B) {

	kvM := Init()
	entries := 500
	kvM.ResetHash()

    kvlist, err := ReadKVRand("NKVList.txt", entries)
    if err != nil {
		b.Errorf("error reading names! %v\n", err)
    }
	if len(*kvlist) -1 != entries {
		b.Errorf("error kvlist length %d entries %d",len(*kvlist), entries)
	}

	kvM.Kvmap = kvlist

    err = kvM.FillVmap(kvlist)
    if err != nil {
        b.Errorf("error filling vmap! %v\n", err)
	}
//		lenkv := len(*kvlist)
		b.ResetTimer()

	for i:=0; i<b.N; i++ {
		idx := seededRand.Intn(entries)+1
		key := (*kvlist)[idx].key
			val, succ := kvM.GetVal(key)
			if !succ { b.Errorf("key %s not found!", key)}
			if val != (*kvlist)[idx].value {
				b.Errorf("vmap value %s is not rec value %s", val, (*kvlist)[idx].value)
			}
	}
}

func BenchmarkGoMap (b *testing.B) {

	kvM := Init()
	entries := 500
	kvM.ResetHash()

    kvlist, err := ReadKVRand("NKVList.txt", entries)
    if err != nil {
		b.Errorf("error reading names! %v\n", err)
    }
	if len(*kvlist) -1 != entries {
		b.Errorf("error kvlist length %d entries %d",len(*kvlist), entries)
	}

   kvmap, err := FillMap2(kvlist)
    if err != nil {
        b.Errorf("error filling map! %v\n", err)
	}
//		lenkv := len(*kvlist)
		b.ResetTimer()

		for i:=0; i<b.N; i++ {
			idx := seededRand.Intn(entries)+1
			key := (*kvlist)[idx].key
			val, ok := kvmap[key]
			if !ok {b.Errorf("key %s not found!", key)}
			if val != (*kvlist)[idx].value {
				b.Errorf("vmap value %s is not rec value %s", val, (*kvlist)[idx].value)
			}
		}
}


func BenchmarkVmapNH (b *testing.B) {

	var keysel [100]string
	entries := 500
	kvM := Init()
	kvM.ResetHash()

    kvlist, err := ReadKVRand("NKVList.txt", entries)
    if err != nil {
		b.Errorf("error reading names! %v\n", err)
    }
	if len(*kvlist) -1 != entries {
		b.Errorf("error kvlist length %d entries %d",len(*kvlist), entries)
	}

	err = kvM.FillVmap(kvlist)
    if err != nil {
        b.Errorf("error filling vmap! %v\n", err)
	}
//		lenkv := len(*kvlist)

	keylen := 5 + seededRand.Intn(10)
	for i:=0; i< 100; i++ {
		keysel[i] = String(keylen)
	}

	b.ResetTimer()

	for i:=0; i<b.N; i++ {
		key := keysel[seededRand.Intn(100)]
		_, succ := kvM.GetVal(key)
			if succ { b.Errorf("key %s should not be in map!", key)}
	}
}

func BenchmarkGoMapNH (b *testing.B) {
	var keysel [100]string

	entries := 500
	kvM := Init()
	kvM.ResetHash()

    kvlist, err := ReadKVRand("NKVList.txt", entries)
    if err != nil {
		b.Errorf("error reading names! %v\n", err)
    }
	if len(*kvlist) -1 != entries {
		b.Errorf("error kvlist length %d entries %d",len(*kvlist), entries)
	}
   kvmap, err := FillMap2(kvlist)
    if err != nil {
        b.Errorf("error filling map! %v\n", err)
	}

	keylen := 5 + seededRand.Intn(10)
	for i:=0; i< 100; i++ {
		keysel[i] = String(keylen)
	}

	b.ResetTimer()

	for i:=0; i<b.N; i++ {
		key := keysel[seededRand.Intn(100)]
		_, ok := kvmap[key]
		if ok {b.Errorf("key %s found!", key)}
	}
}

func BenchmarkCmpStr(b *testing.B) {
	var keysel [100]string

	keylen := 5 + seededRand.Intn(10)
	for i:=0; i< 100; i++ {
		keysel[i] = String(keylen)
	}

	b.ResetTimer()

//	res := false
	for i:=0; i<b.N; i++ {
		key1 := keysel[seededRand.Intn(100)]
		key2 := keysel[seededRand.Intn(100)]

		if key1 == key2 {
//			res = true
		}
	}
//	fmt.Println("res: ", res)
}

func BenchmarkCmpStr2(b *testing.B) {
	var keysel [100]string

	keylen := 5 + seededRand.Intn(10)
	for i:=0; i< 100; i++ {
		keysel[i] = String(keylen)
	}

	b.ResetTimer()

//	res := false
	for i:=0; i<b.N; i++ {
		key1 := keysel[seededRand.Intn(100)]
		key2 := keysel[seededRand.Intn(100)]
		_ = StrCmp(key1, key2)
	}
//	fmt.Println("res: ", res)
}
