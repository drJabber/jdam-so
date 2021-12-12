import ctypes
import json

lib = ctypes.cdll.LoadLibrary("./jdam.so")
lib.jdam_fuzz.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
lib.jdam_fuzz.restype = ctypes.POINTER(ctypes.c_ubyte*8)

options = {
	"Mutators": "",
	"Ignore": "",
	"Output": "",
	"Count": 1,
	"Rounds" : 5,   
	"NilChance" : 0.75,
	"MaxDepth":  100,
	"Seed": 0,
	"List": False,
	"Verbose": False,
	"Version": False
}

tofuzz = {"code": "qwerty", "data": ["asdf", {"alpha":"beta", "gamma": [5,5,"puppet",5]}]}
ptr = lib.jdam_fuzz(json.dumps(options).encode("utf-8"), json.dumps(tofuzz).encode("utf-8"))
length = int.from_bytes(ptr.contents, byteorder="little")
data = bytes(ctypes.cast(ptr,
            ctypes.POINTER(ctypes.c_ubyte*(8 + length))
            ).contents[8:])
print(data)
