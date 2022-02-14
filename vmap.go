package vmap

import (
	"math/rand"
	"time"
	"os"
	"fmt"
	"strings"
	"strconv"
)

type Vmap struct {
	hashList [65536]uint32
	Kvmap *[]KV
}

type Hentry struct {
	Entry uint32
}

type KV struct {
	key string
	value string
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))


func Init()(vmap *Vmap) {
	var kv Vmap
	return &kv
}

func StringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func String(length int) string {
  return StringWithCharset(length, charset)
}

func StrCmp(s1 string, s2 string) (bool) {

	b1 := []byte(s1)
	b2 := []byte(s2)
	if len(b1) != len(b2) {return false}
	for i:=0; i<len(b1); i++ {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
}

func Hashdjb2(s string) (rhash uint16) {
    var hash uint16 = 5381

    for _, c := range s {
        hash = ((hash << 5) + hash) + uint16(c)
        // the above line is an optimized version of the following line:
        //hash = hash * 33 + int64(c)
        // which is easier to read and understand...
    }
    return hash
}

func (vmap *Vmap) ResetHash() {
	for i:=0; i< len(vmap.hashList); i++ {
		vmap.hashList[i] = 0
	}
}

func (vmap *Vmap) MakeKVFile(infilnam string, outfilnam string, minValLen int, maxValLen int)(err error) {
    var bufslice []byte
    var kvlist [2000]KV
    var kv []KV

//	minNamLen := 10
//	maxNamLen := 25

	if len(infilnam) < 1 {infilnam = "SortedNameList.txt"}

    infil, err := os.Open(infilnam)
    if err != nil {
        return fmt.Errorf("error MakeKVFile:: opening input file: %s! %v", infilnam, err)
    }

	fileinfo,_ := os.Stat(infilnam)
	size := fileinfo.Size()

    defer infil.Close()

    outfil, err := os.Create(outfilnam)
    if err != nil {
        return fmt.Errorf("error MakeKVFile:: creating output file! %v", err)
    }

	defer outfil.Close()

	bufslice = make([]byte, size)
    _, err = infil.Read(bufslice)
    if err != nil {
        return fmt.Errorf("error MakeNamFile:: reading from input file! %v", err)
    }

//    fmt.Println("bytes read: ", nb)

    // reading first line
    ist :=0
    for i:=0; i< 50; i++ {
        if bufslice[i] == '\n' {
        ist = i
        break
        }
    }

    istate := 0
    namNum := 0
    for i:=ist; i< len(bufslice); i++ {
//  for i:=ist; i< 100; i++ {
//      fmt.Printf("i: %d char: %q state %d\n", i, bufslice[i], istate)
        switch bufslice[i] {
            case '\n', '\r':
                switch istate {
                    case 0:
                        istate = 0
                    case 1:
                        kvlist[namNum].key = string(bufslice[ist:i])
                        namNum++
                        istate = 0
                }
            case ' ', '\t':
                switch istate {
                    case 0:
                        istate = 0
                    case 1:
                        istate = 1
                }
            default:
                switch istate {
                    case 0:
                        istate = 1
                        ist = i
                    case 1:
                        istate = 1
            }
        }
    }

    kv = kvlist[:namNum]
//    namList := *readNames
	dif := maxValLen - minValLen

    for i:=0; i< namNum; i++ {
		valLen := minValLen + seededRand.Intn(dif)
        kv[i].value = String(valLen)
    }

// eliminate duplicate keys
	dup:=0
	for i:=1; i<namNum; i++ {
		key := kv[i].key
		for j:= i+1; j< namNum; j++ {
			if key == kv[j].key {
				kv[j].key = ""
				dup++
			}
		}
	}


	icount:=0
	outstr :=""
    for j:= 0; j< namNum; j++ {
		if len(kv[j].key) > 0 {
        	outstr += fmt.Sprintf("%s:%s\n",kv[j].key, kv[j].value)
			icount++
		}
    }
    hdstr := fmt.Sprintf("%d\n", icount)
    outfil.WriteString(hdstr)
    outfil.WriteString(outstr)

	return nil
}

func (vmap *Vmap) ReadKV(filnam string, numNames int)(err error) {
    var bufslice []byte

	// function reads the key (names) from the file with file name filnam up to the number numNames and fills an array *namptr
	// the first key with index 0 is never filled so references to a zer index can be used as indicator for a key match fail

//fmt.Println("numNames: ", numNames)
    infil, err := os.Open(filnam)
    if err != nil {
        return fmt.Errorf("error readKV:: opening input file: %s! %v", filnam, err)
    }

	fileinfo,_ := os.Stat(filnam)
	size := fileinfo.Size()

    defer infil.Close()

	bufslice = make([]byte, size)
    _, err = infil.Read(bufslice)
  	if err != nil {
        return fmt.Errorf("error readKV:: reading from input file! %v", err)
    }

    // reading first line
    ist :=0

    for i:=0; i< 80; i++ {
        if bufslice[i] == '\n' {
			ist = i
			break
        }
    }

	entries, err := strconv.Atoi(string(bufslice[:ist]))
  	if err != nil {
        return fmt.Errorf("error readKV:: cannot read number of entries: %s! %v", string(bufslice[:ist]), err)
    }

	if (numNames > entries) || (numNames == 0) {numNames = entries}

	kvList := make([]KV, numNames+1)

    istate := 0
	// we start with 1 rather than 0 in order to be able to  utilise 0 as the no-match reference
    namNum := 1
    for i:=ist; i< len(bufslice); i++ {
//  for i:=ist; i< 100; i++ {
//      fmt.Printf("i: %d char: %q state %d\n", i, bufslice[i], istate)
        switch bufslice[i] {
            case '\n', '\r':
                switch istate {
                    case 0:
                        istate = 0
                    case 3:
                        kvList[namNum].value = string(bufslice[ist:i])
						namNum++
                        istate = 0
					default:
        				return fmt.Errorf("error ReadKV:: parsing error unexpected '\n' nameNum %d", namNum)
                }
			case ':':
				switch istate {
					case 1:
						istate = 2
                        kvList[namNum].key = string(bufslice[ist:i])
					default:
        				return fmt.Errorf("error ReadKV:: parsing error unexpected ':' nameNum %d", namNum)
				}
            case ' ', '\t':
                switch istate {
                    case 0:

                    case 2:

					default:
                }
// letter
            default:
                switch istate {
					case 0:
                        istate = 1
                        ist = i
					case 1:
                    case 2:
                        istate = 3
                        ist = i
					case 3:

            }
        }
//fmt.Printf("namNum: %d\n", namNum)
		if namNum > numNames {break}
    }

	vmap.Kvmap = &kvList
    return nil
}

func ReadKVRand(filnam string, numNames int)(kvptr *[]KV, err error) {
    var bufslice []byte

	// function reads the key (names) from the file with file name filnam up to the number numNames and fills an array *namptr
	// the first key with index 0 is never filled so references to a zer index can be used as indicator for a key match fail

//fmt.Println("numNames: ", numNames)
    infil, err := os.Open(filnam)
    if err != nil {
        return nil, fmt.Errorf("error readNames:: opening input file: %s! %v", filnam, err)
    }

	fileinfo,_ := os.Stat(filnam)
	size := fileinfo.Size()

    defer infil.Close()

	bufslice = make([]byte, size)
    _, err = infil.Read(bufslice)
  	if err != nil {
        return nil, fmt.Errorf("error readNames:: reading from input file! %v", err)
    }

    // reading first line
	// first line contains number of entries

    ist :=0

    for i:=0; i< 80; i++ {
        if bufslice[i] == '\n' {
			ist = i
			break
        }
    }

	s := string(bufslice[0:ist])
//	fmt.Printf("name numbers: %s\n",s)
	FilNamNum, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("error ReadKVRand:: converting %s to int: %v",s, err)
	}

	if FilNamNum < numNames +1 { return nil, fmt.Errorf("error ReadKVRand:: number of names is greater than the stored number: %d", FilNamNum) }
// first entry is always a blank entry, so we need to allocate numManes plus 1

	kvlist := make([]KV, FilNamNum+1)

    istate := 0
	// we start with 1 rather than 0 in order to be able to  utilise 0 as the no-match reference
    namNum := 1
    for i:=ist; i< len(bufslice); i++ {
//  for i:=ist; i< 100; i++ {
//      fmt.Printf("i: %d char: %q state %d\n", i, bufslice[i], istate)
        switch bufslice[i] {
            case '\n', '\r':
                switch istate {
                    case 0:
                        istate = 0
                    case 3:
                        kvlist[namNum].value = string(bufslice[ist:i])
						namNum++
                        istate = 0
					default:
        				return nil, fmt.Errorf("error readNames:: parsing error unexpected '\n' nameNum %d", namNum)
                }
			case ':':
				switch istate {
					case 1:
						istate = 2
                        kvlist[namNum].key = string(bufslice[ist:i])
					default:
        				return nil, fmt.Errorf("error readNames:: parsing error unexpected ':' nameNum %d", namNum)
				}
            case ' ', '\t':
                switch istate {
                    case 0:

                    case 2:

					default:
                }
// letter
            default:
                switch istate {
					case 0:
                        istate = 1
                        ist = i
					case 1:
                    case 2:
                        istate = 3
                        ist = i
					case 3:

            }
        }
    }

// eliminate duplicate names
	dup:=0
	for i:=1; i<namNum; i++ {
		key := kvlist[i].key
		for j:= i+1; j< namNum; j++ {
			if key == kvlist[j].key {
				kvlist[j].key = ""
				dup++
			}
		}
	}


	nkvList := make([]KV, numNames+1)
	icount:=1
	for j:=1; j<FilNamNum; j++ {
//		idx := seededRand.Intn(FilNamNum)+1
		idx :=j
		if len(kvlist[idx].key) > 0 {
			nkvList[icount].key = kvlist[idx].key
			nkvList[icount].value = kvlist[idx].value
			icount++
		}
		if icount > numNames { break}
	}

    return &nkvList, nil
}


func WriteNames(filnam string, names *[]string, maxNumber int)(err error) {

    outfil, err := os.Create(filnam)
    if err != nil {
        return fmt.Errorf("error creating output file! %v", err)
    }

	namslice := *names
	namNum := len(namslice)
    if (maxNumber > 0) && (maxNumber < namNum) {namNum = maxNumber}
    outstr := fmt.Sprintf("%d\n", namNum)
    outfil.WriteString(outstr)
    for j:= 0; j< namNum; j++ {
        outstr = fmt.Sprintf("%s\n",(*names)[j])
        outfil.WriteString(outstr)
    }

	return nil
}


func FillMap (nameptr *[]string) (namMap map[string]string, err error) {
	var s string
	var kv []string
	numName := len(*nameptr)
	namMap = make(map[string]string, numName)

	for i:=0; i<numName; i++ {
		s = (*nameptr)[i]
		kv = strings.Split(s,":")
		namMap[kv[0]] = kv[1]
	}

	return namMap, nil
}

func FillMap2 (kvlist *[]KV) (namMap map[string]string, err error) {

	numKV := len(*kvlist)
	namMap = make(map[string]string, numKV)

	for i:=1; i<numKV; i++ {
		kv := (*kvlist)[i]
		namMap[kv.key] = kv.value
	}

	return namMap, nil
}

func (vmap *Vmap)FillVmap(kvlist *[]KV) (err error) {

	if kvlist == nil {return fmt.Errorf(" error FillVmap:: no list of KV provided!") }

	kv := *kvlist
	// start with 1 rather than 0, so that the 0 value can be used to indicate a failed match
	coll :=0
	for i:=1; i< len(*kvlist); i++ {
		ihash := Hashdjb2(kv[i].key)
		if vmap.hashList[ihash] == 0 {
			vmap.hashList[ihash] = (uint32)(i)
		} else {
			coll++
			entry := vmap.hashList[ihash]+1
			tkey := kv[i].key
			tval := kv[i].value
			kv[i].key = kv[entry].key
			kv[i].value = kv[entry].value
			kv[entry].key = tkey
			kv[entry].value = tval
			for j:=1; j< len(vmap.hashList); j++ {
				if vmap.hashList[j] == entry {
					vmap.hashList[j] = uint32(i)
					break
				}
			}
		}
//{return fmt.Errorf("error FillVmap entry %d collision %d!",i, ihash)}
//fmt.Printf("i: %d ihash: %d Entry: %d\n", i, ihash, hashList[ihash])
	}
	vmap.Kvmap = &kv
	return nil
}

func (vmap *Vmap)GetVal(key string)(value string, ok bool){

	kv := vmap.Kvmap
	ihash := Hashdjb2(key)
	index := vmap.hashList[ihash]
	if index == 0 {return "", false}
	if (*kv)[index].key == key {
		return (*kv)[index].value, true
	} else {
		if (*kv)[index +1].key == key {
			return (*kv)[index+1].value, true
		} else {
			return "", false
		}
	}
}

func (vmap *Vmap)AddKVToVMap(key string, val string) (success bool, err error){

	ihash := Hashdjb2(key)
	index := vmap.hashList[ihash]
	if index != 0 {return false, fmt.Errorf("error AddKVMap:: key %s already exists!", key)}
	var newkv KV
	newkv.key = key
	newkv.value = val

	kvlist := *(vmap.Kvmap)
	kvlist = append(kvlist, newkv)
	vmap.hashList[ihash] = uint32(len(kvlist))

	return true, nil
}

func (vmap *Vmap) DelKVToVMap(key string) (success bool, err error){

	ihash := Hashdjb2(key)
	index := vmap.hashList[ihash]
	if index == 0 {return false, fmt.Errorf("error DelKVMap:: key %s does not exist!", key)}

	kvlist := *(vmap.Kvmap)
	vmap.hashList[ihash] = 0
	kvlist[index].key = ""
	kvlist[index].value = ""
	return true, nil
}

func (vmap *Vmap) VacuumVMap(){
	var icount uint32

	kvlist := *(vmap.Kvmap)
	icount=1
	for i:=1; i<len(kvlist); i++ {
		if len(kvlist[i].key) >0 {
			kvlist[icount].key = kvlist[i].key
			kvlist[icount].value = kvlist[i].value
			ihash := Hashdjb2(kvlist[icount].key)
			vmap.hashList[ihash] = icount
			icount++
		}
	}
	nkvlist := kvlist[:icount]
	vmap.Kvmap = &nkvlist
	return
}


func (vmap *Vmap)SaveVmap(filnam string)(err error) {

    outfil, err := os.Create(filnam)
    if err != nil {
        return fmt.Errorf("error SaveVmap:: creating output file! %v", err)
    }
	defer outfil.Close()

	outstr :=""
	kv := vmap.Kvmap
	icount:=0
	for i:=0; i< len(*kv); i++ {
		if len((*kv)[i].key)<1 {continue}
		outstr += fmt.Sprintf("%s:%s\n",(*kv)[i].key, (*kv)[i].value)
		icount++
    }

	hdstr := fmt.Sprintf("%d/n", icount)
    outfil.WriteString(hdstr)
	outfil.WriteString(outstr)
	return nil
}



func (vmap *Vmap)PrintVmap(filnam string) (err error){

    outfil, err := os.Create(filnam)
    if err != nil {
        return fmt.Errorf("error SaveVmap:: creating output file! %v", err)
    }

	kv := vmap.Kvmap
	outstr := ""
	for i:=0; i< len(*kv); i++ {
        outstr += fmt.Sprintf("entry: %3d key: %-10s value: %-25s\n", i, (*kv)[i].key, (*kv)[i].value)
    }

	outfil.WriteString(outstr)
	outfil.Close()
	return nil
}

func (vmap *Vmap)PrintAltVmap(filnam string) (err error){

    outfil, err := os.Create(filnam)
    if err != nil {
        return fmt.Errorf("error SaveVmap:: creating output file! %v", err)
    }
	kv := vmap.Kvmap
	outstr := ""
	for i:=1; i< len(*kv); i++ {
		key := (*kv)[i].key
		ihash := Hashdjb2(key)
		idx := vmap.hashList[ihash]
        outstr += fmt.Sprintf("entry: %3d key: %-10s value: %-25s idx: %6d Hash: %6d\n", 
		i, (*kv)[i].key, (*kv)[i].value, idx, ihash)
    }

	outfil.WriteString(outstr)
	outfil.Close()
	return nil
}
