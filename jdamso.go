package main

import (
    "C"
    "unsafe"
    "encoding/json"
    "encoding/binary"
    "strings"
    "jdamso/log"
    "gitlab.com/michenriksen/jdam/pkg/jdam" 	
    "gitlab.com/michenriksen/jdam/pkg/jdam/mutation"
)

type options struct {
	Mutators  *string
	Ignore    *string
	Output    *string
	Count     *int
	Rounds    *int
	NilChance *float64
	MaxDepth  *int
	Seed      *int64
	List      *bool
	Verbose   *bool
	Version   *bool
}

//export jdam_fuzz
func jdam_fuzz(config *C.char, target *C.char) unsafe.Pointer {
	logger := log.Logger{}

	options := options{}
  err := json.Unmarshal([]byte(strings.TrimSpace(C.GoString(config))), &options)
  if (err != nil){
    logger.Fatal("Unable to parse input as JSON object: %s\n", err.Error())
  }

	var mutators mutation.MutatorList
	if *options.Mutators != "" {
		ids := strings.Split(*options.Mutators, ",")
		mutators = mutatorListFromIDs(ids)
	} else {
		mutators = mutation.Mutators
	}

	subject := map[string]interface{}{}
	err = json.Unmarshal([]byte(strings.TrimSpace(C.GoString(target))), &subject)
	if err != nil {
		logger.Fatal("Unable to parse input as JSON object: %s\n", err.Error())
	}


  var fuzzer *jdam.Fuzzer

  if *options.Seed != 0 {
    logger.Info("Using PRNG seed: %d\n", *options.Seed)
    fuzzer = jdam.NewWithSeed(*options.Seed, mutators)
  } else {
    fuzzer = jdam.New(mutators)
  }
  fuzzer.IgnoreFields(strings.Split(*options.Ignore, ",")).
    NilChance(*options.NilChance).
    MaxDepth(*options.MaxDepth)

  var fuzzedJSON []byte
  fuzzed := copyMap(subject)

  for j := 0; j < *options.Rounds; j++ {
    logger.Info("[OBJECT #%d] Fuzzing round: %d\n", 1, j+1)
    fuzzed = fuzzer.Fuzz(fuzzed)
  }
  fuzzedJSON, err = json.Marshal(fuzzed)
  if err != nil {
    logger.Fatal("[OBJECT #%d] Unable to marshal fuzzed object as JSON: %s\n", 1, err.Error())
  }

  length := make([]byte, 8)
  binary.LittleEndian.PutUint64(length, uint64(len(fuzzedJSON)))
  return C.CBytes(append(length,fuzzedJSON...))
}

func mutatorListFromIDs(ids []string) mutation.MutatorList {
	var mutators mutation.MutatorList
	for _, m := range mutation.Mutators {
		for _, id := range ids {
			if id == m.ID() {
				mutators = append(mutators, m)
				break
			}
		}
	}
	return mutators
}

func copyMap(orig map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range orig {
		cp[k] = v
	}
	return cp
}

func main(){}